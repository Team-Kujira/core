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
	voteExt := types.VoteExtension2{}
	err = voteExt.Decompress(res.VoteExtension)
	require.NoError(t, err)
	require.Equal(t, voteExt.Height, int64(3))
	require.Equal(t, len(voteExt.Prices), 3)
	exchangeRates := make(map[uint32]string)
	for id, priceBz := range voteExt.Prices {
		exchangeRates[id] = abci.DecompressDecimal(priceBz).String()
	}
	require.Equal(t, exchangeRates[1], "47375.707000000000000000")
	require.Equal(t, exchangeRates[2], "2649.328900000000000000")
	require.Equal(t, exchangeRates[3], "1.000661300000000000")
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

	voteExt := types.VoteExtension2{
		Height: 3,
		Prices: map[uint32][]byte{
			1: abci.CompressDecimal(math.LegacyMustNewDecFromStr("47375.706652541026694000"), 8),
			2: abci.CompressDecimal(math.LegacyMustNewDecFromStr("2649.328939436595054949"), 8),
			3: abci.CompressDecimal(math.LegacyMustNewDecFromStr("1.000661260343873178"), 8),
		},
	}
	voteExtBz, err := voteExt.Compress()
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
