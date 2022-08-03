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
	CreateDenom *CreateDenom `json:"create_denom,omitempty"`
	/// Contracts can change the admin of a denom that they are the admin of.
	ChangeAdmin *ChangeAdmin `json:"change_admin,omitempty"`
	/// Contracts can mint native tokens for an existing factory denom
	/// that they are the admin of.
	MintTokens *MintTokens `json:"mint_tokens,omitempty"`
	/// Contracts can burn native tokens for an existing factory denom
	/// that they are the admin of.
	/// Currently, the burn from address must be the admin contract.
	BurnTokens *BurnTokens `json:"burn_tokens,omitempty"`
}

/// CreateDenom creates a new factory denom, of denomination:
/// factory/{creating contract address}/{Subdenom}
/// Subdenom can be of length at most 44 characters, in [0-9a-zA-Z./]
/// The (creating contract address, subdenom) pair must be unique.
/// The created denom's admin is the creating contract address,
/// but this admin can be changed using the ChangeAdmin binding.
type CreateDenom struct {
	Subdenom string `json:"subdenom"`
}

/// ChangeAdmin changes the admin for a factory denom.
/// If the NewAdminAddress is empty, the denom has no admin.
type ChangeAdmin struct {
	Denom           string `json:"denom"`
	NewAdminAddress string `json:"new_admin_address"`
}

type MintTokens struct {
	Denom         string  `json:"denom"`
	Amount        sdk.Int `json:"amount"`
	MintToAddress string  `json:"mint_to_address"`
}

type BurnTokens struct {
	Denom  string  `json:"denom"`
	Amount sdk.Int `json:"amount"`
	// BurnFromAddress must be set to "" for now.
	BurnFromAddress string `json:"burn_from_address"`
}

// createDenom creates a new token denom
func createDenom(ctx sdk.Context, contractAddr sdk.AccAddress, createDenom *CreateDenom, dk denomkeeper.Keeper, bk bankkeeper.BaseKeeper) ([]sdk.Event, [][]byte, error) {
	err := PerformCreateDenom(dk, bk, ctx, contractAddr, createDenom)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "perform create denom")
	}
	return nil, nil, nil
}

// PerformCreateDenom is used with createDenom to create a token denom; validates the msgCreateDenom.
func PerformCreateDenom(f denomkeeper.Keeper, b bankkeeper.BaseKeeper, ctx sdk.Context, contractAddr sdk.AccAddress, createDenom *CreateDenom) error {
	if createDenom == nil {
		return wasmvmtypes.InvalidRequest{Err: "create denom null create denom"}
	}

	msgServer := denomkeeper.NewMsgServerImpl(f)

	msgCreateDenom := denomtypes.NewMsgCreateDenom(contractAddr.String(), createDenom.Subdenom)

	if err := msgCreateDenom.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "failed validating MsgCreateDenom")
	}

	// Create denom
	_, err := msgServer.CreateDenom(
		sdk.WrapSDKContext(ctx),
		msgCreateDenom,
	)
	if err != nil {
		return sdkerrors.Wrap(err, "creating denom")
	}
	return nil
}

// mintTokens mints tokens of a specified denom to an address.
func mintTokens(ctx sdk.Context, contractAddr sdk.AccAddress, mint *MintTokens, dk denomkeeper.Keeper, bk bankkeeper.BaseKeeper) ([]sdk.Event, [][]byte, error) {
	err := PerformMint(dk, bk, ctx, contractAddr, mint)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "perform mint")
	}
	return nil, nil, nil
}

// PerformMint used with mintTokens to validate the mint message and mint through token factory.
func PerformMint(f denomkeeper.Keeper, b bankkeeper.BaseKeeper, ctx sdk.Context, contractAddr sdk.AccAddress, mint *MintTokens) error {
	if mint == nil {
		return wasmvmtypes.InvalidRequest{Err: "mint token null mint"}
	}
	rcpt, err := parseAddress(mint.MintToAddress)
	if err != nil {
		return err
	}

	coin := sdk.Coin{Denom: mint.Denom, Amount: mint.Amount}
	sdkMsg := denomtypes.NewMsgMint(contractAddr.String(), coin)
	if err = sdkMsg.ValidateBasic(); err != nil {
		return err
	}

	// Mint through token factory / message server
	msgServer := denomkeeper.NewMsgServerImpl(f)
	_, err = msgServer.Mint(sdk.WrapSDKContext(ctx), sdkMsg)
	if err != nil {
		return sdkerrors.Wrap(err, "minting coins from message")
	}
	err = b.SendCoins(ctx, contractAddr, rcpt, sdk.NewCoins(coin))
	if err != nil {
		return sdkerrors.Wrap(err, "sending newly minted coins from message")
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
	newAdminAddr, err := parseAddress(changeAdmin.NewAdminAddress)
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

// burnTokens burns tokens.
func burnTokens(ctx sdk.Context, contractAddr sdk.AccAddress, burn *BurnTokens, dk denomkeeper.Keeper) ([]sdk.Event, [][]byte, error) {
	err := PerformBurn(dk, ctx, contractAddr, burn)
	if err != nil {
		return nil, nil, sdkerrors.Wrap(err, "perform burn")
	}
	return nil, nil, nil
}

// PerformBurn performs token burning after validating tokenBurn message.
func PerformBurn(f denomkeeper.Keeper, ctx sdk.Context, contractAddr sdk.AccAddress, burn *BurnTokens) error {
	if burn == nil {
		return wasmvmtypes.InvalidRequest{Err: "burn token null mint"}
	}
	if burn.BurnFromAddress != "" && burn.BurnFromAddress != contractAddr.String() {
		return wasmvmtypes.InvalidRequest{Err: "BurnFromAddress must be \"\""}
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
	if q.CreateDenom != nil {
		return createDenom(ctx, contractAddr, q.CreateDenom, dk, bk)
	}
	if q.MintTokens != nil {
		return mintTokens(ctx, contractAddr, q.MintTokens, dk, bk)
	}
	if q.ChangeAdmin != nil {
		return changeAdmin(ctx, contractAddr, q.ChangeAdmin, dk)
	}
	if q.BurnTokens != nil {
		return burnTokens(ctx, contractAddr, q.BurnTokens, dk)
	}

	return nil, nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Custom variant"}
}
