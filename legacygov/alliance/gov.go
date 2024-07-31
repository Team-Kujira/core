package alliance

import (
	"time"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const RouterKey = "alliance"

const (
	ProposalTypeCreateAlliance = "msg_create_alliance_proposal"
	ProposalTypeUpdateAlliance = "msg_update_alliance_proposal"
	ProposalTypeDeleteAlliance = "msg_delete_alliance_proposal"
)

var (
	_ govtypes.Content = &MsgCreateAllianceProposal{}
	_ govtypes.Content = &MsgUpdateAllianceProposal{}
	_ govtypes.Content = &MsgDeleteAllianceProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeCreateAlliance)
	govtypes.RegisterProposalType(ProposalTypeUpdateAlliance)
	govtypes.RegisterProposalType(ProposalTypeDeleteAlliance)
}

func NewMsgCreateAllianceProposal(title, description, denom string, rewardWeight math.LegacyDec, rewardWeightRange RewardWeightRange, takeRate math.LegacyDec, rewardChangeRate math.LegacyDec, rewardChangeInterval time.Duration) govtypes.Content {
	return &MsgCreateAllianceProposal{
		Title:                title,
		Description:          description,
		Denom:                denom,
		RewardWeight:         rewardWeight,
		RewardWeightRange:    rewardWeightRange,
		TakeRate:             takeRate,
		RewardChangeRate:     rewardChangeRate,
		RewardChangeInterval: rewardChangeInterval,
	}
}
func (m *MsgCreateAllianceProposal) GetTitle() string       { return m.Title }
func (m *MsgCreateAllianceProposal) GetDescription() string { return m.Description }
func (m *MsgCreateAllianceProposal) ProposalRoute() string  { return RouterKey }
func (m *MsgCreateAllianceProposal) ProposalType() string   { return ProposalTypeCreateAlliance }

func (m *MsgCreateAllianceProposal) ValidateBasic() error {
	if m.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	if err := sdk.ValidateDenom(m.Denom); err != nil {
		return err
	}

	if m.RewardWeight.IsNil() || m.RewardWeight.LT(math.LegacyZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be zero or a positive number")
	}

	if m.RewardWeightRange.Min.IsNil() || m.RewardWeightRange.Min.LT(math.LegacyZeroDec()) ||
		m.RewardWeightRange.Max.IsNil() || m.RewardWeightRange.Max.LT(math.LegacyZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight min and max must be zero or a positive number")
	}

	if m.RewardWeightRange.Min.GT(m.RewardWeightRange.Max) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight min must be less or equal to rewardWeight max")
	}

	if m.RewardWeight.LT(m.RewardWeightRange.Min) || m.RewardWeight.GT(m.RewardWeightRange.Max) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be bounded in RewardWeightRange")
	}

	if m.TakeRate.IsNil() || m.TakeRate.IsNegative() || m.TakeRate.GTE(math.LegacyOneDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance takeRate must be more or equals to 0 but strictly less than 1")
	}

	if m.RewardChangeRate.IsZero() || m.RewardChangeRate.IsNegative() {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardChangeRate must be strictly a positive number")
	}

	if m.RewardChangeInterval < 0 {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardChangeInterval must be strictly a positive number")
	}

	return nil
}

func NewMsgUpdateAllianceProposal(title, description, denom string, rewardWeight math.LegacyDec, rewardWeightRange RewardWeightRange, takeRate math.LegacyDec, rewardChangeRate math.LegacyDec, rewardChangeInterval time.Duration) govtypes.Content {
	return &MsgUpdateAllianceProposal{
		Title:                title,
		Description:          description,
		Denom:                denom,
		RewardWeight:         rewardWeight,
		TakeRate:             takeRate,
		RewardChangeRate:     rewardChangeRate,
		RewardChangeInterval: rewardChangeInterval,
		RewardWeightRange:    rewardWeightRange,
	}
}
func (m *MsgUpdateAllianceProposal) GetTitle() string       { return m.Title }
func (m *MsgUpdateAllianceProposal) GetDescription() string { return m.Description }
func (m *MsgUpdateAllianceProposal) ProposalRoute() string  { return RouterKey }
func (m *MsgUpdateAllianceProposal) ProposalType() string   { return ProposalTypeUpdateAlliance }

func (m *MsgUpdateAllianceProposal) ValidateBasic() error {
	if m.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}

	if m.RewardWeight.IsNil() || m.RewardWeight.LT(math.LegacyZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be zero or a positive number")
	}

	if m.TakeRate.IsNil() || m.TakeRate.IsNegative() || m.TakeRate.GTE(math.LegacyOneDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance takeRate must be more or equals to 0 but strictly less than 1")
	}

	if m.RewardChangeRate.IsZero() || m.RewardChangeRate.IsNegative() {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardChangeRate must be strictly a positive number")
	}

	if m.RewardChangeInterval < 0 {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardChangeInterval must be strictly a positive number")
	}

	if m.RewardWeightRange.Min.IsNil() || m.RewardWeightRange.Min.LT(math.LegacyZeroDec()) ||
		m.RewardWeightRange.Max.IsNil() || m.RewardWeightRange.Max.LT(math.LegacyZeroDec()) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight min and max must be zero or a positive number")
	}

	if m.RewardWeightRange.Min.GT(m.RewardWeightRange.Max) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight min must be less or equal to rewardWeight max")
	}

	if m.RewardWeight.LT(m.RewardWeightRange.Min) || m.RewardWeight.GT(m.RewardWeightRange.Max) {
		return status.Errorf(codes.InvalidArgument, "Alliance rewardWeight must be bounded in RewardWeightRange")
	}

	return nil
}

func NewMsgDeleteAllianceProposal(title, description, denom string) govtypes.Content {
	return &MsgDeleteAllianceProposal{
		Title:       title,
		Description: description,
		Denom:       denom,
	}
}
func (m *MsgDeleteAllianceProposal) GetTitle() string       { return m.Title }
func (m *MsgDeleteAllianceProposal) GetDescription() string { return m.Description }
func (m *MsgDeleteAllianceProposal) ProposalRoute() string  { return RouterKey }
func (m *MsgDeleteAllianceProposal) ProposalType() string   { return ProposalTypeDeleteAlliance }

func (m *MsgDeleteAllianceProposal) ValidateBasic() error {
	if m.Denom == "" {
		return status.Errorf(codes.InvalidArgument, "Alliance denom must have a value")
	}
	return nil
}
