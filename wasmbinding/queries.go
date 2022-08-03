package wasmbinding

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"kujira/wasmbinding/bindings"
	denomkeeper "kujira/x/denom/keeper"

	oraclekeeper "kujira/x/oracle/keeper"

	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

type QueryPlugin struct {
	denomKeeper  denomkeeper.Keeper
	bankkeeper   bankkeeper.Keeper
	oraclekeeper oraclekeeper.Keeper
}

// NewQueryPlugin returns a reference to a new QueryPlugin.
func NewQueryPlugin(bk bankkeeper.Keeper, ok oraclekeeper.Keeper, dk denomkeeper.Keeper) *QueryPlugin {
	return &QueryPlugin{
		denomKeeper:  dk,
		bankkeeper:   bk,
		oraclekeeper: ok,
	}
}

// GetDenomAdmin is a query to get denom admin.
func (qp QueryPlugin) GetDenomAdmin(ctx sdk.Context, denom string) (*bindings.DenomAdminResponse, error) {
	metadata, err := qp.denomKeeper.GetAuthorityMetadata(ctx, denom)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin for denom: %s", denom)
	}

	return &bindings.DenomAdminResponse{Admin: metadata.Admin}, nil
}
