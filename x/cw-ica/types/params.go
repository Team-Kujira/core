package types

import (
	fmt "fmt"
)

// cw-ica params default values
const (
	DefaultMinGasAmountPerAck uint64 = 200_000 // 200k
)

// NewParams creates a new Params instance
func NewParams(minGasAmountPerAck uint64) Params {
	return Params{
		MinGasAmountPerAck: minGasAmountPerAck,
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return NewParams(
		DefaultMinGasAmountPerAck,
	)
}

// validate a set of params
func (p Params) Validate() error {
	if err := validateMinGasAmountPerAck(p.MinGasAmountPerAck); err != nil {
		return err
	}

	return nil
}

func validateMinGasAmountPerAck(i interface{}) error {
	v, ok := i.(uint64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v <= 0 {
		return fmt.Errorf("MinGasAmountPerAck must be positive: %d", v)
	}

	return nil
}
