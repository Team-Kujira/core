package types

import (
	"fmt"

	"cosmossdk.io/math"
	"gopkg.in/yaml.v2"

	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter keys
var (
	KeyVotePeriod        = []byte("VotePeriod")
	KeyVoteThreshold     = []byte("VoteThreshold")
	KeyMaxDeviation      = []byte("MaxDeviation")
	KeyRequiredDenoms    = []byte("RequiredDenoms")
	KeySlashFraction     = []byte("SlashFraction")
	KeySlashWindow       = []byte("SlashWindow")
	KeyMinValidPerWindow = []byte("MinValidPerWindow")
	// Deprecated
	KeyRewardBand = []byte("RewardBand")
	KeyWhitelist  = []byte("Whitelist")
)

// Default parameter values
const (
	DefaultVotePeriod               = uint64(14)       // 30 seconds
	DefaultSlashWindow              = uint64(274000)   // window for a week
	DefaultRewardDistributionWindow = uint64(14250000) // window for a year
)

// Default parameter values
var (
	DefaultVoteThreshold     = math.LegacyNewDecWithPrec(50, 2) // 50%
	DefaultMaxDeviation      = math.LegacyNewDecWithPrec(2, 1)  // 2% (-1, 1)
	DefaultRequiredDenoms    = []string{}
	DefaultSlashFraction     = math.LegacyNewDecWithPrec(1, 4) // 0.01%
	DefaultMinValidPerWindow = math.LegacyNewDecWithPrec(5, 2) // 5%
)

var _ paramstypes.ParamSet = &Params{}

// DefaultParams creates default oracle module parameters
func DefaultParams() Params {
	return Params{
		VotePeriod:        DefaultVotePeriod,
		VoteThreshold:     DefaultVoteThreshold,
		MaxDeviation:      DefaultMaxDeviation,
		RequiredDenoms:    DefaultRequiredDenoms,
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
	return paramstypes.ParamSetPairs{
		paramstypes.NewParamSetPair(KeyVotePeriod, &p.VotePeriod, validateVotePeriod),
		paramstypes.NewParamSetPair(KeyVoteThreshold, &p.VoteThreshold, validateVoteThreshold),
		paramstypes.NewParamSetPair(KeyMaxDeviation, &p.MaxDeviation, validateMaxDeviation),
		paramstypes.NewParamSetPair(KeyRequiredDenoms, &p.RequiredDenoms, validateRequiredDenoms),
		paramstypes.NewParamSetPair(KeySlashFraction, &p.SlashFraction, validateSlashFraction),
		paramstypes.NewParamSetPair(KeySlashWindow, &p.SlashWindow, validateSlashWindow),
		paramstypes.NewParamSetPair(KeyMinValidPerWindow, &p.MinValidPerWindow, validateMinValidPerWindow),
		paramstypes.NewParamSetPair(KeyRewardBand, &p.RewardBand, func(_ interface{}) error { return nil }),
		paramstypes.NewParamSetPair(KeyWhitelist, &p.Whitelist, func(_ interface{}) error { return nil }),
	}
}

// String implements fmt.Stringer interface
func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// Validate performs basic validation on oracle parameters.
func (p Params) Validate() error {
	if p.VotePeriod == 0 {
		return fmt.Errorf("oracle parameter VotePeriod must be > 0, is %d", p.VotePeriod)
	}
	if p.VoteThreshold.LTE(math.LegacyNewDecWithPrec(33, 2)) {
		return fmt.Errorf("oracle parameter VoteThreshold must be greater than 33 percent")
	}

	if p.MaxDeviation.GT(math.LegacyOneDec()) || p.MaxDeviation.IsNegative() {
		return fmt.Errorf("oracle parameter MaxDeviation must be between [0, 1]")
	}

	if p.SlashFraction.GT(math.LegacyOneDec()) || p.SlashFraction.IsNegative() {
		return fmt.Errorf("oracle parameter SlashFraction must be between [0, 1]")
	}

	if p.SlashWindow < p.VotePeriod {
		return fmt.Errorf("oracle parameter SlashWindow must be greater than or equal with VotePeriod")
	}

	if p.MinValidPerWindow.GT(math.LegacyOneDec()) || p.MinValidPerWindow.IsNegative() {
		return fmt.Errorf("oracle parameter MinValidPerWindow must be between [0, 1]")
	}

	for _, denom := range p.RequiredDenoms {
		if len(denom) == 0 {
			return fmt.Errorf("oracle parameter RequiredDenoms Denom must not be ''")
		}
	}
	return nil
}

func validateVotePeriod(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("vote period must be positive: %d", v)
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

func validateRequiredDenoms(i interface{}) error {
	v, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	for _, d := range v {
		if len(d) == 0 {
			return fmt.Errorf("oracle parameter RequiredDenoms Denom must not be ''")
		}
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
