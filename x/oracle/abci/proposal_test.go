package abci_test

import (
	"encoding/json"
	"sort"
	"testing"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/Team-Kujira/core/x/oracle/types"

	"github.com/Team-Kujira/core/x/oracle/abci"
	"github.com/Team-Kujira/core/x/oracle/keeper"
	cometabci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
)

var ValAddrs = keeper.ValAddrs
var ValPubKeys = keeper.ValPubKeys

func TestGetBallotByDenom(t *testing.T) {
	input := keeper.CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power, sdk.DefaultPowerReduction)
	sh := stakingkeeper.NewMsgServerImpl(&input.StakingKeeper)
	ctx := input.Ctx

	// Validator created
	_, err := sh.CreateValidator(ctx, keeper.NewTestMsgCreateValidator(ValAddrs[0], ValPubKeys[0], amt))
	require.NoError(t, err)
	_, err = sh.CreateValidator(ctx, keeper.NewTestMsgCreateValidator(ValAddrs[1], ValPubKeys[1], amt))
	require.NoError(t, err)
	_, err = sh.CreateValidator(ctx, keeper.NewTestMsgCreateValidator(ValAddrs[2], ValPubKeys[2], amt))
	require.NoError(t, err)
	input.StakingKeeper.EndBlocker(ctx)

	h := abci.NewProposalHandler(
		input.Ctx.Logger(),
		input.OracleKeeper,
		input.StakingKeeper,
		nil, // module manager
		nil, // mempool
		nil, // bApp
	)

	// organize votes by denom
	voteExt1 := abci.OracleVoteExtension{
		Height: 1,
		Prices: map[string]math.LegacyDec{
			"BTC": math.LegacyNewDec(25000),
			"ETH": math.LegacyNewDec(2200),
		},
	}
	voteExt2 := abci.OracleVoteExtension{
		Height: 1,
		Prices: map[string]math.LegacyDec{
			"BTC": math.LegacyNewDec(25030),
			"ETH": math.LegacyNewDec(2180),
		},
	}
	voteExt1Bytes, err := json.Marshal(voteExt1)
	require.NoError(t, err)
	voteExt2Bytes, err := json.Marshal(voteExt2)
	require.NoError(t, err)

	consAddrMap := map[string]sdk.ValAddress{
		sdk.ConsAddress(ValPubKeys[0].Address().Bytes()).String(): ValAddrs[0],
		sdk.ConsAddress(ValPubKeys[1].Address().Bytes()).String(): ValAddrs[1],
		sdk.ConsAddress(ValPubKeys[2].Address().Bytes()).String(): ValAddrs[2],
	}

	ballotMap := h.GetBallotByDenom(cometabci.ExtendedCommitInfo{
		Votes: []cometabci.ExtendedVoteInfo{
			{
				Validator: cometabci.Validator{
					Address: ValPubKeys[0].Address().Bytes(),
				},
				VoteExtension: voteExt1Bytes,
			},
			{
				Validator: cometabci.Validator{
					Address: ValPubKeys[1].Address().Bytes(),
				},
				VoteExtension: voteExt2Bytes,
			},
		},
	}, map[string]types.Claim{
		ValAddrs[0].String(): {
			Power:     power,
			WinCount:  0,
			Recipient: ValAddrs[0],
		},
		ValAddrs[1].String(): {
			Power:     power,
			WinCount:  0,
			Recipient: ValAddrs[1],
		},
		ValAddrs[2].String(): {
			Power:     power,
			WinCount:  0,
			Recipient: ValAddrs[2],
		},
	}, consAddrMap)

	ethBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(math.LegacyNewDec(2200), "ETH", ValAddrs[0], power),
		types.NewVoteForTally(math.LegacyNewDec(2180), "ETH", ValAddrs[1], power),
	}
	btcBallot := types.ExchangeRateBallot{
		types.NewVoteForTally(math.LegacyNewDec(25000), "BTC", ValAddrs[0], power),
		types.NewVoteForTally(math.LegacyNewDec(25030), "BTC", ValAddrs[1], power),
	}
	sort.Sort(ethBallot)
	sort.Sort(btcBallot)
	sort.Sort(ballotMap["ETH"])
	sort.Sort(ballotMap["BTC"])

	require.Equal(t, ethBallot, ballotMap["ETH"])
	require.Equal(t, btcBallot, ballotMap["BTC"])
}

func TestComputeStakeWeightedPricesAndMissMap(t *testing.T) {
	input := keeper.CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power, sdk.DefaultPowerReduction)
	sh := stakingkeeper.NewMsgServerImpl(&input.StakingKeeper)
	ctx := input.Ctx

	// Validator created
	_, err := sh.CreateValidator(ctx, keeper.NewTestMsgCreateValidator(ValAddrs[0], ValPubKeys[0], amt))
	require.NoError(t, err)
	_, err = sh.CreateValidator(ctx, keeper.NewTestMsgCreateValidator(ValAddrs[1], ValPubKeys[1], amt))
	require.NoError(t, err)
	_, err = sh.CreateValidator(ctx, keeper.NewTestMsgCreateValidator(ValAddrs[2], ValPubKeys[2], amt))
	require.NoError(t, err)
	input.StakingKeeper.EndBlocker(ctx)

	h := abci.NewProposalHandler(
		input.Ctx.Logger(),
		input.OracleKeeper,
		input.StakingKeeper,
		nil, // module manager
		nil, // mempool
		nil, // bApp
	)

	// organize votes by denom
	voteExt1 := abci.OracleVoteExtension{
		Height: 1,
		Prices: map[string]math.LegacyDec{
			"BTC": math.LegacyNewDec(25000),
			"ETH": math.LegacyNewDec(2200),
		},
	}
	voteExt2 := abci.OracleVoteExtension{
		Height: 1,
		Prices: map[string]math.LegacyDec{
			"BTC": math.LegacyNewDec(25030),
			"ETH": math.LegacyNewDec(2180),
		},
	}
	voteExt1Bytes, err := json.Marshal(voteExt1)
	require.NoError(t, err)
	voteExt2Bytes, err := json.Marshal(voteExt2)
	require.NoError(t, err)

	params := types.DefaultParams()
	params.RequiredDenoms = []string{"BTC", "ETH"}
	input.OracleKeeper.SetParams(ctx, params)

	stakeWeightedPrices, missMap, err := h.ComputeStakeWeightedPricesAndMissMap(input.Ctx, cometabci.ExtendedCommitInfo{
		Votes: []cometabci.ExtendedVoteInfo{
			{
				Validator: cometabci.Validator{
					Address: ValPubKeys[0].Address().Bytes(),
					Power:   1,
				},
				VoteExtension: voteExt1Bytes,
			},
			{
				Validator: cometabci.Validator{
					Address: ValPubKeys[1].Address().Bytes(),
					Power:   1,
				},
				VoteExtension: voteExt2Bytes,
			},
		},
	})

	require.Equal(t, math.LegacyNewDec(2180).String(), stakeWeightedPrices["ETH"].String())
	require.Equal(t, math.LegacyNewDec(25000).String(), stakeWeightedPrices["BTC"].String())
	require.Nil(t, missMap[ValAddrs[0].String()])
	require.Nil(t, missMap[ValAddrs[1].String()])
	require.NotNil(t, missMap[ValAddrs[2].String()])
}
