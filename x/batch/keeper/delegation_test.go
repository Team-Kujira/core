package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"
	batchkeeper "github.com/Team-Kujira/core/x/batch/keeper"
	batchtypes "github.com/Team-Kujira/core/x/batch/types"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtestutil "github.com/cosmos/cosmos-sdk/x/distribution/testutil"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func generateValidatorKeysAndAddresses(num int) (pubKeys []cryptotypes.PubKey, addrs []sdk.ConsAddress) {
	pubKeys = simtestutil.CreateTestPubKeys(num)
	for _, pubKey := range pubKeys {
		addrs = append(addrs, sdk.ConsAddress(pubKey.Address()))
	}
	return
}

func (suite *KeeperTestSuite) TestWithdrawAllDelegationRewards() {
	// Generating validator key pubkey and addresses
	totalVals := 75
	totalTokens := 100
	valConsPks, valConsAddrs := generateValidatorKeysAndAddresses(totalVals)

	// Setting up msg servers
	distrMsgServer := distrkeeper.NewMsgServerImpl(suite.app.DistrKeeper)
	batchMsgServer := batchkeeper.NewMsgServerImpl(suite.app.BatchKeeper)

	// Setting up validators and delegations
	valTokens := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)

	delPub := secp256k1.GenPrivKey().PubKey()
	delAddr := sdk.AccAddress(delPub.Address())
	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(delAddr, nil, 0, 0))

	delTokens := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)

	for i := 0; i < totalVals; i++ {
		// Setup the validator
		validator, err := distrtestutil.CreateValidator(valConsPks[i], delTokens)
		require.NoError(suite.T(), err)

		validator, _ = validator.SetInitialCommission(stakingtypes.NewCommission(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), math.LegacyNewDec(0)))
		validator, _ = validator.AddTokensFromDel(valTokens)
		suite.app.StakingKeeper.SetValidator(suite.ctx, validator)
		suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
		suite.app.StakingKeeper.SetValidatorByPowerIndex(suite.ctx, validator)

		// Call the after-creation hook
		err = suite.app.StakingKeeper.Hooks().AfterValidatorCreated(suite.ctx, validator.GetOperator())
		require.NoError(suite.T(), err)

		// Delegate to the validator
		delAmount := sdk.NewCoin(sdk.DefaultBondDenom, delTokens)
		err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{delAmount})
		suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, delAddr, sdk.Coins{delAmount})
		require.NoError(suite.T(), err)

		_, err = suite.app.StakingKeeper.Delegate(suite.ctx, delAddr, delAmount.Amount, stakingtypes.Unbonded, validator, true)
		require.NoError(suite.T(), err)
	}

	// =====================Next block =====================
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

	// Add tokens to distribution module for reward distribution
	distrModuleTokens := sdk.Coins{}
	for i := 0; i < totalTokens; i++ {
		distrModuleTokens = append(distrModuleTokens, sdk.NewCoin(fmt.Sprintf("stake-%d", i), delTokens))
	}
	distrModuleTokens = distrModuleTokens.Sort()

	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, distrModuleTokens)
	suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, minttypes.ModuleName, distrtypes.ModuleName, distrModuleTokens)
	require.NoError(suite.T(), err)

	// Allocate rewards to validators
	valRewardTokens := sdk.DecCoins{}
	for i := 0; i < totalTokens; i++ {
		valRewardTokens = append(valRewardTokens, sdk.NewInt64DecCoin(fmt.Sprintf("stake-%d", i), 10))
	}
	valRewardTokens = valRewardTokens.Sort()
	for i := 0; i < totalVals; i++ {
		valAddr := sdk.ValAddress(valConsAddrs[i])
		suite.app.DistrKeeper.AllocateTokensToValidator(suite.ctx, suite.app.StakingKeeper.Validator(suite.ctx, valAddr), valRewardTokens)
	}

	// Withdraw all rewards using a single batch transaction
	gasForBatchWithdrawal := suite.ctx.GasMeter().GasConsumed()
	res, err := batchMsgServer.WithdrawAllDelegatorRewards(suite.ctx, batchtypes.NewMsgWithdrawAllDelegatorRewards(delAddr))
	gasForBatchWithdrawal = suite.ctx.GasMeter().GasConsumed() - gasForBatchWithdrawal
	require.NoError(suite.T(), err)
	require.False(suite.T(), res.Amount.IsZero())

	// ===================== Next block =====================
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

	// Allocate rewards to validators
	for i := 0; i < totalVals; i++ {
		valAddr := sdk.ValAddress(valConsAddrs[i])
		suite.app.DistrKeeper.AllocateTokensToValidator(suite.ctx, suite.app.StakingKeeper.Validator(suite.ctx, valAddr), valRewardTokens)
	}

	totalGasForIndividualWithdrawals := suite.ctx.GasMeter().GasConsumed()
	// Withdraw all rewards using multiple individual transactions
	for i := 0; i < totalVals; i++ {
		valAddr := sdk.ValAddress(valConsAddrs[i])
		// Withdraw rewards
		res, err := distrMsgServer.WithdrawDelegatorReward(suite.ctx, distrtypes.NewMsgWithdrawDelegatorReward(delAddr, valAddr))
		require.NoError(suite.T(), err)
		require.False(suite.T(), res.Amount.IsZero())

	}
	totalGasForIndividualWithdrawals = suite.ctx.GasMeter().GasConsumed() - totalGasForIndividualWithdrawals

	require.True(suite.T(), gasForBatchWithdrawal < totalGasForIndividualWithdrawals)
	suite.T().Log(">>>>>>> Gas usage for batch withdrawals is ", gasForBatchWithdrawal)
	suite.T().Log(">>>>>>> Gas usage for individual withdrawals is ", totalGasForIndividualWithdrawals)
}

func (suite *KeeperTestSuite) TestBatchResetDelegation() {
	// Generating validator key pubkey and addresses
	totalVals := 75
	valConsPks, valConsAddrs := generateValidatorKeysAndAddresses(totalVals)

	// Setting up msg servers
	stakingMsgServer := stakingkeeper.NewMsgServerImpl(suite.app.StakingKeeper)
	batchMsgServer := batchkeeper.NewMsgServerImpl(suite.app.BatchKeeper)

	// Setting up validators and delegations
	valTokens := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)

	delPub := secp256k1.GenPrivKey().PubKey()
	delAddr := sdk.AccAddress(delPub.Address())
	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(delAddr, nil, 0, 0))

	delTokens := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)

	validators := []string{}
	amounts := []sdk.Int{}
	for i := 0; i < totalVals; i++ {
		// Setup the validator
		validator, err := distrtestutil.CreateValidator(valConsPks[i], delTokens)
		require.NoError(suite.T(), err)

		validator, _ = validator.SetInitialCommission(stakingtypes.NewCommission(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), math.LegacyNewDec(0)))
		validator, _ = validator.AddTokensFromDel(valTokens)
		suite.app.StakingKeeper.SetValidator(suite.ctx, validator)
		suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
		suite.app.StakingKeeper.SetValidatorByPowerIndex(suite.ctx, validator)

		// Call the after-creation hook
		err = suite.app.StakingKeeper.Hooks().AfterValidatorCreated(suite.ctx, validator.GetOperator())
		require.NoError(suite.T(), err)

		validators = append(validators, validator.GetOperator().String())
		amounts = append(amounts, sdk.NewInt(1000))

		// Mint coins for delegation
		delAmount := sdk.NewCoin(sdk.DefaultBondDenom, delTokens)
		err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{delAmount})
		suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, delAddr, sdk.Coins{delAmount})
		require.NoError(suite.T(), err)
	}

	// ===================== Next block =====================
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

	// Set all delegations using a single batch transaction
	gasForBatchDelegation := suite.ctx.GasMeter().GasConsumed()
	_, err := batchMsgServer.BatchResetDelegation(suite.ctx, batchtypes.NewMsgBatchResetDelegation(delAddr, validators, amounts))
	gasForBatchDelegation = suite.ctx.GasMeter().GasConsumed() - gasForBatchDelegation
	require.NoError(suite.T(), err)

	// ===================== Next block =====================
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

	totalGasForIndividualDelegations := suite.ctx.GasMeter().GasConsumed()
	// Delegate using multiple individual transactions
	for i := 0; i < totalVals; i++ {
		valAddr := sdk.ValAddress(valConsAddrs[i])
		_, err := stakingMsgServer.Delegate(suite.ctx, stakingtypes.NewMsgDelegate(delAddr, valAddr, sdk.NewCoin(sdk.DefaultBondDenom, amounts[i])))
		require.NoError(suite.T(), err)
	}
	totalGasForIndividualDelegations = suite.ctx.GasMeter().GasConsumed() - totalGasForIndividualDelegations

	require.True(suite.T(), gasForBatchDelegation < totalGasForIndividualDelegations)
	suite.T().Log(">>>>>>> Gas usage for batch delegation is ", gasForBatchDelegation)
	suite.T().Log(">>>>>>> Gas usage for individual delegations is ", totalGasForIndividualDelegations)
}
