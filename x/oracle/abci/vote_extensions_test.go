package abci_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Team-Kujira/core/x/oracle/abci"
	"github.com/Team-Kujira/core/x/oracle/keeper"
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
	require.Equal(t, string(res.VoteExtension), `{"Height":3,"Prices":{"BTC":"47375.706652541026694000","ETH":"2649.328939436595054949","USDT":"1.000661260343873178"}}`)
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

	// Height's same
	res, err := handler(input.Ctx, &cometabci.RequestVerifyVoteExtension{
		Hash:             []byte{},
		Height:           3,
		VoteExtension:    []byte(`{"Height":3,"Prices":{"BTC":"47375.706652541026694000","ETH":"2649.328939436595054949","USDT":"1.000661260343873178"}}`),
		ValidatorAddress: []byte{},
	})
	require.NoError(t, err)
	require.Equal(t, res.Status, cometabci.ResponseVerifyVoteExtension_ACCEPT)

	// Height different case
	_, err = handler(input.Ctx, &cometabci.RequestVerifyVoteExtension{
		Hash:             []byte{},
		Height:           2,
		VoteExtension:    []byte(`{"Height":3,"Prices":{"BTC":"47375.706652541026694000","ETH":"2649.328939436595054949","USDT":"1.000661260343873178"}}`),
		ValidatorAddress: []byte{},
	})
	require.Error(t, err)
}
