package bindings

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Messages
type AllianceMsg struct {
	// Contracts can delegate from the contract.
	Delegate *Delegate `json:"delegate,omitempty"`
	// Contracts can redelegate from the contract.
	Redelegate *Redelegate `json:"redelegate,omitempty"`
	// Contracts can undelegate from the contract.
	Undelegate *Undelegate `json:"undelegate,omitempty"`
	// Contracts can claim delegation rewards from the contract.
	ClaimDelegationRewards *ClaimDelegationRewards `json:"claim_delegation_rewards,omitempty"`
}

type Delegate struct {
	ValidatorAddress string
	Amount           sdk.Coin
}

type Redelegate struct {
	ValidatorSrcAddress string
	ValidatorDstAddress string
	Amount              sdk.Coin
}

type Undelegate struct {
	ValidatorAddress string
	Amount           sdk.Coin
}

type ClaimDelegationRewards struct {
	ValidatorAddress string
	Denom            string
}

// Queries
type AllianceQuery struct {
	Params      *Params      `json:"params,omitempty"`
	Alliance    *Alliance    `json:"alliance,omitempty"`
	IBCAlliance *IBCAlliance `json:"ibc_alliance,omitempty"`
}

type Params struct{}

type Alliance struct {
	Denom string
}

type IBCAlliance struct {
	Hash string
}
