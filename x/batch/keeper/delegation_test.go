package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	distrtestutil "github.com/cosmos/cosmos-sdk/x/distribution/testutil"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestWithdrawAllDelegationRewards() {
	// Setting up validators and delegations
	valTokens := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)

	delPub := secp256k1.GenPrivKey().PubKey()
	delAddr := sdk.AccAddress(delPub.Address())
	suite.app.AccountKeeper.SetAccount(suite.ctx, authtypes.NewBaseAccount(delAddr, nil, 0, 0))

	delTokens := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)

	for i := 0; i < 3; i++ {
		validator, err := distrtestutil.CreateValidator(valConsPks[i], delTokens)
		require.NoError(suite.T(), err)

		validator, _ = validator.SetInitialCommission(stakingtypes.NewCommission(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), math.LegacyNewDec(0)))
		validator, _ = validator.AddTokensFromDel(valTokens)
		suite.app.StakingKeeper.SetValidator(suite.ctx, validator)
		suite.app.StakingKeeper.SetValidatorByConsAddr(suite.ctx, validator)
		suite.app.StakingKeeper.SetValidatorByPowerIndex(suite.ctx, validator)

		// call the after-creation hook
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

	// next block
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

	// Add tokens to distribution module for reward distribution
	distrModuleToken := sdk.NewCoin(sdk.DefaultBondDenom, delTokens)
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{distrModuleToken})
	suite.app.BankKeeper.SendCoinsFromModuleToModule(suite.ctx, minttypes.ModuleName, distrtypes.ModuleName, sdk.Coins{distrModuleToken})
	require.NoError(suite.T(), err)

	// Simulate rewards for the delegations
	reward := sdk.DecCoins{
		sdk.NewInt64DecCoin("stake", 10),
	}
	for i := 0; i < 3; i++ {
		valAddr := sdk.ValAddress(valConsAddrs[i])
		suite.app.DistrKeeper.AllocateTokensToValidator(suite.ctx, suite.app.StakingKeeper.Validator(suite.ctx, valAddr), reward)
	}

	// Withdraw all rewards
	rewards, err := suite.app.BatchKeeper.WithdrawAllDelegationRewards(suite.ctx, delAddr)
	require.NoError(suite.T(), err)

	// Validate if rewards are properly claimed
	require.False(suite.T(), rewards.IsZero())

	// Try second withdrawal at same block
	rewards, err = suite.app.BatchKeeper.WithdrawAllDelegationRewards(suite.ctx, delAddr)
	require.NoError(suite.T(), err)

	// Validate if remaining rewards are zero
	require.True(suite.T(), rewards.IsZero())
}
