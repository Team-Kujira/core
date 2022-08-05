package wasm

import (
	denomkeeper "kujira/x/denom/keeper"

	denomtypes "kujira/x/denom/types"

	wasmvmtypes "github.com/CosmWasm/wasmvm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	bankkeeper "github.com/cosmos/cosmos-sdk/x/bank/keeper"
)

type DenomMsg struct {
	/// Contracts can create denoms, namespaced under the contract's address.
	/// A contract may create any number of independent sub-denoms.
	Create *Create `json:"create,omitempty"`
	/// Contracts can change the admin of a denom that they are the admin of.
	ChangeAdmin *ChangeAdmin `json:"change_admin,omitempty"`
	/// Contracts can mint native tokens for an existing factory denom
	/// that they are the admin of.
	Mint *Mint `json:"mint,omitempty"`
	/// Contracts can burn native tokens for an existing factory denom
	/// that they are the admin of.
	/// Currently, the burn from address must be the admin contract.
	Burn *Burn `json:"burn,omitempty"`
}

/// Create creates a new factory denom, of denomination:
/// factory/{creating contract address}/{Subdenom}
/// Subdenom can be of length at most 44 characters, in [0-9a-zA-Z./]
/// The (creating contract address, subdenom) pair must be unique.
/// The created denom's admin is the creating contract address,
/// but this admin can be changed using the ChangeAdmin binding.
type Create struct {
	Subdenom string `json:"subdenom"`
}

/// ChangeAdmin changes the admin for a factory denom.
/// If the Address is empty, the denom has no admin.
type ChangeAdmin struct {
	Denom   string `json:"denom"`
	Address string `json:"address"`
}

type Mint struct {
	Denom     string  `json:"denom"`
	Amount    sdk.Int `json:"amount"`
	Recipient string  `json:"recipient"`
}

type Burn struct {
	Denom  string  `json:"denom"`
	Amount sdk.Int `json:"amount"`
}

// create creates a new token denom
func create(ctx sdk.Context, contractAddr sdk.AccAddress, create *Create, dk denomkeeper.Keeper, bk bankkeeper.BaseKeeper) ([]sdk.Event, [][]byte, error) {
	err := PerformCreate(dk, bk, ctx, contractAddr, create)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "perform create denom")
	}
	return nil, nil, nil
}

// PerformCreate is used with create to create a token denom; validates the msgCreate.
func PerformCreate(f denomkeeper.Keeper, b bankkeeper.BaseKeeper, ctx sdk.Context, contractAddr sdk.AccAddress, create *Create) error {
	if create == nil {
		return wasmvmtypes.InvalidRequest{Err: "create denom null create denom"}
	}

	msgServer := denomkeeper.NewMsgServerImpl(f)

	msgCreate := denomtypes.NewMsgCreateDenom(contractAddr.String(), create.Subdenom)

	if err := msgCreate.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "failed validating MsgCreate")
	}

	// Create denom
	_, err := msgServer.CreateDenom(
		sdk.WrapSDKContext(ctx),
		msgCreate,
	)
	if err != nil {
		return sdkerrors.Wrap(err, "creating denom")
	}
	return nil
}

// mint mints tokens of a specified denom to an address.
func mint(ctx sdk.Context, contractAddr sdk.AccAddress, mint *Mint, dk denomkeeper.Keeper, bk bankkeeper.BaseKeeper) ([]sdk.Event, [][]byte, error) {
	err := PerformMint(dk, bk, ctx, contractAddr, mint)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "perform mint")
	}
	return nil, nil, nil
}

// PerformMint used with mint to validate the mint message and mint through token factory.
func PerformMint(f denomkeeper.Keeper, b bankkeeper.BaseKeeper, ctx sdk.Context, contractAddr sdk.AccAddress, mint *Mint) error {
	if mint == nil {
		return wasmvmtypes.InvalidRequest{Err: "mint token null mint"}
	}
	_, err := parseAddress(mint.Recipient)
	if err != nil {
		return err
	}

	coin := sdk.Coin{Denom: mint.Denom, Amount: mint.Amount}
	sdkMsg := denomtypes.NewMsgMint(contractAddr.String(), coin, mint.Recipient)
	if err = sdkMsg.ValidateBasic(); err != nil {
		return err
	}

	// Mint through token factory / message server
	msgServer := denomkeeper.NewMsgServerImpl(f)
	_, err = msgServer.Mint(sdk.WrapSDKContext(ctx), sdkMsg)
	if err != nil {
		return sdkerrors.Wrap(err, "minting coins from message")
	}
	return nil
}

// changeAdmin changes the admin.
func changeAdmin(ctx sdk.Context, contractAddr sdk.AccAddress, changeAdmin *ChangeAdmin, dk denomkeeper.Keeper) ([]sdk.Event, [][]byte, error) {
	err := PerformChangeAdmin(dk, ctx, contractAddr, changeAdmin)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "failed to change admin")
	}
	return nil, nil, nil
}

// ChangeAdmin is used with changeAdmin to validate changeAdmin messages and to dispatch.
func PerformChangeAdmin(f denomkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, changeAdmin *ChangeAdmin) error {
	if changeAdmin == nil {
		return wasmvmtypes.InvalidRequest{Err: "changeAdmin is nil"}
	}
	newAdminAddr, err := parseAddress(changeAdmin.Address)
	if err != nil {
		return err
	}

	changeAdminMsg := denomtypes.NewMsgChangeAdmin(contractAddr.String(), changeAdmin.Denom, newAdminAddr.String())
	if err := changeAdminMsg.ValidateBasic(); err != nil {
		return err
	}

	msgServer := denomkeeper.NewMsgServerImpl(f)
	_, err = msgServer.ChangeAdmin(sdk.WrapSDKContext(ctx), changeAdminMsg)
	if err != nil {
		return sdkerrors.Wrap(err, "failed changing admin from message")
	}
	return nil
}

// burn burns tokens.
func burn(ctx sdk.Context, contractAddr sdk.AccAddress, burn *Burn, dk denomkeeper.Keeper) ([]sdk.Event, [][]byte, error) {
	err := PerformBurn(dk, ctx, contractAddr, burn)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "perform burn")
	}
	return nil, nil, nil
}

// PerformBurn performs token burning after validating tokenBurn message.
func PerformBurn(f denomkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, burn *Burn) error {
	if burn == nil {
		return wasmvmtypes.InvalidRequest{Err: "burn token null mint"}
	}

	coin := sdk.Coin{Denom: burn.Denom, Amount: burn.Amount}
	sdkMsg := denomtypes.NewMsgBurn(contractAddr.String(), coin)
	if err := sdkMsg.ValidateBasic(); err != nil {
		return err
	}

	// Burn through token factory / message server
	msgServer := denomkeeper.NewMsgServerImpl(f)
	_, err := msgServer.Burn(sdk.WrapSDKContext(ctx), sdkMsg)
	if err != nil {
		return sdkerrors.Wrap(err, "burning coins from message")
	}
	return nil
}

// QueryCustom implements custom query interface
func HandleMsg(dk denomkeeper.Keeper, bk bankkeeper.BaseKeeper, contractAddr sdk.AccAddress, ctx sdk.Context, q *DenomMsg) ([]sdk.Event, [][]byte, error) {
	if q.Create != nil {
		return create(ctx, contractAddr, q.Create, dk, bk)
	}
	if q.Mint != nil {
		return mint(ctx, contractAddr, q.Mint, dk, bk)
	}
	if q.ChangeAdmin != nil {
		return changeAdmin(ctx, contractAddr, q.ChangeAdmin, dk)
	}
	if q.Burn != nil {
		return burn(ctx, contractAddr, q.Burn, dk)
	}

	return nil, nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Custom variant"}
}
