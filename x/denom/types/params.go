package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Parameter store keys.
var (
	KeyCreationFee   = []byte("CreationFee")
	KeyNoFeeAccounts = []byte("NoFeeAccounts")
)

// ParamTable for gamm module.
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

func NewParams(creationFee sdk.Coins) Params {
	return Params{
		CreationFee: creationFee,
	}
}

// default gamm module parameters.
func DefaultParams() Params {
	return Params{
		CreationFee: sdk.NewCoins(sdk.NewInt64Coin("ukuji", 10_000_000)),
	}
}

// validate params.
func (p Params) Validate() error {
	err := validateCreationFee(p.CreationFee)
	if err != nil {
		return err
	}
	return nil
}

// Implements params.ParamSet.
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyCreationFee, &p.CreationFee, validateCreationFee),
	}
}

func validateCreationFee(i interface{}) error {
	v, ok := i.(sdk.Coins)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.Validate() != nil {
		return fmt.Errorf("invalid denom creation fee: %+v", i)
	}

	return nil
}
