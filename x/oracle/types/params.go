package types

import (
	"fmt"

	"cosmossdk.io/math"
	"gopkg.in/yaml.v2"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyVoteThreshold     = []byte("VoteThreshold")
	KeyMaxDeviation      = []byte("MaxDeviation")
	KeyRequiredSymbols   = []byte("RequiredSymbols")
	KeySlashFraction     = []byte("SlashFraction")
	KeySlashWindow       = []byte("SlashWindow")
	KeyMinValidPerWindow = []byte("MinValidPerWindow")
	// Deprecated
	KeyRewardBand = []byte("RewardBand")
	KeyWhitelist  = []byte("Whitelist")
)

// Default parameter values
const (
	DefaultSlashWindow              = uint64(274000)   // window for a week
	DefaultRewardDistributionWindow = uint64(14250000) // window for a year
)

// Default parameter values
var (
	DefaultVoteThreshold     = math.LegacyNewDecWithPrec(5, 1) // 0.5   (50%)
	DefaultMaxDeviation      = math.LegacyNewDecWithPrec(2, 2) // 0.02  ( 2%)
	DefaultRequiredSymbols   = []Symbol{}
	DefaultSlashFraction     = math.LegacyNewDecWithPrec(1, 4) // 0.0001 (0.01%)
	DefaultMinValidPerWindow = math.LegacyNewDecWithPrec(5, 2) // 0.05   (5%)
)

var _ paramstypes.ParamSet = &Params{}

// DefaultParams creates default oracle module parameters
func DefaultParams() Params {
	return Params{
		VoteThreshold:     DefaultVoteThreshold,
		MaxDeviation:      DefaultMaxDeviation,
		RequiredSymbols:   DefaultRequiredSymbols,
		SlashFraction:     DefaultSlashFraction,
		SlashWindow:       DefaultSlashWindow,
		MinValidPerWindow: DefaultMinValidPerWindow,
	}
}

// ParamKeyTable returns the parameter key table.
func ParamKeyTable() paramstypes.KeyTable {
	return paramstypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of oracle module's parameters.
func (p *Params) ParamSetPairs() paramstypes.ParamSetPairs {
	rewardBand := math.LegacyDec{}
	whitelist := DenomList{}
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyVoteThreshold, &p.VoteThreshold, validateVoteThreshold),
		paramstypes.NewParamSetPair(KeyMaxDeviation, &p.MaxDeviation, validateMaxDeviation),
		paramstypes.NewParamSetPair(KeyRequiredSymbols, &p.RequiredSymbols, validateRequiredSymbols),
		paramstypes.NewParamSetPair(KeySlashFraction, &p.SlashFraction, validateSlashFraction),
		paramstypes.NewParamSetPair(KeySlashWindow, &p.SlashWindow, validateSlashWindow),
		paramstypes.NewParamSetPair(KeyMinValidPerWindow, &p.MinValidPerWindow, validateMinValidPerWindow),
		paramstypes.NewParamSetPair(KeyRewardBand, &rewardBand, validateRewardBand),
		paramstypes.NewParamSetPair(KeyWhitelist, &whitelist, validateWhitelist),
	}
}

// String implements fmt.Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate performs basic validation on oracle parameters.
func (p Params) Validate() error {
	if p.VoteThreshold.LTE(math.LegacyNewDecWithPrec(33, 2)) {
		return fmt.Errorf("oracle parameter VoteThreshold must be greater than 33 percent")
	}

	if p.MaxDeviation.GT(math.LegacyOneDec()) || p.MaxDeviation.IsNegative() {
		return fmt.Errorf("oracle parameter MaxDeviation must be between [0, 1]")
	}

	if p.SlashFraction.GT(math.LegacyOneDec()) || p.SlashFraction.IsNegative() {
		return fmt.Errorf("oracle parameter SlashFraction must be between [0, 1]")
	}

	if p.SlashWindow == 0 {
		return fmt.Errorf("oracle parameter SlashWindow must be positive")
	}

	if p.MinValidPerWindow.GT(math.LegacyOneDec()) || p.MinValidPerWindow.IsNegative() {
		return fmt.Errorf("oracle parameter MinValidPerWindow must be between [0, 1]")
	}

	if err := validateRequiredSymbols(p.RequiredSymbols); err != nil {
		return err
	}

	return nil
}

func validateVoteThreshold(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.LT(math.LegacyNewDecWithPrec(33, 2)) {
		return fmt.Errorf("vote threshold must be bigger than 33%%: %s", v)
	}

	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("vote threshold too large: %s", v)
	}

	return nil
}

func validateMaxDeviation(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("reward band must be positive: %s", v)
	}

	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("reward band is too large: %s", v)
	}

	return nil
}

func validateRequiredSymbols(i interface{}) error {
	v, ok := i.([]Symbol)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, d := range v {
		if len(d.Symbol) == 0 {
			return fmt.Errorf("oracle parameter RequiredSymbols Denom must not be ''")
		}
	}

	registeredSymbols := make(map[string]bool)
	registeredIDs := make(map[uint32]bool)
	for _, denom := range v {
		if registeredSymbols[denom.Symbol] {
			return fmt.Errorf("oracle parameter denom should be unique")
		}
		if registeredIDs[denom.Id] {
			return fmt.Errorf("oracle parameter denom id should be unique")
		}
		registeredSymbols[denom.Symbol] = true
		registeredIDs[denom.Id] = true
	}

	return nil
}

func validateSlashFraction(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("slash fraction must be positive: %s", v)
	}

	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("slash fraction is too large: %s", v)
	}

	return nil
}

func validateSlashWindow(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("slash window must be positive: %d", v)
	}

	return nil
}

func validateMinValidPerWindow(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("min valid per window must be positive: %s", v)
	}

	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("min valid per window is too large: %s", v)
	}

	return nil
}

func validateRewardBand(i interface{}) error {
	v, ok := i.(math.LegacyDec)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNegative() {
		return fmt.Errorf("reward band must be positive: %s", v)
	}

	if v.GT(math.LegacyOneDec()) {
		return fmt.Errorf("reward band is too large: %s", v)
	}

	return nil
}

func validateWhitelist(i interface{}) error {
	v, ok := i.(DenomList)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, d := range v {
		if len(d.Name) == 0 {
			return fmt.Errorf("oracle parameter Whitelist Denom must have name")
		}
	}

	return nil
}
