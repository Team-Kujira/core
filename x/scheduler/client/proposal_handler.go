package client

import (
	"github.com/Team-Kujira/core/x/scheduler/client/cli"

	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"
)

// ProposalHandlers define the wasm cli proposal types and rest handler.
var (
	CreateHookProposalHandler = govclient.NewProposalHandler(cli.CreateHookProposalCmd)
	UpdateHookProposalHandler = govclient.NewProposalHandler(cli.UpdateHookProposalCmd)
	DeleteHookProposalHandler = govclient.NewProposalHandler(cli.DeleteHookProposalCmd)
)
