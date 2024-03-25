package types

import (
	errorsmod "cosmossdk.io/errors"
)

var (
	ErrIBCAccountAlreadyExist      = errorsmod.Register(ModuleName, 2, "interchain account already registered")
	ErrIBCAccountNotExist          = errorsmod.Register(ModuleName, 3, "interchain account does not exist")
	ErrLowerThanMinGasAmountPerAck = errorsmod.Register(ModuleName, 4, "lower than min gas amount per acknowledgement params")
)
