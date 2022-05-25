package rest

import (
	"encoding/json"
	"kujira/x/scheduler/types"
	"net/http"

	wasmtypes "github.com/CosmWasm/wasmd/x/wasm/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	govrest "github.com/cosmos/cosmos-sdk/x/gov/client/rest"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
)

type CreateHookProposalJSONReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Proposer string    `json:"proposer" yaml:"proposer"`
	Deposit  sdk.Coins `json:"deposit" yaml:"deposit"`

	Contract string          `json:"contract" yaml:"contract"`
	Msg      json.RawMessage `json:"msg" yaml:"msg"`
	// Executor is the role that is passed to the contract's environment
	Executor string `json:"executor" yaml:"executor"`
	// Call this message every N blocks
	Frequency int64     `json:"frequency" yaml:"frequency"`
	Funds     sdk.Coins `json:"funds" yaml:"funds"`
}

func (s CreateHookProposalJSONReq) Content() govtypes.Content {
	return &types.CreateHookProposal{
		Title:       s.Title,
		Description: s.Description,
		Contract:    s.Contract,
		Msg:         wasmtypes.RawContractMessage(s.Msg),
		Executor:    s.Executor,
		Frequency:   s.Frequency,
		Funds:       s.Funds,
	}
}

func (s CreateHookProposalJSONReq) GetProposer() string {
	return s.Proposer
}

func (s CreateHookProposalJSONReq) GetDeposit() sdk.Coins {
	return s.Deposit
}

func (s CreateHookProposalJSONReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}

func CreateHookProposalHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "scheduler_create",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req CreateHookProposalJSONReq
			if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type UpdateHookProposalJSONReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Proposer string    `json:"proposer" yaml:"proposer"`
	Deposit  sdk.Coins `json:"deposit" yaml:"deposit"`

	Id       uint64          `json:"id" yaml:"id"`
	Contract string          `json:"contract" yaml:"contract"`
	Msg      json.RawMessage `json:"msg" yaml:"msg"`
	// Executor is the role that is passed to the contract's environment
	Executor string `json:"executor" yaml:"executor"`
	// Call this message every N blocks
	Frequency int64     `json:"frequency" yaml:"frequency"`
	Funds     sdk.Coins `json:"funds" yaml:"funds"`
}

func (s UpdateHookProposalJSONReq) Content() govtypes.Content {
	return &types.UpdateHookProposal{
		Title:       s.Title,
		Description: s.Description,
		Id:          s.Id,
		Contract:    s.Contract,
		Msg:         wasmtypes.RawContractMessage(s.Msg),
		Executor:    s.Executor,
		Frequency:   s.Frequency,
		Funds:       s.Funds,
	}
}

func (s UpdateHookProposalJSONReq) GetProposer() string {
	return s.Proposer
}

func (s UpdateHookProposalJSONReq) GetDeposit() sdk.Coins {
	return s.Deposit
}

func (s UpdateHookProposalJSONReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}

func UpdateHookProposalHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "scheduler_update",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req UpdateHookProposalJSONReq
			if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type DeleteHookProposalJSONReq struct {
	BaseReq rest.BaseReq `json:"base_req" yaml:"base_req"`

	Title       string `json:"title" yaml:"title"`
	Description string `json:"description" yaml:"description"`

	Proposer string    `json:"proposer" yaml:"proposer"`
	Deposit  sdk.Coins `json:"deposit" yaml:"deposit"`

	Id uint64 `json:"id" yaml:"id"`
}

func (s DeleteHookProposalJSONReq) Content() govtypes.Content {
	return &types.DeleteHookProposal{
		Title:       s.Title,
		Description: s.Description,
		Id:          s.Id,
	}
}

func (s DeleteHookProposalJSONReq) GetProposer() string {
	return s.Proposer
}

func (s DeleteHookProposalJSONReq) GetDeposit() sdk.Coins {
	return s.Deposit
}

func (s DeleteHookProposalJSONReq) GetBaseReq() rest.BaseReq {
	return s.BaseReq
}

func DeleteHookProposalHandler(cliCtx client.Context) govrest.ProposalRESTHandler {
	return govrest.ProposalRESTHandler{
		SubRoute: "scheduler_delete",
		Handler: func(w http.ResponseWriter, r *http.Request) {
			var req DeleteHookProposalJSONReq
			if !rest.ReadRESTReq(w, r, cliCtx.LegacyAmino, &req) {
				return
			}
			toStdTxResponse(cliCtx, w, req)
		},
	}
}

type wasmProposalData interface {
	Content() govtypes.Content
	GetProposer() string
	GetDeposit() sdk.Coins
	GetBaseReq() rest.BaseReq
}

func toStdTxResponse(cliCtx client.Context, w http.ResponseWriter, data wasmProposalData) {
	proposerAddr, err := sdk.AccAddressFromBech32(data.GetProposer())
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	msg, err := govtypes.NewMsgSubmitProposal(data.Content(), data.GetDeposit(), proposerAddr)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := msg.ValidateBasic(); err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}
	baseReq := data.GetBaseReq().Sanitize()
	if !baseReq.ValidateBasic(w) {
		return
	}
	tx.WriteGeneratedTxResponse(cliCtx, w, baseReq, msg)
}
