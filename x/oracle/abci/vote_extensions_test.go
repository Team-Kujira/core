package abci_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"cosmossdk.io/math"
	"github.com/Team-Kujira/core/x/oracle/abci"
	"github.com/Team-Kujira/core/x/oracle/keeper"
	"github.com/Team-Kujira/core/x/oracle/types"
	cometabci "github.com/cometbft/cometbft/abci/types"
	"github.com/stretchr/testify/require"
)

func TestDecoding(t *testing.T) {
	resBody := []byte(`{"prices":{"BTC":"47375.706652541026694000","ETH":"2649.328939436595054949","USDT":"1.000661260343873178"}}`)
	prices := abci.PricesResponse{}
	err := json.Unmarshal(resBody, &prices)
	require.NoError(t, err)
}

func TestExtendVoteHandler(t *testing.T) {
	input := keeper.CreateTestInput(t)

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{"prices":{"BTC":"47375.706652541026694000","ETH":"2649.328939436595054949","USDT":"1.000661260343873178"}}`))
	}))
	defer func() { testServer.Close() }()

	h := abci.NewVoteExtHandler(input.Ctx.Logger(), input.OracleKeeper)
	handler := h.ExtendVoteHandler(abci.OracleConfig{
		Endpoint: testServer.URL,
	})
	res, err := handler(input.Ctx, &cometabci.RequestExtendVote{
		Hash:               []byte{},
		Height:             3,
		Time:               time.Time{},
		Txs:                [][]byte{},
		ProposedLastCommit: cometabci.CommitInfo{},
		Misbehavior:        []cometabci.Misbehavior{},
		NextValidatorsHash: []byte{},
		ProposerAddress:    []byte{},
	})
	require.NoError(t, err)
	voteExt := types.VoteExtension{}
	err = voteExt.Unmarshal(res.VoteExtension)
	require.NoError(t, err)
	require.Equal(t, voteExt.Height, int64(3))
	require.Len(t, voteExt.Prices, 3)
	require.Equal(t, voteExt.Prices[0].Denom, "BTC")
	require.Equal(t, voteExt.Prices[0].ExchangeRate.String(), "47375.706652541026694000")
	require.Equal(t, voteExt.Prices[1].Denom, "ETH")
	require.Equal(t, voteExt.Prices[1].ExchangeRate.String(), "2649.328939436595054949")
	require.Equal(t, voteExt.Prices[2].Denom, "USDT")
	require.Equal(t, voteExt.Prices[2].ExchangeRate.String(), "1.000661260343873178")
}

func TestVerifyVoteExtensionHandler(t *testing.T) {
	input := keeper.CreateTestInput(t)

	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(http.StatusOK)
		res.Write([]byte(`{"prices":{"BTC":"47375.706652541026694000","ETH":"2649.328939436595054949","USDT":"1.000661260343873178"}}`))
	}))
	defer func() { testServer.Close() }()

	h := abci.NewVoteExtHandler(input.Ctx.Logger(), input.OracleKeeper)
	handler := h.VerifyVoteExtensionHandler(abci.OracleConfig{
		Endpoint: testServer.URL,
	})

	voteExt := types.VoteExtension{
		Height: 3,
		Prices: []types.ExchangeRateTuple{
			{
				Denom:        "BTC",
				ExchangeRate: math.LegacyMustNewDecFromStr("47375.706652541026694000"),
			},
			{
				Denom:        "ETH",
				ExchangeRate: math.LegacyMustNewDecFromStr("2649.328939436595054949"),
			},
			{
				Denom:        "USDT",
				ExchangeRate: math.LegacyMustNewDecFromStr("1.000661260343873178"),
			},
		},
	}
	voteExtBz, err := voteExt.Marshal()
	require.NoError(t, err)
	// Height's same
	res, err := handler(input.Ctx, &cometabci.RequestVerifyVoteExtension{
		Hash:             []byte{},
		Height:           3,
		VoteExtension:    voteExtBz,
		ValidatorAddress: []byte{},
	})
	require.NoError(t, err)
	require.Equal(t, res.Status, cometabci.ResponseVerifyVoteExtension_ACCEPT)

	// Height different case
	_, err = handler(input.Ctx, &cometabci.RequestVerifyVoteExtension{
		Hash:             []byte{},
		Height:           2,
		VoteExtension:    voteExtBz,
		ValidatorAddress: []byte{},
	})
	require.Error(t, err)
}
