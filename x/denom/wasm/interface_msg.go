package wasm

import (
	"cosmossdk.io/errors"
	"cosmossdk.io/math"
	wasmvmtypes "github.com/CosmWasm/wasmvm/v2/types"
	denomkeeper "github.com/Team-Kujira/core/x/denom/keeper"
	denomtypes "github.com/Team-Kujira/core/x/denom/types"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/gogoproto/proto"

	// bankkeeper "github.com/terra-money/alliance/custom/bank/keeper"
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

// / Create creates a new factory denom, of denomination:
// / factory/{creating contract address}/{Subdenom}
// / Subdenom can be of length at most 44 characters, in [0-9a-zA-Z./]
// / The (creating contract address, subdenom) pair must be unique.
// / The created denom's admin is the creating contract address,
// / but this admin can be changed using the ChangeAdmin binding.
type Create struct {
	Subdenom string `json:"subdenom"`
}

// / ChangeAdmin changes the admin for a factory denom.
// / If the Address is empty, the denom has no admin.
type ChangeAdmin struct {
	Denom   string `json:"denom"`
	Address string `json:"address"`
}

type Mint struct {
	Denom     string   `json:"denom"`
	Amount    math.Int `json:"amount"`
	Recipient string   `json:"recipient"`
}

type Burn struct {
	Denom  string   `json:"denom"`
	Amount math.Int `json:"amount"`
}

// create creates a new token denom
func create(
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	create *Create,
	dk denomkeeper.Keeper,
	bk bankkeeper.Keeper,
) (*denomtypes.MsgCreateDenomResponse, error) {
	res, err := PerformCreate(dk, bk, ctx, contractAddr, create)
	if err != nil {
		return nil, errors.Wrap(err, "perform create denom")
	}
	return res, nil
}

// PerformCreate is used with create to create a token denom; validates the msgCreate.
func PerformCreate(
	f denomkeeper.Keeper,
	_ bankkeeper.Keeper,
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	create *Create,
) (*denomtypes.MsgCreateDenomResponse, error) {
	if create == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "create denom null create denom"}
	}

	msgServer := denomkeeper.NewMsgServerImpl(f)

	msgCreate := denomtypes.NewMsgCreateDenom(contractAddr.String(), create.Subdenom)

	if err := msgCreate.ValidateBasic(); err != nil {
		return nil, errors.Wrap(err, "failed validating MsgCreate")
	}

	// Create denom
	res, err := msgServer.CreateDenom(
		ctx,
		msgCreate,
	)
	if err != nil {
		return nil, errors.Wrap(err, "creating denom")
	}
	return res, nil
}

// mint mints tokens of a specified denom to an address.
func mint(
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	mint *Mint,
	dk denomkeeper.Keeper,
	bk bankkeeper.Keeper,
) (*denomtypes.MsgMintResponse, error) {
	res, err := PerformMint(dk, bk, ctx, contractAddr, mint)
	if err != nil {
		return nil, errors.Wrap(err, "perform mint")
	}
	return res, nil
}

// PerformMint used with mint to validate the mint message and mint through token factory.
func PerformMint(
	f denomkeeper.Keeper,
	_ bankkeeper.Keeper,
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	mint *Mint,
) (*denomtypes.MsgMintResponse, error) {
	if mint == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "mint token null mint"}
	}
	_, err := parseAddress(mint.Recipient)
	if err != nil {
		return nil, err
	}

	coin := sdk.Coin{Denom: mint.Denom, Amount: mint.Amount}
	sdkMsg := denomtypes.NewMsgMint(contractAddr.String(), coin, mint.Recipient)
	if err = sdkMsg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Mint through token factory / message server
	msgServer := denomkeeper.NewMsgServerImpl(f)
	res, err := msgServer.Mint(ctx, sdkMsg)
	if err != nil {
		return nil, errors.Wrap(err, "minting coins from message")
	}
	return res, nil
}

// changeAdmin changes the admin.
func changeAdmin(
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	changeAdmin *ChangeAdmin,
	dk denomkeeper.Keeper,
) (*denomtypes.MsgChangeAdminResponse, error) {
	res, err := PerformChangeAdmin(dk, ctx, contractAddr, changeAdmin)
	if err != nil {
		return nil, errors.Wrap(err, "failed to change admin")
	}
	return res, nil
}

// ChangeAdmin is used with changeAdmin to validate changeAdmin messages and to dispatch.
func PerformChangeAdmin(
	f denomkeeper.Keeper,
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	changeAdmin *ChangeAdmin,
) (*denomtypes.MsgChangeAdminResponse, error) {
	if changeAdmin == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "changeAdmin is nil"}
	}
	newAdminAddr, err := parseAddress(changeAdmin.Address)
	if err != nil {
		return nil, err
	}

	changeAdminMsg := denomtypes.NewMsgChangeAdmin(contractAddr.String(), changeAdmin.Denom, newAdminAddr.String())
	if err := changeAdminMsg.ValidateBasic(); err != nil {
		return nil, err
	}

	msgServer := denomkeeper.NewMsgServerImpl(f)
	res, err := msgServer.ChangeAdmin(ctx, changeAdminMsg)
	if err != nil {
		return nil, errors.Wrap(err, "failed changing admin from message")
	}
	return res, nil
}

// burn burns tokens.
func burn(
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	burn *Burn,
	dk denomkeeper.Keeper,
) (*denomtypes.MsgBurnResponse, error) {
	res, err := PerformBurn(dk, ctx, contractAddr, burn)
	if err != nil {
		return nil, errors.Wrap(err, "perform burn")
	}

	return res, nil
}

// PerformBurn performs token burning after validating tokenBurn message.
func PerformBurn(
	f denomkeeper.Keeper,
	ctx sdk.Context,
	contractAddr sdk.AccAddress,
	burn *Burn,
) (*denomtypes.MsgBurnResponse, error) {
	if burn == nil {
		return nil, wasmvmtypes.InvalidRequest{Err: "burn token null mint"}
	}

	coin := sdk.Coin{Denom: burn.Denom, Amount: burn.Amount}
	sdkMsg := denomtypes.NewMsgBurn(contractAddr.String(), coin)
	if err := sdkMsg.ValidateBasic(); err != nil {
		return nil, err
	}

	// Burn through token factory / message server
	msgServer := denomkeeper.NewMsgServerImpl(f)
	res, err := msgServer.Burn(ctx, sdkMsg)
	if err != nil {
		return nil, errors.Wrap(err, "burning coins from message")
	}

	return res, nil
}

// QueryCustom implements custom query interface
func HandleMsg(
	dk denomkeeper.Keeper,
	bk bankkeeper.Keeper,
	contractAddr sdk.AccAddress,
	ctx sdk.Context,
	q *DenomMsg,
) ([]sdk.Event, [][]byte, [][]*codectypes.Any, error) {
	var res proto.Message
	var err error

	if q.Create != nil {
		res, err = create(ctx, contractAddr, q.Create, dk, bk)
	}
	if q.Mint != nil {
		res, err = mint(ctx, contractAddr, q.Mint, dk, bk)
	}
	if q.ChangeAdmin != nil {
		res, err = changeAdmin(ctx, contractAddr, q.ChangeAdmin, dk)
	}
	if q.Burn != nil {
		res, err = burn(ctx, contractAddr, q.Burn, dk)
	}
	if err != nil {
		return nil, nil, nil, err
	}
	if res == nil {
		return nil, nil, nil, wasmvmtypes.UnsupportedRequest{Kind: "unknown Custom variant"}
	}

	x, err := codectypes.NewAnyWithValue(res)
	if err != nil {
		return nil, nil, nil, err
	}
	msgResponses := [][]*codectypes.Any{{x}}

	return nil, nil, msgResponses, err
}
