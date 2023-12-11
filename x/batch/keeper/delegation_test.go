package keeper_test

import (
	"fmt"

	"cosmossdk.io/math"
	batchkeeper "github.com/Team-Kujira/core/x/batch/keeper"
	batchtypes "github.com/Team-Kujira/core/x/batch/types"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrkeeper "github.com/cosmos/cosmos-sdk/x/distribution/keeper"
	distrtestutil "github.com/cosmos/cosmos-sdk/x/distribution/testutil"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	simtestutil "github.com/cosmos/cosmos-sdk/testutil/sims"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
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

	delTokens := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)

	for i := 0; i < totalVals; i++ {
		// Setup the validator
		validator, err := distrtestutil.CreateValidator(valConsPks[i], delTokens)
		suite.Require().NoError(err)

		validator, _ = validator.SetInitialCommission(stakingtypes.NewCommission(math.LegacyNewDecWithPrec(5, 1), math.LegacyNewDecWithPrec(5, 1), math.LegacyNewDec(0)))
		validator, _ = validator.AddTokensFromDel(valTokens)
		err = suite.app.StakingKeeper.SetValidator(sdk.WrapSDKContext(suite.ctx), validator)
		suite.NoError(err)
		err = suite.app.StakingKeeper.SetValidatorByConsAddr(sdk.WrapSDKContext(suite.ctx), validator)
		suite.NoError(err)
		err = suite.app.StakingKeeper.SetValidatorByPowerIndex(sdk.WrapSDKContext(suite.ctx), validator)
		suite.NoError(err)

		// Call the after-creation hook
		valAddr, err := sdk.ValAddressFromBech32(validator.GetOperator())
		suite.Require().NoError(err)
		err = suite.app.StakingKeeper.Hooks().AfterValidatorCreated(sdk.WrapSDKContext(suite.ctx), valAddr)
		suite.NoError(err)

		// Delegate to the validator
		delAmount := sdk.NewCoin(sdk.DefaultBondDenom, delTokens)
		err = suite.app.BankKeeper.MintCoins(sdk.WrapSDKContext(suite.ctx), minttypes.ModuleName, sdk.Coins{delAmount})
		suite.app.BankKeeper.SendCoinsFromModuleToAccount(sdk.WrapSDKContext(suite.ctx), minttypes.ModuleName, delAddr, sdk.Coins{delAmount})
		suite.Require().NoError(err)

		_, err = suite.app.StakingKeeper.Delegate(sdk.WrapSDKContext(suite.ctx), delAddr, delAmount.Amount, stakingtypes.Unbonded, validator, true)
		suite.Require().NoError(err)
	}

	// =====================Next block =====================
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

	// Add tokens to distribution module for reward distribution
	distrModuleTokens := sdk.Coins{}
	for i := 0; i < totalTokens; i++ {
		distrModuleTokens = append(distrModuleTokens, sdk.NewCoin(fmt.Sprintf("stake-%d", i), delTokens))
	}
	distrModuleTokens = distrModuleTokens.Sort()

	err := suite.app.BankKeeper.MintCoins(sdk.WrapSDKContext(suite.ctx), minttypes.ModuleName, distrModuleTokens)
	suite.app.BankKeeper.SendCoinsFromModuleToModule(sdk.WrapSDKContext(suite.ctx), minttypes.ModuleName, distrtypes.ModuleName, distrModuleTokens)
	suite.Require().NoError(err)

	// Allocate rewards to validators
	valRewardTokens := sdk.DecCoins{}
	expDelegatorRewards := sdk.Coins{}
	for i := 0; i < totalTokens; i++ {
		valRewardTokens = valRewardTokens.Add(sdk.NewInt64DecCoin(fmt.Sprintf("stake-%d", i), 10))
		expDelegatorRewards = expDelegatorRewards.Add(sdk.NewInt64Coin(fmt.Sprintf("stake-%d", i), 1))
	}
	for i := 0; i < totalVals; i++ {
		valAddr := sdk.ValAddress(valConsAddrs[i])
		validator, err := suite.app.StakingKeeper.Validator(sdk.WrapSDKContext(suite.ctx), valAddr)
		suite.Require().NoError(err)
		suite.app.DistrKeeper.AllocateTokensToValidator(sdk.WrapSDKContext(suite.ctx), validator, valRewardTokens)
	}

	// Withdraw all rewards using a single batch transaction
	gasForBatchWithdrawal := suite.ctx.GasMeter().GasConsumed()
	res, err := batchMsgServer.WithdrawAllDelegatorRewards(sdk.WrapSDKContext(suite.ctx), batchtypes.NewMsgWithdrawAllDelegatorRewards(delAddr))
	gasForBatchWithdrawal = suite.ctx.GasMeter().GasConsumed() - gasForBatchWithdrawal
	suite.Require().NoError(err)
	suite.Require().False(res.Amount.IsZero())
	totalBatchRewards := res.Amount

	// check if there are no pending rewards for any validator
	for i := 0; i < totalVals; i++ {
		valAddr := sdk.ValAddress(valConsAddrs[i])
		val, err := suite.app.StakingKeeper.Validator(sdk.WrapSDKContext(suite.ctx), valAddr)
		suite.Require().NoError(err)
		delegation, err := suite.app.StakingKeeper.Delegation(sdk.WrapSDKContext(suite.ctx), delAddr, valAddr)
		suite.Require().NoError(err)
		endingPeriod, err := suite.app.DistrKeeper.IncrementValidatorPeriod(sdk.WrapSDKContext(suite.ctx), val)
		suite.Require().NoError(err)
		rewards, err := suite.app.DistrKeeper.CalculateDelegationRewards(sdk.WrapSDKContext(suite.ctx), val, delegation, endingPeriod)
		suite.Require().NoError(err)
		suite.Require().True(rewards.IsZero())
	}

	// ===================== Next block =====================
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

	// Allocate rewards to validators
	for i := 0; i < totalVals; i++ {
		valAddr := sdk.ValAddress(valConsAddrs[i])
		validator, err := suite.app.StakingKeeper.Validator(sdk.WrapSDKContext(suite.ctx), valAddr)
		suite.Require().NoError(err)
		suite.app.DistrKeeper.AllocateTokensToValidator(sdk.WrapSDKContext(suite.ctx), validator, valRewardTokens)
		suite.Require().NoError(err)
	}

	totalGasForIndividualWithdrawals := suite.ctx.GasMeter().GasConsumed()
	totalIndividualRewards := sdk.Coins{}
	// Withdraw all rewards using multiple individual transactions
	for i := 0; i < totalVals; i++ {
		valAddr := sdk.ValAddress(valConsAddrs[i])
		// Withdraw rewards
		res, err := distrMsgServer.WithdrawDelegatorReward(sdk.WrapSDKContext(suite.ctx), distrtypes.NewMsgWithdrawDelegatorReward(delAddr.String(), valAddr.String()))
		suite.Require().NoError(err)
		suite.Require().False(res.Amount.IsZero())
		// check individual rewards are accurate
		suite.Require().Equal(res.Amount.String(), expDelegatorRewards.String())
		totalIndividualRewards = totalIndividualRewards.Add(res.Amount...)
	}
	// check if rewards are same for batch execution and individual executions
	suite.Require().Equal(totalIndividualRewards.String(), totalBatchRewards.String())
	// check gas being reduced using batch operation
	totalGasForIndividualWithdrawals = suite.ctx.GasMeter().GasConsumed() - totalGasForIndividualWithdrawals
	suite.Require().True(gasForBatchWithdrawal < totalGasForIndividualWithdrawals)
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

	delTokens := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)

	validators := []string{}
	amounts := []math.Int{}
	for i := 0; i < totalVals; i++ {
		// Setup the validator
		validator, err := distrtestutil.CreateValidator(valConsPks[i], delTokens)
		suite.Require().NoError(err)

		validator, _ = validator.SetInitialCommission(stakingtypes.NewCommission(math.LegacyNewDecWithPrec(5, 1), math.LegacyNewDecWithPrec(5, 1), math.LegacyNewDec(0)))
		validator, _ = validator.AddTokensFromDel(valTokens)
		suite.app.StakingKeeper.SetValidator(sdk.WrapSDKContext(suite.ctx), validator)
		suite.app.StakingKeeper.SetValidatorByConsAddr(sdk.WrapSDKContext(suite.ctx), validator)
		suite.app.StakingKeeper.SetValidatorByPowerIndex(sdk.WrapSDKContext(suite.ctx), validator)

		// Call the after-creation hook
		valAddr, err := sdk.ValAddressFromBech32(validator.GetOperator())
		suite.Require().NoError(err)
		err = suite.app.StakingKeeper.Hooks().AfterValidatorCreated(sdk.WrapSDKContext(suite.ctx), valAddr)
		suite.Require().NoError(err)

		validators = append(validators, validator.GetOperator())
		amounts = append(amounts, math.NewInt(int64(500*(i%5)))) // 0, 500, 1000, 1500, 2000

		// Mint coins for delegation
		delAmount := sdk.NewCoin(sdk.DefaultBondDenom, delTokens)
		err = suite.app.BankKeeper.MintCoins(sdk.WrapSDKContext(suite.ctx), minttypes.ModuleName, sdk.Coins{delAmount})
		suite.app.BankKeeper.SendCoinsFromModuleToAccount(sdk.WrapSDKContext(suite.ctx), minttypes.ModuleName, delAddr, sdk.Coins{delAmount})
		suite.Require().NoError(err)

		// setup initial delegation
		_, err = suite.app.StakingKeeper.Delegate(sdk.WrapSDKContext(suite.ctx), delAddr, math.NewInt(1000), stakingtypes.Unbonded, validator, true)
		suite.Require().NoError(err)
	}

	// ===================== First cache context =====================
	cacheCtx1, _ := suite.ctx.CacheContext()

	// Set all delegations using a single batch transaction
	gasForBatchDelegation := cacheCtx1.GasMeter().GasConsumed()
	_, err := batchMsgServer.BatchResetDelegation(sdk.WrapSDKContext(cacheCtx1), batchtypes.NewMsgBatchResetDelegation(delAddr, validators, amounts))
	gasForBatchDelegation = cacheCtx1.GasMeter().GasConsumed() - gasForBatchDelegation
	suite.Require().NoError(err)

	// ===================== Second cache context =====================
	cacheCtx2, _ := suite.ctx.CacheContext()

	totalGasForIndividualDelegations := cacheCtx2.GasMeter().GasConsumed()
	// Delegate using multiple individual transactions
	for i := 0; i < totalVals; i++ {
		valAddr := sdk.ValAddress(valConsAddrs[i])
		existingDelegation := math.NewInt(1000)
		if amounts[i].GT(existingDelegation) {
			_, err := stakingMsgServer.Delegate(
				sdk.WrapSDKContext(cacheCtx2),
				stakingtypes.NewMsgDelegate(
					delAddr.String(),
					valAddr.String(),
					sdk.NewCoin(sdk.DefaultBondDenom, amounts[i].Sub(existingDelegation)),
				),
			)
			suite.Require().NoError(err)
		} else if amounts[i].LT(existingDelegation) {
			_, err := stakingMsgServer.Undelegate(
				sdk.WrapSDKContext(cacheCtx2),
				stakingtypes.NewMsgUndelegate(
					delAddr.String(),
					valAddr.String(),
					sdk.NewCoin(sdk.DefaultBondDenom, existingDelegation.Sub(amounts[i])),
				),
			)
			suite.Require().NoError(err)
		}
	}
	totalGasForIndividualDelegations = cacheCtx2.GasMeter().GasConsumed() - totalGasForIndividualDelegations

	suite.T().Log(">>>>>>> Gas usage for batch delegation is ", gasForBatchDelegation)
	suite.T().Log(">>>>>>> Gas usage for individual delegations is ", totalGasForIndividualDelegations)
}
