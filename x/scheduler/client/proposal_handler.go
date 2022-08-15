package client

import (
	"github.com/Team-Kujira/core/x/scheduler/client/cli"
	"github.com/Team-Kujira/core/x/scheduler/client/rest"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// ProposalHandlers define the wasm cli proposal types and rest handler.
var CreateHookProposalHandler = govclient.NewProposalHandler(cli.CreateHookProposalCmd, rest.CreateHookProposalHandler)
var UpdateHookProposalHandler = govclient.NewProposalHandler(cli.UpdateHookProposalCmd, rest.UpdateHookProposalHandler)
var DeleteHookProposalHandler = govclient.NewProposalHandler(cli.DeleteHookProposalCmd, rest.DeleteHookProposalHandler)
