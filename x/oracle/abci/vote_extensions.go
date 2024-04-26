package abci

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/Team-Kujira/core/x/oracle/keeper"
	"github.com/Team-Kujira/core/x/oracle/types"
	abci "github.com/cometbft/cometbft/abci/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type VoteExtHandler struct {
	logger          log.Logger
	currentBlock    int64     // current block height
	lastPriceSyncTS time.Time // last time we synced prices

	Keeper keeper.Keeper
}

func NewVoteExtHandler(
	logger log.Logger,
	keeper keeper.Keeper,
) *VoteExtHandler {
	return &VoteExtHandler{
		logger: logger,
		Keeper: keeper,
	}
}

type PricesResponse struct {
	Prices map[string]math.LegacyDec `json:"prices"`
}

func (h *VoteExtHandler) ExtendVoteHandler(oracleConfig OracleConfig) sdk.ExtendVoteHandler {
	return func(_ sdk.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
		h.currentBlock = req.Height
		h.lastPriceSyncTS = time.Now()

		h.logger.Info("computing oracle prices for vote extension", "height", req.Height, "time", h.lastPriceSyncTS, "endpoint", oracleConfig.Endpoint)

		emptyVoteExt := types.VoteExtension{
			Height: req.Height,
			Prices: []types.ExchangeRateTuple{},
		}

		// Encode vote extension to bytes
		emptyVoteExtBz, err := emptyVoteExt.Marshal()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal vote extension: %w", err)
		}

		res, err := http.Get(oracleConfig.Endpoint)
		if err != nil {
			h.logger.Error("failed to query endpoint", err)
			return &abci.ResponseExtendVote{VoteExtension: emptyVoteExtBz}, nil
		}
		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			h.logger.Error("failed to read response body", err)
			return &abci.ResponseExtendVote{VoteExtension: emptyVoteExtBz}, nil
		}

		prices := PricesResponse{}
		err = json.Unmarshal(resBody, &prices)
		if err != nil {
			h.logger.Error("failed to unmarshal prices", err)
			return &abci.ResponseExtendVote{VoteExtension: emptyVoteExtBz}, nil
		}

		computedPrices := []types.ExchangeRateTuple{}
		for denom, rate := range prices.Prices {
			computedPrices = append(computedPrices, types.ExchangeRateTuple{
				Denom:        denom,
				ExchangeRate: rate,
			})
		}

		// produce a canonical vote extension
		voteExt := types.VoteExtension{
			Height: req.Height,
			Prices: computedPrices,
		}

		h.logger.Info("computed prices", "prices", computedPrices)

		// Encode vote extension to bytes
		bz, err := voteExt.Marshal()
		if err != nil {
			return nil, fmt.Errorf("failed to marshal vote extension: %w", err)
		}

		return &abci.ResponseExtendVote{VoteExtension: bz}, nil
	}
}

func (h *VoteExtHandler) VerifyVoteExtensionHandler(_ OracleConfig) sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
		var voteExt types.VoteExtension

		if len(req.VoteExtension) == 0 {
			return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
		}

		err := voteExt.Unmarshal(req.VoteExtension)
		if err != nil {
			// NOTE: It is safe to return an error as the Cosmos SDK will capture all
			// errors, log them, and reject the proposal.
			return nil, fmt.Errorf("failed to unmarshal vote extension: %w", err)
		}

		if voteExt.Height != req.Height {
			return nil, fmt.Errorf("vote extension height does not match request height; expected: %d, got: %d", req.Height, voteExt.Height)
		}

		// Verify incoming prices from a validator are valid. Note, verification during
		// VerifyVoteExtensionHandler MUST be deterministic.
		if err := h.verifyOraclePrices(ctx, voteExt.Prices); err != nil {
			return nil, fmt.Errorf("failed to verify oracle prices from validator %X: %w", req.ValidatorAddress, err)
		}

		return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}

func (h *VoteExtHandler) verifyOraclePrices(_ sdk.Context, _ []types.ExchangeRateTuple) error {
	return nil
}
