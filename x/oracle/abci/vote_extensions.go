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

// OracleVoteExtension defines the canonical vote extension structure.
type OracleVoteExtension struct {
	Height int64
	Prices map[string]math.LegacyDec
}

type PricesResponse struct {
	Prices map[string]math.LegacyDec `json:"prices"`
}

func (h *VoteExtHandler) ExtendVoteHandler(oracleConfig OracleConfig) sdk.ExtendVoteHandler {
	return func(ctx sdk.Context, req *abci.RequestExtendVote) (*abci.ResponseExtendVote, error) {
		h.currentBlock = req.Height
		h.lastPriceSyncTS = time.Now()

		h.logger.Info("computing oracle prices for vote extension", "height", req.Height, "time", h.lastPriceSyncTS, "endpoint", oracleConfig.Endpoint)

		res, err := http.Get(oracleConfig.Endpoint)
		if err != nil {
			return nil, err
		}
		defer res.Body.Close()

		resBody, err := io.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}

		prices := PricesResponse{}
		err = json.Unmarshal(resBody, &prices)
		if err != nil {
			return nil, err
		}

		computedPrices := prices.Prices

		// produce a canonical vote extension
		voteExt := OracleVoteExtension{
			Height: req.Height,
			Prices: computedPrices,
		}

		h.logger.Info("computed prices", "prices", computedPrices)

		// Encode vote extension to bytes
		bz, err := json.Marshal(voteExt)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal vote extension: %w", err)
		}

		return &abci.ResponseExtendVote{VoteExtension: bz}, nil
	}
}

func (h *VoteExtHandler) VerifyVoteExtensionHandler(_ OracleConfig) sdk.VerifyVoteExtensionHandler {
	return func(ctx sdk.Context, req *abci.RequestVerifyVoteExtension) (*abci.ResponseVerifyVoteExtension, error) {
		var voteExt OracleVoteExtension

		if len(req.VoteExtension) == 0 {
			return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
		}

		err := json.Unmarshal(req.VoteExtension, &voteExt)
		if err != nil {
			// NOTE: It is safe to return an error as the Cosmos SDK will capture all
			// errors, log them, and reject the proposal.
			return nil, fmt.Errorf("failed to unmarshal vote extension: %w", err)
		}

		if voteExt.Height != req.Height {
			return nil, fmt.Errorf("vote extension height does not match request height; expected: %d, got: %d", req.Height, voteExt.Height)
		}

		// Verify incoming prices from a validator are valid. Note, verification during
		// VerifyVoteExtensionHandler MUST be deterministic. For brevity and demo
		// purposes, we omit implementation.
		if err := h.verifyOraclePrices(ctx, voteExt.Prices); err != nil {
			return nil, fmt.Errorf("failed to verify oracle prices from validator %X: %w", req.ValidatorAddress, err)
		}

		return &abci.ResponseVerifyVoteExtension{Status: abci.ResponseVerifyVoteExtension_ACCEPT}, nil
	}
}

func (h *VoteExtHandler) verifyOraclePrices(_ sdk.Context, _ map[string]math.LegacyDec) error {
	return nil
}
