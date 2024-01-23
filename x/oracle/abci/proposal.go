package abci

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"

	"cosmossdk.io/core/appmodule"
	"cosmossdk.io/log"
	"cosmossdk.io/math"
	"github.com/Team-Kujira/core/x/oracle/keeper"
	"github.com/Team-Kujira/core/x/oracle/types"
	abci "github.com/cometbft/cometbft/abci/types"
	"github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/mempool"
	"github.com/cosmos/cosmos-sdk/types/module"
)

// StakeWeightedPrices defines the structure a proposer should use to calculate
// and submit the stake-weighted prices for a given set of supported currency
// pairs, in addition to the vote extensions used to calculate them. This is so
// validators can verify the proposer's calculations.
type StakeWeightedPrices struct {
	StakeWeightedPrices map[string]math.LegacyDec
	ExtendedCommitInfo  abci.ExtendedCommitInfo
	MissCounter         map[string]sdk.ValAddress
}

type ProposalHandler struct {
	logger        log.Logger
	keeper        keeper.Keeper
	valStore      baseapp.ValidatorStore
	mempool       mempool.Mempool
	txVerifier    baseapp.ProposalTxVerifier
	txSelector    baseapp.TxSelector
	ModuleManager *module.Manager
}

func NewProposalHandler(logger log.Logger, keeper keeper.Keeper, valStore baseapp.ValidatorStore, ModuleManager *module.Manager, mp mempool.Mempool, txVerifier baseapp.ProposalTxVerifier) *ProposalHandler {
	return &ProposalHandler{
		logger:        logger,
		keeper:        keeper,
		valStore:      valStore,
		ModuleManager: ModuleManager,
		mempool:       mp,
		txVerifier:    txVerifier,
		txSelector:    baseapp.NewDefaultTxSelector(),
	}
}

// cosmos-sdk/baseapp/abci_utils.go#L191
// PrepareProposalHandler returns the default implementation for processing an
// ABCI proposal. The application's mempool is enumerated and all valid
// transactions are added to the proposal. Transactions are valid if they:
//
// 1) Successfully encode to bytes.
// 2) Are valid (i.e. pass runTx, AnteHandler only).
//
// Enumeration is halted once RequestPrepareProposal.MaxBytes of transactions is
// reached or the mempool is exhausted.
//
// Note:
//
// - Step (2) is identical to the validation step performed in
// DefaultProcessProposal. It is very important that the same validation logic
// is used in both steps, and applications must ensure that this is the case in
// non-default handlers.
//
// - If no mempool is set or if the mempool is a no-op mempool, the transactions
// requested from CometBFT will simply be returned, which, by default, are in
// FIFO order.
func (h *ProposalHandler) PrepareProposal() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		var maxBlockGas uint64
		if b := ctx.ConsensusParams().Block; b != nil {
			maxBlockGas = uint64(b.MaxGas)
		}

		defer h.txSelector.Clear()

		proposalTxs := [][]byte{}

		// If the mempool is nil or NoOp we simply return the transactions
		// requested from CometBFT, which, by default, should be in FIFO order.
		//
		// Note, we still need to ensure the transactions returned respect req.MaxTxBytes.
		_, isNoOp := h.mempool.(mempool.NoOpMempool)
		if h.mempool == nil || isNoOp {
			for _, txBz := range req.Txs {
				tx, err := h.txVerifier.TxDecode(txBz)
				if err != nil {
					return nil, err
				}

				stop := h.txSelector.SelectTxForProposal(ctx, uint64(req.MaxTxBytes), maxBlockGas, tx, txBz)
				if stop {
					break
				}
			}

			proposalTxs = h.txSelector.SelectedTxs(ctx)
		} else {
			iterator := h.mempool.Select(ctx, req.Txs)
			for iterator != nil {
				memTx := iterator.Tx()

				// NOTE: Since transaction verification was already executed in CheckTx,
				// which calls mempool.Insert, in theory everything in the pool should be
				// valid. But some mempool implementations may insert invalid txs, so we
				// check again.
				txBz, err := h.txVerifier.PrepareProposalVerifyTx(memTx)
				if err != nil {
					err := h.mempool.Remove(memTx)
					if err != nil && !errors.Is(err, mempool.ErrTxNotFound) {
						return nil, err
					}
				} else {
					stop := h.txSelector.SelectTxForProposal(ctx, uint64(req.MaxTxBytes), maxBlockGas, memTx, txBz)
					if stop {
						break
					}
				}

				iterator = iterator.Next()
			}

			proposalTxs = h.txSelector.SelectedTxs(ctx)
		}

		err := baseapp.ValidateVoteExtensions(ctx, h.valStore, req.Height, ctx.ChainID(), req.LocalLastCommit)
		if err != nil {
			return nil, err
		}

		if req.Height >= ctx.ConsensusParams().Abci.VoteExtensionsEnableHeight {
			// stakeWeightedPrices, err := h.computeStakeWeightedOraclePrices(ctx, req.LocalLastCommit)
			// if err != nil {
			// 	return nil, errors.New("failed to compute stake-weighted oracle prices")
			// }

			stakeWeightedPrices, missMap, err := h.ComputeStakeWeightedPricesAndMissMap(ctx, req.LocalLastCommit)
			if err != nil {
				return nil, errors.New("failed to compute stake-weighted oracle prices")
			}

			injectedVoteExtTx := StakeWeightedPrices{
				StakeWeightedPrices: stakeWeightedPrices,
				ExtendedCommitInfo:  req.LocalLastCommit,
				MissCounter:         missMap,
			}

			// NOTE: We use stdlib JSON encoding, but an application may choose to use
			// a performant mechanism. This is for demo purposes only.
			bz, err := json.Marshal(injectedVoteExtTx)
			if err != nil {
				h.logger.Error("failed to encode injected vote extension tx", "err", err)
				return nil, errors.New("failed to encode injected vote extension tx")
			}

			// Inject a "fake" tx into the proposal s.t. validators can decode, verify,
			// and store the canonical stake-weighted average prices.
			proposalTxs = append(proposalTxs, bz)
		}

		// proceed with normal block proposal construction, e.g. POB, normal txs, etc...
		return &abci.ResponsePrepareProposal{
			Txs: proposalTxs,
		}, nil
	}
}

// cosmos-sdk/baseapp/abci_utils.go#L260
// ProcessProposalHandler returns the default implementation for processing an
// ABCI proposal. Every transaction in the proposal must pass 2 conditions:
//
// 1. The transaction bytes must decode to a valid transaction.
// 2. The transaction must be valid (i.e. pass runTx, AnteHandler only)
//
// If any transaction fails to pass either condition, the proposal is rejected.
// Note that step (2) is identical to the validation step performed in
// DefaultPrepareProposal. It is very important that the same validation logic
// is used in both steps, and applications must ensure that this is the case in
// non-default handlers.
func (h *ProposalHandler) ProcessProposal() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		var totalTxGas uint64

		var maxBlockGas int64
		if b := ctx.ConsensusParams().Block; b != nil {
			maxBlockGas = b.MaxGas
		}

		for _, txBytes := range req.Txs {
			var injectedVoteExtTx StakeWeightedPrices
			if err := json.Unmarshal(txBytes, &injectedVoteExtTx); err == nil {
				h.logger.Error("failed to decode injected vote extension tx", "err", err)
				err := baseapp.ValidateVoteExtensions(ctx, h.valStore, req.Height, ctx.ChainID(), injectedVoteExtTx.ExtendedCommitInfo)
				if err != nil {
					return nil, err
				}

				// Verify the proposer's stake-weighted oracle prices & miss counter by computing the same
				// calculation and comparing the results.
				// stakeWeightedPrices, err := h.computeStakeWeightedOraclePrices(ctx, injectedVoteExtTx.ExtendedCommitInfo)
				// if err != nil {
				// 	return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
				// }
				stakeWeightedPrices, missMap, err := h.ComputeStakeWeightedPricesAndMissMap(ctx, injectedVoteExtTx.ExtendedCommitInfo)
				if err != nil {
					return nil, errors.New("failed to compute stake-weighted oracle prices")
				}

				// compare stakeWeightedPrices
				if err := compareOraclePrices(injectedVoteExtTx.StakeWeightedPrices, stakeWeightedPrices); err != nil {
					return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
				}

				// compare missMap
				if err := compareMissMap(injectedVoteExtTx.MissCounter, missMap); err != nil {
					return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
				}
				continue
			}

			_, isNoOp := h.mempool.(mempool.NoOpMempool)
			if h.mempool == nil || isNoOp {
				continue
			}

			tx, err := h.txVerifier.ProcessProposalVerifyTx(txBytes)
			if err != nil {
				return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
			}

			if maxBlockGas > 0 {
				gasTx, ok := tx.(baseapp.GasTx)
				if ok {
					totalTxGas += gasTx.GetGas()
				}

				if totalTxGas > uint64(maxBlockGas) {
					return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
				}
			}
		}

		return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_ACCEPT}, nil
	}
}

// cosmos-sdk/types/module/module.go#L753
// PreBlock performs begin block functionality for upgrade module.
// It takes the current context as a parameter and returns a boolean value
// indicating whether the migration was successfully executed or not.
func (h *ProposalHandler) PreBlocker(ctx sdk.Context, req *abci.RequestFinalizeBlock) (*sdk.ResponsePreBlock, error) {
	ctx = ctx.WithEventManager(sdk.NewEventManager())
	paramsChanged := false
	for _, moduleName := range h.ModuleManager.OrderPreBlockers {
		if module, ok := h.ModuleManager.Modules[moduleName].(appmodule.HasPreBlocker); ok {
			rsp, err := module.PreBlock(ctx)
			if err != nil {
				return nil, err
			}
			if rsp.IsConsensusParamsChanged() {
				paramsChanged = true
			}
		}
	}

	for _, txBytes := range req.Txs {
		var injectedVoteExtTx StakeWeightedPrices
		if err := json.Unmarshal(txBytes, &injectedVoteExtTx); err != nil {
			h.logger.Error("failed to decode injected vote extension tx", "err", err)
			continue
		}

		for _, valAddr := range injectedVoteExtTx.MissCounter {
			h.keeper.SetMissCounter(ctx, valAddr, h.keeper.GetMissCounter(ctx, valAddr)+1)
		}

		// set oracle prices using the passed in context, which will make these prices available in the current block
		if err := h.keeper.SetOraclePrices(ctx, injectedVoteExtTx.StakeWeightedPrices); err != nil {
			return nil, err
		}
	}

	return &sdk.ResponsePreBlock{
		ConsensusParamsChanged: paramsChanged,
	}, nil
}

// func (h *ProposalHandler) computeStakeWeightedOraclePrices(ctx sdk.Context, ci abci.ExtendedCommitInfo) (map[string]math.LegacyDec, error) {
// 	stakeWeightedPrices := make(map[string]math.LegacyDec) // base -> average stake-weighted price

// 	var totalStake int64
// 	for _, v := range ci.Votes {
// 		if v.BlockIdFlag != cmtproto.BlockIDFlagCommit {
// 			continue
// 		}

// 		var voteExt OracleVoteExtension
// 		if err := json.Unmarshal(v.VoteExtension, &voteExt); err != nil {
// 			h.logger.Error("failed to decode vote extension", "err", err, "validator", fmt.Sprintf("%x", v.Validator.Address))
// 			return nil, err
// 		}

// 		totalStake += v.Validator.Power

// 		// Compute stake-weighted average of prices, i.e.
// 		// (P1)(W1) + (P2)(W2) + ... + (Pn)(Wn) / (W1 + W2 + ... + Wn)
// 		//
// 		// NOTE: These are the prices computed at the PREVIOUS height, i.e. H-1
// 		for base, price := range voteExt.Prices {
// 			if _, ok := stakeWeightedPrices[base]; !ok {
// 				stakeWeightedPrices[base] = math.LegacyZeroDec()
// 			}
// 			stakeWeightedPrices[base] = stakeWeightedPrices[base].Add(price.MulInt64(v.Validator.Power))
// 		}
// 	}

// 	if totalStake == 0 {
// 		return nil, nil
// 	}

// 	// finalize average by dividing by total stake, i.e. total weights
// 	for base, price := range stakeWeightedPrices {
// 		stakeWeightedPrices[base] = price.QuoInt64(totalStake)
// 	}

// 	return stakeWeightedPrices, nil
// }

func compareOraclePrices(p1, p2 map[string]math.LegacyDec) error {
	return nil
}

func compareMissMap(m1, m2 map[string]sdk.ValAddress) error {
	return nil
}

func (h *ProposalHandler) GetBallotByDenom(ctx sdk.Context, ci abci.ExtendedCommitInfo, validatorClaimMap map[string]types.Claim, validatorConsensusAddrMap map[string]sdk.ValAddress) (votes map[string]types.ExchangeRateBallot) {
	votes = map[string]types.ExchangeRateBallot{}

	for _, v := range ci.Votes {
		valAddr := validatorConsensusAddrMap[sdk.ConsAddress(v.Validator.Address).String()]
		claim, ok := validatorClaimMap[valAddr.String()]
		if ok {
			power := claim.Power

			var voteExt OracleVoteExtension
			if err := json.Unmarshal(v.VoteExtension, &voteExt); err != nil {
				h.logger.Error("failed to decode vote extension", "err", err, "validator", fmt.Sprintf("%x", v.Validator.Address))
				return votes
			}

			for base, price := range voteExt.Prices {
				tmpPower := power
				if !price.IsPositive() {
					// Make the power of abstain vote zero
					tmpPower = 0
				}

				votes[base] = append(votes[base],
					types.NewVoteForTally(
						price,
						base,
						v.Validator.Address,
						tmpPower,
					),
				)
			}
		}
	}

	// sort created ballot
	for denom, ballot := range votes {
		sort.Sort(ballot)
		votes[denom] = ballot
	}

	return votes
}

func (h *ProposalHandler) ComputeStakeWeightedPricesAndMissMap(ctx sdk.Context, ci abci.ExtendedCommitInfo) (map[string]math.LegacyDec, map[string]sdk.ValAddress, error) {
	params := h.keeper.GetParams(ctx)

	// Build claim map over all validators in active set
	stakeWeightedPrices := make(map[string]math.LegacyDec) // base -> average stake-weighted price
	validatorClaimMap := make(map[string]types.Claim)
	validatorConsensusAddrMap := make(map[string]sdk.ValAddress)

	maxValidators, err := h.keeper.StakingKeeper.MaxValidators(ctx)
	if err != nil {
		return nil, nil, err
	}
	iterator, err := h.keeper.StakingKeeper.ValidatorsPowerStoreIterator(ctx)
	if err != nil {
		return nil, nil, err
	}

	defer iterator.Close()

	powerReduction := h.keeper.StakingKeeper.PowerReduction(ctx)

	i := 0
	for ; iterator.Valid() && i < int(maxValidators); iterator.Next() {
		validator, err := h.keeper.StakingKeeper.Validator(ctx, iterator.Value())
		if err != nil {
			return nil, nil, err
		}

		// Exclude not bonded validator
		if validator.IsBonded() {
			valAddrStr := validator.GetOperator()
			valAddr, err := sdk.ValAddressFromBech32(valAddrStr)
			if err != nil {
				return nil, nil, err
			}

			validatorClaimMap[valAddr.String()] = types.NewClaim(validator.GetConsensusPower(powerReduction), 0, 0, valAddr)

			consAddr, err := validator.GetConsAddr()
			if err != nil {
				return nil, nil, err
			}
			validatorConsensusAddrMap[sdk.ConsAddress(consAddr).String()] = valAddr
			i++
		}
	}

	voteMap := h.GetBallotByDenom(ctx, ci, validatorClaimMap, validatorConsensusAddrMap)

	// Keep track, if a voter submitted a price deviating too much
	missMap := map[string]sdk.ValAddress{}

	// Iterate through ballots and update exchange rates; drop if not enough votes have been achieved.
	for denom, ballot := range voteMap {
		bondedTokens, err := h.keeper.StakingKeeper.TotalBondedTokens(ctx)
		if err != nil {
			return nil, nil, err
		}

		totalBondedPower := sdk.TokensToConsensusPower(bondedTokens, h.keeper.StakingKeeper.PowerReduction(ctx))
		voteThreshold := h.keeper.VoteThreshold(ctx)
		thresholdVotes := voteThreshold.MulInt64(totalBondedPower).RoundInt()
		ballotPower := math.NewInt(ballot.Power())

		if !ballotPower.IsZero() && ballotPower.GTE(thresholdVotes) {
			exchangeRate, err := keeper.Tally(
				ctx, ballot, params.MaxDeviation, validatorClaimMap, missMap,
			)
			if err != nil {
				return nil, nil, err
			}

			stakeWeightedPrices[denom] = exchangeRate
		}
	}

	//---------------------------
	// Do miss counting & slashing
	denomMap := map[string]map[string]struct{}{}
	var voteTargets []string
	voteTargets = append(voteTargets, params.RequiredDenoms...)

	for _, denom := range voteTargets {
		denomMap[denom] = map[string]struct{}{}
	}

	for denom, votes := range voteMap {
		for _, vote := range votes {
			// ignore denoms, not requested in voteTargets
			_, ok := denomMap[denom]
			if !ok {
				continue
			}

			denomMap[denom][vote.Voter.String()] = struct{}{}
		}
	}

	// Check if each validator is missing a required denom price
	for _, claim := range validatorClaimMap {
		for _, denom := range voteTargets {
			_, ok := denomMap[denom][claim.Recipient.String()]
			if !ok {
				missMap[claim.Recipient.String()] = claim.Recipient
				break
			}
		}
	}

	return stakeWeightedPrices, missMap, nil
}
