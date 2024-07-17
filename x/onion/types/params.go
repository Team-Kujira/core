package types

import (
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	_ paramtypes.ParamSet = &Params{}
)

func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable()
}

func NewParams() Params {
	return Params{}
}

// DefaultParams returns default concentrated-liquidity module parameters.
func DefaultParams() Params {
	return Params{}
}

// ParamSetPairs implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{}
}

// Validate params.
func (p Params) Validate() error {
	return nil
}
