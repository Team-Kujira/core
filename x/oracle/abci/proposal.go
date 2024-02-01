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
	logger   log.Logger
	keeper   keeper.Keeper
	valStore baseapp.ValidatorStore
	baseapp.DefaultProposalHandler
	ModuleManager *module.Manager
}

func NewProposalHandler(logger log.Logger, keeper keeper.Keeper, valStore baseapp.ValidatorStore, ModuleManager *module.Manager, mp mempool.Mempool, txVerifier baseapp.ProposalTxVerifier) *ProposalHandler {
	return &ProposalHandler{
		logger:                 logger,
		keeper:                 keeper,
		valStore:               valStore,
		ModuleManager:          ModuleManager,
		DefaultProposalHandler: *baseapp.NewDefaultProposalHandler(mp, txVerifier),
	}
}

// PrepareProposalHandler returns the implementation for processing an
// ABCI proposal.
// - Default PrepareProposalHandler selects regular txs
// - Appends vote extension tx at the end
func (h *ProposalHandler) PrepareProposal() sdk.PrepareProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestPrepareProposal) (*abci.ResponsePrepareProposal, error) {
		defaultHandler := h.DefaultProposalHandler.PrepareProposalHandler()
		defaultResponse, err := defaultHandler(ctx, req)
		if err != nil {
			return nil, err
		}

		proposalTxs := defaultResponse.Txs

		// Note: Upgrade height should be equal to vote extension enable height
		if req.Height >= ctx.ConsensusParams().Abci.VoteExtensionsEnableHeight {
			err = baseapp.ValidateVoteExtensions(ctx, h.valStore, req.Height, ctx.ChainID(), req.LocalLastCommit)
			if err != nil {
				return nil, err
			}

			stakeWeightedPrices, missMap, err := h.ComputeStakeWeightedPricesAndMissMap(ctx, req.LocalLastCommit)
			if err != nil {
				return nil, errors.New("failed to compute stake-weighted oracle prices")
			}

			injectedVoteExtTx := StakeWeightedPrices{
				StakeWeightedPrices: stakeWeightedPrices,
				ExtendedCommitInfo:  req.LocalLastCommit,
				MissCounter:         missMap,
			}

			// Encode vote extension to bytes
			bz, err := json.Marshal(injectedVoteExtTx)
			if err != nil {
				h.logger.Error("failed to encode injected vote extension tx", "err", err)
				return nil, errors.New("failed to encode injected vote extension tx")
			}

			// Inject vote extension tx into the proposal s.t. validators can decode, verify,
			// and store the canonical stake-weighted average prices.
			proposalTxs = append(proposalTxs, bz)
		}

		// proceed with normal block proposal construction, e.g. POB, normal txs, etc...
		return &abci.ResponsePrepareProposal{
			Txs: proposalTxs,
		}, nil
	}
}

// ProcessProposalHandler returns the implementation for processing an
// ABCI proposal
// - Validate vote extension tx
// - Validate regular tx with default PrepareProposalHandler
func (h *ProposalHandler) ProcessProposal() sdk.ProcessProposalHandler {
	return func(ctx sdk.Context, req *abci.RequestProcessProposal) (*abci.ResponseProcessProposal, error) {
		reReq := *req
		var injectedVoteExtTx StakeWeightedPrices
		if len(req.Txs) > 0 {
			lastTx := req.Txs[len(req.Txs)-1]
			if err := json.Unmarshal(lastTx, &injectedVoteExtTx); err == nil {
				h.logger.Debug("handling injected vote extension tx")
				err := baseapp.ValidateVoteExtensions(ctx, h.valStore, req.Height, ctx.ChainID(), injectedVoteExtTx.ExtendedCommitInfo)
				if err != nil {
					return nil, err
				}

				// Verify the proposer's stake-weighted oracle prices & miss counter by computing the same
				// calculation and comparing the results.
				stakeWeightedPrices, missMap, err := h.ComputeStakeWeightedPricesAndMissMap(ctx, injectedVoteExtTx.ExtendedCommitInfo)
				if err != nil {
					return nil, errors.New("failed to compute stake-weighted oracle prices")
				}

				// compare stakeWeightedPrices
				if err := CompareOraclePrices(injectedVoteExtTx.StakeWeightedPrices, stakeWeightedPrices); err != nil {
					return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
				}

				// compare missMap
				if err := CompareMissMap(injectedVoteExtTx.MissCounter, missMap); err != nil {
					return &abci.ResponseProcessProposal{Status: abci.ResponseProcessProposal_REJECT}, nil
				}

				// Exclude last tx if it's vote extension tx
				reReq.Txs = reReq.Txs[:len(reReq.Txs)-1]
			}
		}

		defaultHandler := h.DefaultProposalHandler.ProcessProposalHandler()
		return defaultHandler(ctx, &reReq)
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
			h.logger.Debug("Skipping regular tx, not one of our vote", "tx", string(txBytes))
			continue
		}

		// Clear all exchange rates
		h.keeper.IterateExchangeRates(ctx, func(denom string, _ math.LegacyDec) (stop bool) {
			h.keeper.DeleteExchangeRate(ctx, denom)
			return false
		})

		for _, valAddr := range injectedVoteExtTx.MissCounter {
			h.keeper.SetMissCounter(ctx, valAddr, h.keeper.GetMissCounter(ctx, valAddr)+1)
		}

		// set oracle prices using the passed in context, which will make these prices available in the current block
		if err := h.keeper.SetOraclePrices(ctx, injectedVoteExtTx.StakeWeightedPrices); err != nil {
			return nil, err
		}

		// Do slash who did miss voting over threshold and
		// reset miss counters of all validators at the last block of slash window
		params := h.keeper.GetParams(ctx)
		if IsPeriodLastBlock(ctx, params.SlashWindow) {
			h.keeper.SlashAndResetMissCounters(ctx)
		}
	}

	return &sdk.ResponsePreBlock{
		ConsensusParamsChanged: paramsChanged,
	}, nil
}

func IsPeriodLastBlock(ctx sdk.Context, blocksPerPeriod uint64) bool {
	return (uint64(ctx.BlockHeight())+1)%blocksPerPeriod == 0
}

func CompareOraclePrices(p1, p2 map[string]math.LegacyDec) error {
	for denom, p := range p1 {
		if p2[denom].IsNil() || !p.Equal(p2[denom]) {
			return errors.New("oracle prices mismatch")
		}
	}

	for denom, p := range p2 {
		if p1[denom].IsNil() || !p.Equal(p1[denom]) {
			return errors.New("oracle prices mismatch")
		}
	}

	return nil
}

func CompareMissMap(m1, m2 map[string]sdk.ValAddress) error {
	for valAddrStr, valAddr := range m1 {
		if _, ok := m2[valAddrStr]; !ok {
			return errors.New("oracle missMap mismatch")
		}
		if valAddr.String() != valAddrStr {
			return errors.New("invalid oracle missMap")
		}
	}

	for valAddrStr, valAddr := range m2 {
		if _, ok := m1[valAddrStr]; !ok {
			return errors.New("oracle missMap mismatch")
		}
		if valAddr.String() != valAddrStr {
			return errors.New("invalid oracle missMap")
		}
	}

	return nil
}

func (h *ProposalHandler) GetBallotByDenom(ci abci.ExtendedCommitInfo, validatorClaimMap map[string]types.Claim, validatorConsensusAddrMap map[string]sdk.ValAddress) (votes map[string]types.ExchangeRateBallot) {
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
						valAddr,
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

	voteMap := h.GetBallotByDenom(ci, validatorClaimMap, validatorConsensusAddrMap)

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
