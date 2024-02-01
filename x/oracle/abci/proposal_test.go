package abci_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/stretchr/testify/require"

	"github.com/Team-Kujira/core/x/oracle/types"

	"github.com/Team-Kujira/core/x/oracle/abci"
	"github.com/Team-Kujira/core/x/oracle/keeper"
	cometabci "github.com/cometbft/cometbft/abci/types"
	cmtproto "github.com/cometbft/cometbft/proto/tendermint/types"
	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/ed25519"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/module"
	stakingkeeper "github.com/cosmos/cosmos-sdk/x/staking/keeper"
	protoio "github.com/cosmos/gogoproto/io"
	"github.com/cosmos/gogoproto/proto"
)

var ValAddrs = keeper.ValAddrs
var ValPubKeys []cryptotypes.PubKey
var ValPrivKeys []*ed25519.PrivKey

func init() {
	for i := 0; i < 5; i++ {
		privKey := ed25519.GenPrivKey()
		ValPrivKeys = append(ValPrivKeys, privKey)
		pubKey := &ed25519.PubKey{Key: privKey.PubKey().Bytes()}
		ValPubKeys = append(ValPubKeys, pubKey)
	}
}

func SetupTest(t *testing.T) (keeper.TestInput, *abci.ProposalHandler) {
	input := keeper.CreateTestInput(t)

	power := int64(100)
	amt := sdk.TokensFromConsensusPower(power, sdk.DefaultPowerReduction)
	sh := stakingkeeper.NewMsgServerImpl(&input.StakingKeeper)
	ctx := input.Ctx
	mm := module.NewManager()

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
		mm,  // module manager
		nil, // mempool
		nil, // bApp
	)

	params := types.DefaultParams()
	params.RequiredDenoms = []string{"BTC", "ETH"}
	input.OracleKeeper.SetParams(ctx, params)

	return input, h
}

func TestGetBallotByDenom(t *testing.T) {
	_, h := SetupTest(t)
	power := int64(100)

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
	input, h := SetupTest(t)

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
	input.OracleKeeper.SetParams(input.Ctx, params)

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

func TestCompareOraclePrices(t *testing.T) {
	testCases := []struct {
		p1       map[string]math.LegacyDec
		p2       map[string]math.LegacyDec
		expEqual bool
	}{
		{
			p1: map[string]math.LegacyDec{
				"ATOM": math.LegacyNewDec(10),
				"ETH":  math.LegacyNewDec(2300),
			},
			p2: map[string]math.LegacyDec{
				"ATOM": math.LegacyNewDec(10),
				"ETH":  math.LegacyNewDec(2300),
			},
			expEqual: true,
		},
		{
			p1: map[string]math.LegacyDec{
				"ATOM": math.LegacyNewDec(10),
				"ETH":  math.LegacyNewDec(2300),
			},
			p2: map[string]math.LegacyDec{
				"ATOM": math.LegacyNewDec(10),
				"ETH":  math.LegacyNewDec(2200),
			},
			expEqual: false,
		},
		{
			p1: map[string]math.LegacyDec{
				"ATOM": math.LegacyNewDec(10),
				"ETH":  math.LegacyNewDec(2300),
			},
			p2: map[string]math.LegacyDec{
				"ATOM": math.LegacyNewDec(10),
				"ETH":  math.LegacyNewDec(2300),
				"BTC":  math.LegacyNewDec(43000),
			},
			expEqual: false,
		},
		{
			p1: map[string]math.LegacyDec{
				"ATOM": math.LegacyNewDec(10),
				"ETH":  math.LegacyNewDec(2300),
				"BTC":  math.LegacyNewDec(43000),
			},
			p2: map[string]math.LegacyDec{
				"ATOM": math.LegacyNewDec(10),
				"ETH":  math.LegacyNewDec(2300),
			},
			expEqual: false,
		},
		{
			p1: nil,
			p2: map[string]math.LegacyDec{
				"ATOM": math.LegacyNewDec(10),
				"ETH":  math.LegacyNewDec(2300),
			},
			expEqual: false,
		},
		{
			p1: map[string]math.LegacyDec{
				"ATOM": math.LegacyNewDec(10),
				"ETH":  math.LegacyNewDec(2300),
			},
			p2:       nil,
			expEqual: false,
		},
		{
			p1:       nil,
			p2:       nil,
			expEqual: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			err := abci.CompareOraclePrices(tc.p1, tc.p2)
			if tc.expEqual {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestCompareMissMap(t *testing.T) {
	testCases := []struct {
		m1       map[string]sdk.ValAddress
		m2       map[string]sdk.ValAddress
		expEqual bool
	}{
		{
			m1: map[string]sdk.ValAddress{
				ValAddrs[0].String(): ValAddrs[0],
				ValAddrs[1].String(): ValAddrs[1],
			},
			m2: map[string]sdk.ValAddress{
				ValAddrs[0].String(): ValAddrs[0],
				ValAddrs[1].String(): ValAddrs[1],
			},
			expEqual: true,
		},
		{
			m1: map[string]sdk.ValAddress{
				ValAddrs[0].String(): ValAddrs[0],
				ValAddrs[1].String(): ValAddrs[1],
			},
			m2: map[string]sdk.ValAddress{
				ValAddrs[0].String(): ValAddrs[0],
				ValAddrs[2].String(): ValAddrs[2],
			},
			expEqual: false,
		},
		{
			m1: map[string]sdk.ValAddress{
				ValAddrs[0].String(): ValAddrs[0],
				ValAddrs[1].String(): ValAddrs[1],
			},
			m2: map[string]sdk.ValAddress{
				ValAddrs[0].String(): ValAddrs[0],
				ValAddrs[1].String(): ValAddrs[2],
			},
			expEqual: false,
		},
		{
			m1: map[string]sdk.ValAddress{
				ValAddrs[0].String(): ValAddrs[0],
				ValAddrs[1].String(): ValAddrs[1],
			},
			m2: map[string]sdk.ValAddress{
				ValAddrs[0].String(): ValAddrs[0],
				ValAddrs[2].String(): ValAddrs[1],
			},
			expEqual: false,
		},
		{
			m1: map[string]sdk.ValAddress{
				ValAddrs[0].String(): ValAddrs[0],
				ValAddrs[1].String(): ValAddrs[1],
			},
			m2: map[string]sdk.ValAddress{
				ValAddrs[0].String(): ValAddrs[0],
			},
			expEqual: false,
		},
		{
			m1: map[string]sdk.ValAddress{
				ValAddrs[0].String(): ValAddrs[0],
			},
			m2: map[string]sdk.ValAddress{
				ValAddrs[0].String(): ValAddrs[0],
				ValAddrs[1].String(): ValAddrs[1],
			},
			expEqual: false,
		},
		{
			m1:       map[string]sdk.ValAddress{},
			m2:       nil,
			expEqual: true,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("Case %d", i), func(t *testing.T) {
			err := abci.CompareMissMap(tc.m1, tc.m2)
			if tc.expEqual {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestPrepareProposal(t *testing.T) {
	input, h := SetupTest(t)

	params := types.DefaultParams()
	params.RequiredDenoms = []string{"BTC", "ETH"}
	input.OracleKeeper.SetParams(input.Ctx, params)

	handler := h.PrepareProposal()

	consParams := input.Ctx.ConsensusParams()
	consParams.Abci = &tmproto.ABCIParams{
		VoteExtensionsEnableHeight: 2,
	}
	input.Ctx = input.Ctx.WithConsensusParams(consParams)

	// Handler before vote extension enable
	res, err := handler(input.Ctx, &cometabci.RequestPrepareProposal{
		Height:          1,
		Txs:             [][]byte{},
		LocalLastCommit: cometabci.ExtendedCommitInfo{},
	})
	require.NoError(t, err)
	require.Len(t, res.Txs, 0)

	// Invalid vote extension data
	invalidLocalLastCommit := cometabci.ExtendedCommitInfo{
		Votes: []cometabci.ExtendedVoteInfo{
			{
				Validator: cometabci.Validator{
					Address: ValPubKeys[0].Address().Bytes(),
					Power:   1,
				},
				VoteExtension: []byte{},
			},
			{
				Validator: cometabci.Validator{
					Address: ValPubKeys[1].Address().Bytes(),
					Power:   1,
				},
				VoteExtension: []byte{},
			},
		},
	}
	_, err = handler(input.Ctx, &cometabci.RequestPrepareProposal{
		Height:          2,
		Txs:             [][]byte{},
		LocalLastCommit: invalidLocalLastCommit,
	})
	require.Error(t, err)

	// Valid vote extension data
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
	marshalDelimitedFn := func(msg proto.Message) ([]byte, error) {
		var buf bytes.Buffer
		if err := protoio.NewDelimitedWriter(&buf).WriteMsg(msg); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}

	cve := cmtproto.CanonicalVoteExtension{
		Extension: voteExt1Bytes,
		Height:    3 - 1, // the vote extension was signed in the previous height
		Round:     1,
		ChainId:   input.Ctx.ChainID(),
	}

	ext1SignBytes, err := marshalDelimitedFn(&cve)
	require.NoError(t, err)

	signature1, err := ValPrivKeys[0].Sign(ext1SignBytes)
	require.NoError(t, err)

	cve = cmtproto.CanonicalVoteExtension{
		Extension: voteExt2Bytes,
		Height:    3 - 1, // the vote extension was signed in the previous height
		Round:     1,
		ChainId:   input.Ctx.ChainID(),
	}

	ext2SignBytes, err := marshalDelimitedFn(&cve)
	require.NoError(t, err)

	signature2, err := ValPrivKeys[1].Sign(ext2SignBytes)
	require.NoError(t, err)

	localLastCommit := cometabci.ExtendedCommitInfo{
		Round: 1,
		Votes: []cometabci.ExtendedVoteInfo{
			{
				BlockIdFlag: cmtproto.BlockIDFlagCommit,
				Validator: cometabci.Validator{
					Address: ValPubKeys[0].Address().Bytes(),
					Power:   1,
				},
				VoteExtension:      voteExt1Bytes,
				ExtensionSignature: signature1,
			},
			{
				BlockIdFlag: cmtproto.BlockIDFlagCommit,
				Validator: cometabci.Validator{
					Address: ValPubKeys[1].Address().Bytes(),
					Power:   1,
				},
				VoteExtension:      voteExt2Bytes,
				ExtensionSignature: signature2,
			},
		},
	}
	res, err = handler(input.Ctx, &cometabci.RequestPrepareProposal{
		Height:          3,
		Txs:             [][]byte{},
		LocalLastCommit: localLastCommit,
	})
	require.NoError(t, err)
	require.Len(t, res.Txs, 1) // Check injectedVoteExtTx
	injectedVoteExtTx := abci.StakeWeightedPrices{
		StakeWeightedPrices: map[string]math.LegacyDec{
			"ETH": math.LegacyNewDec(2180),
			"BTC": math.LegacyNewDec(25000),
		},
		ExtendedCommitInfo: localLastCommit,
		MissCounter: map[string]sdk.ValAddress{
			ValAddrs[2].String(): ValAddrs[2],
		},
	}
	injectedBytes, err := json.Marshal(injectedVoteExtTx)
	require.NoError(t, err)
	require.Equal(t, string(injectedBytes), string(res.Txs[0]))
}

func TestProcessProposal(t *testing.T) {
	input, h := SetupTest(t)

	handler := h.ProcessProposal()

	consParams := input.Ctx.ConsensusParams()
	consParams.Abci = &tmproto.ABCIParams{
		VoteExtensionsEnableHeight: 2,
	}
	input.Ctx = input.Ctx.WithConsensusParams(consParams)

	// Valid vote extension data
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
	marshalDelimitedFn := func(msg proto.Message) ([]byte, error) {
		var buf bytes.Buffer
		if err := protoio.NewDelimitedWriter(&buf).WriteMsg(msg); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}

	cve := cmtproto.CanonicalVoteExtension{
		Extension: voteExt1Bytes,
		Height:    3 - 1, // the vote extension was signed in the previous height
		Round:     1,
		ChainId:   input.Ctx.ChainID(),
	}

	ext1SignBytes, err := marshalDelimitedFn(&cve)
	require.NoError(t, err)

	signature1, err := ValPrivKeys[0].Sign(ext1SignBytes)
	require.NoError(t, err)

	cve = cmtproto.CanonicalVoteExtension{
		Extension: voteExt2Bytes,
		Height:    3 - 1, // the vote extension was signed in the previous height
		Round:     1,
		ChainId:   input.Ctx.ChainID(),
	}

	ext2SignBytes, err := marshalDelimitedFn(&cve)
	require.NoError(t, err)

	signature2, err := ValPrivKeys[1].Sign(ext2SignBytes)
	require.NoError(t, err)

	localLastCommit := cometabci.ExtendedCommitInfo{
		Round: 1,
		Votes: []cometabci.ExtendedVoteInfo{
			{
				BlockIdFlag: cmtproto.BlockIDFlagCommit,
				Validator: cometabci.Validator{
					Address: ValPubKeys[0].Address().Bytes(),
					Power:   1,
				},
				VoteExtension:      voteExt1Bytes,
				ExtensionSignature: signature1,
			},
			{
				BlockIdFlag: cmtproto.BlockIDFlagCommit,
				Validator: cometabci.Validator{
					Address: ValPubKeys[1].Address().Bytes(),
					Power:   1,
				},
				VoteExtension:      voteExt2Bytes,
				ExtensionSignature: signature2,
			},
		},
	}

	// Invalid missMap calculation
	injectedVoteExtTx := abci.StakeWeightedPrices{
		StakeWeightedPrices: map[string]math.LegacyDec{
			"ETH": math.LegacyNewDec(2180),
			"BTC": math.LegacyNewDec(25000),
		},
		ExtendedCommitInfo: localLastCommit,
		MissCounter: map[string]sdk.ValAddress{
			ValAddrs[1].String(): ValAddrs[1],
		},
	}
	injectedBytes, err := json.Marshal(injectedVoteExtTx)
	require.NoError(t, err)

	res, err := handler(input.Ctx, &cometabci.RequestProcessProposal{
		Txs:                [][]byte{injectedBytes},
		ProposedLastCommit: cometabci.CommitInfo{},
		Misbehavior:        []cometabci.Misbehavior{},
		Hash:               []byte{},
		Height:             3,
		Time:               time.Time{},
		NextValidatorsHash: []byte{},
		ProposerAddress:    []byte{},
	})
	require.NoError(t, err)
	require.Equal(t, res.Status, cometabci.ResponseProcessProposal_REJECT)

	// Invalid stake weighted prices calculation
	injectedVoteExtTx = abci.StakeWeightedPrices{
		StakeWeightedPrices: map[string]math.LegacyDec{
			"ETH": math.LegacyNewDec(2180),
			"BTC": math.LegacyNewDec(25500),
		},
		ExtendedCommitInfo: localLastCommit,
		MissCounter: map[string]sdk.ValAddress{
			ValAddrs[2].String(): ValAddrs[2],
		},
	}
	injectedBytes, err = json.Marshal(injectedVoteExtTx)
	require.NoError(t, err)

	res, err = handler(input.Ctx, &cometabci.RequestProcessProposal{
		Txs:                [][]byte{injectedBytes},
		ProposedLastCommit: cometabci.CommitInfo{},
		Misbehavior:        []cometabci.Misbehavior{},
		Hash:               []byte{},
		Height:             3,
		Time:               time.Time{},
		NextValidatorsHash: []byte{},
		ProposerAddress:    []byte{},
	})
	require.NoError(t, err)
	require.Equal(t, res.Status, cometabci.ResponseProcessProposal_REJECT)

	// Empty txs
	res, err = handler(input.Ctx, &cometabci.RequestProcessProposal{
		Txs:                [][]byte{},
		ProposedLastCommit: cometabci.CommitInfo{},
		Misbehavior:        []cometabci.Misbehavior{},
		Hash:               []byte{},
		Height:             3,
		Time:               time.Time{},
		NextValidatorsHash: []byte{},
		ProposerAddress:    []byte{},
	})
	require.NoError(t, err)
	require.Equal(t, res.Status, cometabci.ResponseProcessProposal_ACCEPT)

	// Not decode-able last tx
	res, err = handler(input.Ctx, &cometabci.RequestProcessProposal{
		Txs:                [][]byte{{0x0}},
		ProposedLastCommit: cometabci.CommitInfo{},
		Misbehavior:        []cometabci.Misbehavior{},
		Hash:               []byte{},
		Height:             3,
		Time:               time.Time{},
		NextValidatorsHash: []byte{},
		ProposerAddress:    []byte{},
	})
	require.NoError(t, err)
	require.Equal(t, res.Status, cometabci.ResponseProcessProposal_ACCEPT)

	// Accurate vote extension
	injectedVoteExtTx = abci.StakeWeightedPrices{
		StakeWeightedPrices: map[string]math.LegacyDec{
			"ETH": math.LegacyNewDec(2180),
			"BTC": math.LegacyNewDec(25000),
		},
		ExtendedCommitInfo: localLastCommit,
		MissCounter: map[string]sdk.ValAddress{
			ValAddrs[2].String(): ValAddrs[2],
		},
	}
	injectedBytes, err = json.Marshal(injectedVoteExtTx)
	require.NoError(t, err)

	res, err = handler(input.Ctx, &cometabci.RequestProcessProposal{
		Txs:                [][]byte{injectedBytes},
		ProposedLastCommit: cometabci.CommitInfo{},
		Misbehavior:        []cometabci.Misbehavior{},
		Hash:               []byte{},
		Height:             3,
		Time:               time.Time{},
		NextValidatorsHash: []byte{},
		ProposerAddress:    []byte{},
	})
	require.NoError(t, err)
	require.Equal(t, res.Status, cometabci.ResponseProcessProposal_ACCEPT)
}

func TestPreBlocker(t *testing.T) {
	input, h := SetupTest(t)

	input.OracleKeeper.SetOraclePrices(input.Ctx, map[string]math.LegacyDec{
		"LTC": math.LegacyNewDec(100),
	})
	_, err := input.OracleKeeper.GetExchangeRate(input.Ctx, "LTC")
	require.NoError(t, err)

	consParams := input.Ctx.ConsensusParams()
	consParams.Abci = &tmproto.ABCIParams{
		VoteExtensionsEnableHeight: 2,
	}
	input.Ctx = input.Ctx.WithConsensusParams(consParams)

	// Valid vote extension data
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
	marshalDelimitedFn := func(msg proto.Message) ([]byte, error) {
		var buf bytes.Buffer
		if err := protoio.NewDelimitedWriter(&buf).WriteMsg(msg); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}

	cve := cmtproto.CanonicalVoteExtension{
		Extension: voteExt1Bytes,
		Height:    3 - 1, // the vote extension was signed in the previous height
		Round:     1,
		ChainId:   input.Ctx.ChainID(),
	}

	ext1SignBytes, err := marshalDelimitedFn(&cve)
	require.NoError(t, err)

	signature1, err := ValPrivKeys[0].Sign(ext1SignBytes)
	require.NoError(t, err)

	cve = cmtproto.CanonicalVoteExtension{
		Extension: voteExt2Bytes,
		Height:    3 - 1, // the vote extension was signed in the previous height
		Round:     1,
		ChainId:   input.Ctx.ChainID(),
	}

	ext2SignBytes, err := marshalDelimitedFn(&cve)
	require.NoError(t, err)

	signature2, err := ValPrivKeys[1].Sign(ext2SignBytes)
	require.NoError(t, err)

	localLastCommit := cometabci.ExtendedCommitInfo{
		Round: 1,
		Votes: []cometabci.ExtendedVoteInfo{
			{
				BlockIdFlag: cmtproto.BlockIDFlagCommit,
				Validator: cometabci.Validator{
					Address: ValPubKeys[0].Address().Bytes(),
					Power:   1,
				},
				VoteExtension:      voteExt1Bytes,
				ExtensionSignature: signature1,
			},
			{
				BlockIdFlag: cmtproto.BlockIDFlagCommit,
				Validator: cometabci.Validator{
					Address: ValPubKeys[1].Address().Bytes(),
					Power:   1,
				},
				VoteExtension:      voteExt2Bytes,
				ExtensionSignature: signature2,
			},
		},
	}

	// Accurate vote extension
	injectedVoteExtTx := abci.StakeWeightedPrices{
		StakeWeightedPrices: map[string]math.LegacyDec{
			"ETH": math.LegacyNewDec(2180),
			"BTC": math.LegacyNewDec(25000),
		},
		ExtendedCommitInfo: localLastCommit,
		MissCounter: map[string]sdk.ValAddress{
			ValAddrs[2].String(): ValAddrs[2],
		},
	}
	injectedBytes, err := json.Marshal(injectedVoteExtTx)
	require.NoError(t, err)

	res, err := h.PreBlocker(input.Ctx, &cometabci.RequestFinalizeBlock{
		Txs:                [][]byte{injectedBytes},
		DecidedLastCommit:  cometabci.CommitInfo{},
		Misbehavior:        []cometabci.Misbehavior{},
		Hash:               []byte{},
		Height:             3,
		Time:               time.Time{},
		NextValidatorsHash: []byte{},
		ProposerAddress:    []byte{},
	})
	require.NoError(t, err)
	require.Equal(t, res.ConsensusParamsChanged, false)

	ethPrice, err := input.OracleKeeper.GetExchangeRate(input.Ctx, "ETH")
	require.Equal(t, ethPrice.String(), injectedVoteExtTx.StakeWeightedPrices["ETH"].String())
	btcPrice, err := input.OracleKeeper.GetExchangeRate(input.Ctx, "BTC")
	require.Equal(t, btcPrice.String(), injectedVoteExtTx.StakeWeightedPrices["BTC"].String())
	_, err = input.OracleKeeper.GetExchangeRate(input.Ctx, "LTC")
	require.Error(t, err)

	val0MissCount := input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[0])
	require.Equal(t, val0MissCount, uint64(0))
	val1MissCount := input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[1])
	require.Equal(t, val1MissCount, uint64(0))
	val2MissCount := input.OracleKeeper.GetMissCounter(input.Ctx, ValAddrs[2])
	require.Equal(t, val2MissCount, uint64(1))
}
