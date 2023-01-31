package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	govtypesv1beta "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

type ProposalType string

const (
	ProposalTypeCreateHook ProposalType = "CreateHook"
	ProposalTypeUpdateHook ProposalType = "UpdateHook"
	ProposalTypeDeleteHook ProposalType = "DeleteHook"
)

func init() { // register new content types with the sdk
	govtypesv1beta.RegisterProposalType(string(ProposalTypeCreateHook))
	govtypesv1beta.RegisterProposalType(string(ProposalTypeUpdateHook))
	govtypesv1beta.RegisterProposalType(string(ProposalTypeDeleteHook))
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p CreateHookProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type
func (p CreateHookProposal) ProposalType() string { return string(ProposalTypeCreateHook) }

// ValidateBasic validates the proposal
func (p CreateHookProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	if _, err := sdk.AccAddressFromBech32(p.Contract); err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	if _, err := sdk.AccAddressFromBech32(p.Executor); err != nil {
		return sdkerrors.Wrap(err, "executor")
	}
	if !p.Funds.IsValid() {
		return sdkerrors.ErrInvalidCoins
	}
	if err := p.Msg.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "payload msg")
	}
	return nil
}

// String implements the Stringer interface.
func (p CreateHookProposal) String() string {
	return fmt.Sprintf(`Create Hook Proposal:
  Title:       %s
  Description: %s
  Contract:    %s
  Executor:    %s
  Msg:         %q
  Funds:       %s
`, p.Title, p.Description, p.Contract, p.Executor, p.Msg, p.Funds)
}

// MarshalYAML pretty prints the wasm byte code
func (p CreateHookProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Title       string    `yaml:"title"`
		Description string    `yaml:"description"`
		Contract    string    `yaml:"contract"`
		Executor    string    `yaml:"executor"`
		Msg         string    `yaml:"msg"`
		Funds       sdk.Coins `yaml:"funds"`
	}{
		Title:       p.Title,
		Description: p.Description,
		Contract:    p.Contract,
		Executor:    p.Executor,
		Msg:         string(p.Msg),
		Funds:       p.Funds,
	}, nil
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p UpdateHookProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type
func (p UpdateHookProposal) ProposalType() string {
	return string(ProposalTypeUpdateHook)
}

// ValidateBasic validates the proposal
func (p UpdateHookProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}

	if p.Id == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "ID is required")
	}

	if _, err := sdk.AccAddressFromBech32(p.Contract); err != nil {
		return sdkerrors.Wrap(err, "contract")
	}
	if _, err := sdk.AccAddressFromBech32(p.Executor); err != nil {
		return sdkerrors.Wrap(err, "executor")
	}
	if !p.Funds.IsValid() {
		return sdkerrors.ErrInvalidCoins
	}
	if err := p.Msg.ValidateBasic(); err != nil {
		return sdkerrors.Wrap(err, "payload msg")
	}
	return nil
}

// String implements the Stringer interface.
func (p UpdateHookProposal) String() string {
	return fmt.Sprintf(`Update Hook Proposal:
  Title:       %s
  Description: %s
  ID:          %d
  Contract:    %s
  Executor:    %s
  Msg:         %q
  Funds:       %s
  `, p.Title, p.Description, p.Id, p.Contract, p.Executor, p.Msg, p.Funds)
}

// MarshalYAML pretty prints the init message
//
//nolint:revive
func (p UpdateHookProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Id          uint64    `yaml:"id"` //nolint:stylecheck
		Title       string    `yaml:"title"`
		Description string    `yaml:"description"`
		Contract    string    `yaml:"contract"`
		Executor    string    `yaml:"executor"`
		Msg         string    `yaml:"msg"`
		Funds       sdk.Coins `yaml:"funds"`
	}{
		Id:          p.Id,
		Title:       p.Title,
		Description: p.Description,
		Contract:    p.Contract,
		Executor:    p.Executor,
		Msg:         string(p.Msg),
		Funds:       p.Funds,
	}, nil
}

// ProposalRoute returns the routing key of a parameter change proposal.
func (p DeleteHookProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type
func (p DeleteHookProposal) ProposalType() string { return string(ProposalTypeDeleteHook) }

// ValidateBasic validates the proposal
func (p DeleteHookProposal) ValidateBasic() error {
	if err := validateProposalCommons(p.Title, p.Description); err != nil {
		return err
	}
	return nil
}

// String implements the Stringer interface.
func (p DeleteHookProposal) String() string {
	return fmt.Sprintf(`Migrate Contract Proposal:
  Title:       %s
  Description: %s
  ID:    %d
`, p.Title, p.Description, p.Id)
}

// MarshalYAML pretty prints the migrate message
func (p DeleteHookProposal) MarshalYAML() (interface{}, error) {
	return struct {
		Title       string `yaml:"title"`
		Description string `yaml:"description"`
		Id          uint64 `yaml:"id"` //nolint:revive,stylecheck
	}{
		Title:       p.Title,
		Description: p.Description,
		Id:          p.Id,
	}, nil
}

func validateProposalCommons(title, description string) error {
	if strings.TrimSpace(title) != title {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal title must not start/end with white spaces")
	}
	if len(title) == 0 {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal title cannot be blank")
	}
	if len(title) > govtypesv1beta.MaxTitleLength {
		return sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal title is longer than max length of %d", govtypesv1beta.MaxTitleLength)
	}
	if strings.TrimSpace(description) != description {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal description must not start/end with white spaces")
	}
	if len(description) == 0 {
		return sdkerrors.Wrap(govtypes.ErrInvalidProposalContent, "proposal description cannot be blank")
	}
	if len(description) > govtypesv1beta.MaxDescriptionLength {
		return sdkerrors.Wrapf(govtypes.ErrInvalidProposalContent, "proposal description is longer than max length of %d", govtypesv1beta.MaxDescriptionLength)
	}
	return nil
}
