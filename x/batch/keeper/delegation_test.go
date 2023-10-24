package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtestutil "github.com/cosmos/cosmos-sdk/x/distribution/testutil"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/stretchr/testify/require"
)

func (suite *KeeperTestSuite) TestWithdrawAllDelegationRewards() {
	// Setting up validators and delegations
	valTokens := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)

	delAddr := sdk.AccAddress([]byte("delegator"))
	delTokens := sdk.TokensFromConsensusPower(100, sdk.DefaultPowerReduction)

	for i := 0; i < 3; i++ {
		validator, err := distrtestutil.CreateValidator(valConsPks[i], delTokens)
		require.NoError(suite.T(), err)

		validator, _ = validator.SetInitialCommission(stakingtypes.NewCommission(sdk.NewDecWithPrec(5, 1), sdk.NewDecWithPrec(5, 1), math.LegacyNewDec(0)))
		validator, _ = validator.AddTokensFromDel(valTokens)
		suite.app.StakingKeeper.SetValidator(suite.ctx, validator)

		// Delegate to the validator
		delAmount := sdk.NewCoin(sdk.DefaultBondDenom, delTokens)
		err = suite.app.BankKeeper.MintCoins(suite.ctx, stakingtypes.ModuleName, sdk.Coins{delAmount})
		require.NoError(suite.T(), err)

		_, err = suite.app.StakingKeeper.Delegate(suite.ctx, delAddr, delAmount.Amount, stakingtypes.Unbonded, validator, true)
		require.NoError(suite.T(), err)
	}

	// Simulate rewards for the delegations
	// This is a simplified way to generate rewards, in reality, you might need to simulate blocks or other activities
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
	expectedRewards := sdk.NewCoins(sdk.NewInt64Coin("stake", 30)) // 3 validators * 10 stake each
	require.True(suite.T(), rewards.IsEqual(expectedRewards))
}
