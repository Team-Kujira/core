// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kujira/oracle/oracle.proto

package types

import (
	cosmossdk_io_math "cosmossdk.io/math"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// Params defines the parameters for the oracle module.
type Params struct {
	VotePeriod        uint64                      `protobuf:"varint,1,opt,name=vote_period,json=votePeriod,proto3" json:"vote_period,omitempty" yaml:"vote_period"`
	VoteThreshold     cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=vote_threshold,json=voteThreshold,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"vote_threshold" yaml:"vote_threshold"`
	MaxDeviation      cosmossdk_io_math.LegacyDec `protobuf:"bytes,3,opt,name=max_deviation,json=maxDeviation,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"max_deviation" yaml:"max_deviation"`
	RequiredDenoms    []string                    `protobuf:"bytes,4,rep,name=required_denoms,json=requiredDenoms,proto3" json:"required_denoms,omitempty" yaml:"required_denoms"`
	SlashFraction     cosmossdk_io_math.LegacyDec `protobuf:"bytes,5,opt,name=slash_fraction,json=slashFraction,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"slash_fraction" yaml:"slash_fraction"`
	SlashWindow       uint64                      `protobuf:"varint,6,opt,name=slash_window,json=slashWindow,proto3" json:"slash_window,omitempty" yaml:"slash_window"`
	MinValidPerWindow cosmossdk_io_math.LegacyDec `protobuf:"bytes,7,opt,name=min_valid_per_window,json=minValidPerWindow,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"min_valid_per_window" yaml:"min_valid_per_window"`
	// Deprecated
	RewardBand cosmossdk_io_math.LegacyDec `protobuf:"bytes,8,opt,name=reward_band,json=rewardBand,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"reward_band" yaml:"reward_band"`
	Whitelist  DenomList                   `protobuf:"bytes,9,rep,name=whitelist,proto3,castrepeated=DenomList" json:"whitelist" yaml:"whitelist"`
}

func (m *Params) Reset()      { *m = Params{} }
func (*Params) ProtoMessage() {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_8fffe8fb5ee63325, []int{0}
}
func (m *Params) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Params) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Params.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Params) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Params.Merge(m, src)
}
func (m *Params) XXX_Size() int {
	return m.Size()
}
func (m *Params) XXX_DiscardUnknown() {
	xxx_messageInfo_Params.DiscardUnknown(m)
}

var xxx_messageInfo_Params proto.InternalMessageInfo

func (m *Params) GetVotePeriod() uint64 {
	if m != nil {
		return m.VotePeriod
	}
	return 0
}

func (m *Params) GetRequiredDenoms() []string {
	if m != nil {
		return m.RequiredDenoms
	}
	return nil
}

func (m *Params) GetSlashWindow() uint64 {
	if m != nil {
		return m.SlashWindow
	}
	return 0
}

func (m *Params) GetWhitelist() DenomList {
	if m != nil {
		return m.Whitelist
	}
	return nil
}

// Denom - the object to hold configurations of each denom
type Denom struct {
	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty" yaml:"name"`
}

func (m *Denom) Reset()      { *m = Denom{} }
func (*Denom) ProtoMessage() {}
func (*Denom) Descriptor() ([]byte, []int) {
	return fileDescriptor_8fffe8fb5ee63325, []int{1}
}
func (m *Denom) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *Denom) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_Denom.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *Denom) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Denom.Merge(m, src)
}
func (m *Denom) XXX_Size() int {
	return m.Size()
}
func (m *Denom) XXX_DiscardUnknown() {
	xxx_messageInfo_Denom.DiscardUnknown(m)
}

var xxx_messageInfo_Denom proto.InternalMessageInfo

// ExchangeRateTuple - struct to store interpreted exchange rates data to store
type ExchangeRateTuple struct {
	Denom        string                      `protobuf:"bytes,1,opt,name=denom,proto3" json:"denom,omitempty" yaml:"denom"`
	ExchangeRate cosmossdk_io_math.LegacyDec `protobuf:"bytes,2,opt,name=exchange_rate,json=exchangeRate,proto3,customtype=cosmossdk.io/math.LegacyDec" json:"exchange_rate" yaml:"exchange_rate"`
}

func (m *ExchangeRateTuple) Reset()      { *m = ExchangeRateTuple{} }
func (*ExchangeRateTuple) ProtoMessage() {}
func (*ExchangeRateTuple) Descriptor() ([]byte, []int) {
	return fileDescriptor_8fffe8fb5ee63325, []int{2}
}
func (m *ExchangeRateTuple) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *ExchangeRateTuple) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_ExchangeRateTuple.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *ExchangeRateTuple) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExchangeRateTuple.Merge(m, src)
}
func (m *ExchangeRateTuple) XXX_Size() int {
	return m.Size()
}
func (m *ExchangeRateTuple) XXX_DiscardUnknown() {
	xxx_messageInfo_ExchangeRateTuple.DiscardUnknown(m)
}

var xxx_messageInfo_ExchangeRateTuple proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Params)(nil), "kujira.oracle.Params")
	proto.RegisterType((*Denom)(nil), "kujira.oracle.Denom")
	proto.RegisterType((*ExchangeRateTuple)(nil), "kujira.oracle.ExchangeRateTuple")
}

func init() { proto.RegisterFile("kujira/oracle/oracle.proto", fileDescriptor_8fffe8fb5ee63325) }

var fileDescriptor_8fffe8fb5ee63325 = []byte{
	// 639 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0x3d, 0x6f, 0xd3, 0x40,
	0x18, 0xb6, 0x69, 0x1b, 0x9a, 0x4b, 0xd2, 0x0f, 0x13, 0x8a, 0x95, 0x4a, 0x76, 0x74, 0x08, 0x14,
	0x21, 0x61, 0xab, 0x30, 0x20, 0x02, 0x93, 0x09, 0x2c, 0x74, 0xa8, 0xac, 0x0a, 0x44, 0x17, 0x73,
	0xb1, 0x8f, 0xf8, 0xa8, 0xed, 0x0b, 0xe7, 0xcb, 0x47, 0xff, 0x01, 0x23, 0x23, 0x63, 0x07, 0x26,
	0xf8, 0x23, 0x1d, 0x3b, 0x22, 0x06, 0x83, 0xda, 0x85, 0xd9, 0xbf, 0x00, 0xf9, 0xce, 0x29, 0x09,
	0x62, 0x28, 0x4c, 0xf6, 0xf3, 0x3c, 0xef, 0xfb, 0x3c, 0xaf, 0xee, 0x0b, 0xb4, 0x0e, 0x47, 0x6f,
	0x09, 0x43, 0x36, 0x65, 0xc8, 0x8f, 0x70, 0xf9, 0xb1, 0x86, 0x8c, 0x72, 0xaa, 0x35, 0xa4, 0x66,
	0x49, 0xb2, 0xd5, 0x1c, 0xd0, 0x01, 0x15, 0x8a, 0x5d, 0xfc, 0xc9, 0xa2, 0x96, 0xe1, 0xd3, 0x34,
	0xa6, 0xa9, 0xdd, 0x47, 0x29, 0xb6, 0xc7, 0x3b, 0x7d, 0xcc, 0xd1, 0x8e, 0xed, 0x53, 0x92, 0x48,
	0x1d, 0x7e, 0xaa, 0x80, 0xca, 0x1e, 0x62, 0x28, 0x4e, 0xb5, 0x07, 0xa0, 0x36, 0xa6, 0x1c, 0x7b,
	0x43, 0xcc, 0x08, 0x0d, 0x74, 0xb5, 0xad, 0x76, 0x96, 0x9d, 0xad, 0x3c, 0x33, 0xb5, 0x23, 0x14,
	0x47, 0x5d, 0x38, 0x27, 0x42, 0x17, 0x14, 0x68, 0x4f, 0x00, 0xcd, 0x07, 0x6b, 0x42, 0xe3, 0x21,
	0xc3, 0x69, 0x48, 0xa3, 0x40, 0xbf, 0xd2, 0x56, 0x3b, 0x55, 0xe7, 0xf1, 0x49, 0x66, 0x2a, 0xdf,
	0x32, 0x73, 0x5b, 0xce, 0x90, 0x06, 0x87, 0x16, 0xa1, 0x76, 0x8c, 0x78, 0x68, 0xed, 0xe2, 0x01,
	0xf2, 0x8f, 0x7a, 0xd8, 0xcf, 0x33, 0xf3, 0xfa, 0x9c, 0xfd, 0x85, 0x05, 0x74, 0x1b, 0x05, 0xb1,
	0x3f, 0xc3, 0xda, 0x6b, 0xd0, 0x88, 0xd1, 0xd4, 0x0b, 0xf0, 0x98, 0x20, 0x4e, 0x68, 0xa2, 0x2f,
	0x89, 0x8c, 0x47, 0x97, 0xcb, 0x68, 0xca, 0x8c, 0x05, 0x07, 0xe8, 0xd6, 0x63, 0x34, 0xed, 0xcd,
	0xa0, 0xf6, 0x04, 0xac, 0x33, 0xfc, 0x6e, 0x44, 0x18, 0x0e, 0xbc, 0x00, 0x27, 0x34, 0x4e, 0xf5,
	0xe5, 0xf6, 0x52, 0xa7, 0xea, 0xb4, 0xf2, 0xcc, 0xdc, 0x92, 0x06, 0x7f, 0x14, 0x40, 0x77, 0x6d,
	0xc6, 0xf4, 0x04, 0x51, 0xac, 0x45, 0x1a, 0xa1, 0x34, 0xf4, 0xde, 0x30, 0xe4, 0x8b, 0x39, 0x57,
	0xfe, 0x63, 0x2d, 0x16, 0x2d, 0xa0, 0xdb, 0x10, 0xc4, 0xb3, 0x12, 0x6b, 0x5d, 0x50, 0x97, 0x15,
	0x13, 0x92, 0x04, 0x74, 0xa2, 0x57, 0xc4, 0x56, 0xdd, 0xc8, 0x33, 0xf3, 0xda, 0x7c, 0xbf, 0x54,
	0xa1, 0x5b, 0x13, 0xf0, 0xa5, 0x40, 0x5a, 0x0a, 0x9a, 0x31, 0x49, 0xbc, 0x31, 0x8a, 0x48, 0x50,
	0xec, 0xe6, 0xcc, 0xe3, 0xaa, 0x18, 0xd3, 0xb9, 0xdc, 0x98, 0xdb, 0xe5, 0x72, 0xfe, 0xc5, 0x08,
	0xba, 0x9b, 0x31, 0x49, 0x5e, 0x14, 0xec, 0x1e, 0x66, 0x65, 0xe8, 0x01, 0xa8, 0x31, 0x3c, 0x41,
	0x2c, 0xf0, 0xfa, 0x28, 0x09, 0xf4, 0x55, 0x91, 0xf5, 0xf0, 0x72, 0x59, 0xda, 0x6c, 0xe5, 0x2f,
	0xfa, 0xa1, 0x0b, 0x24, 0x72, 0x50, 0x12, 0x68, 0xaf, 0x40, 0x75, 0x12, 0x12, 0x8e, 0x23, 0x92,
	0x72, 0xbd, 0xda, 0x5e, 0xea, 0xd4, 0xee, 0x35, 0xad, 0x85, 0xab, 0x61, 0x89, 0xbd, 0x71, 0x6e,
	0x15, 0x79, 0x79, 0x66, 0x6e, 0x48, 0xc3, 0x8b, 0x26, 0xf8, 0xf9, 0xbb, 0x59, 0x15, 0x25, 0xbb,
	0x24, 0xe5, 0xee, 0x6f, 0xb7, 0xee, 0xea, 0xc7, 0x63, 0x53, 0xf9, 0x79, 0x6c, 0xaa, 0xb0, 0x0b,
	0x56, 0x44, 0x85, 0x76, 0x13, 0x2c, 0x27, 0x28, 0xc6, 0xe2, 0x76, 0x54, 0x9d, 0xf5, 0x3c, 0x33,
	0x6b, 0xd2, 0xae, 0x60, 0xa1, 0x2b, 0xc4, 0x6e, 0xfd, 0xfd, 0xb1, 0xa9, 0x94, 0xbd, 0x0a, 0xfc,
	0xa2, 0x82, 0xcd, 0xa7, 0x53, 0x3f, 0x44, 0xc9, 0x00, 0xbb, 0x88, 0xe3, 0xfd, 0xd1, 0x30, 0xc2,
	0xda, 0x6d, 0xb0, 0x22, 0xce, 0x50, 0xe9, 0xb4, 0x91, 0x67, 0x66, 0x5d, 0x3a, 0x09, 0x1a, 0xba,
	0x52, 0x2e, 0xce, 0x3d, 0x2e, 0x9b, 0x3d, 0x86, 0x38, 0x2e, 0xef, 0xd6, 0xbf, 0x9d, 0xfb, 0x05,
	0x07, 0xe8, 0xd6, 0xf1, 0xdc, 0x38, 0x8b, 0xd3, 0x3a, 0xbd, 0x93, 0x33, 0x43, 0x3d, 0x3d, 0x33,
	0xd4, 0x1f, 0x67, 0x86, 0xfa, 0xe1, 0xdc, 0x50, 0x4e, 0xcf, 0x0d, 0xe5, 0xeb, 0xb9, 0xa1, 0x1c,
	0xdc, 0x19, 0x10, 0x1e, 0x8e, 0xfa, 0x96, 0x4f, 0x63, 0x7b, 0x1f, 0xa3, 0xf8, 0xee, 0x73, 0xf9,
	0x36, 0xf9, 0x94, 0x61, 0x7b, 0x3a, 0x7b, 0xa2, 0xf8, 0xd1, 0x10, 0xa7, 0xfd, 0x8a, 0x78, 0x5d,
	0xee, 0xff, 0x0a, 0x00, 0x00, 0xff, 0xff, 0xf9, 0xcd, 0x96, 0xa6, 0xc0, 0x04, 0x00, 0x00,
}

func (this *Params) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*Params)
	if !ok {
		that2, ok := that.(Params)
		if ok {
			that1 = &that2
		} else {
			return false
		}
	}
	if that1 == nil {
		return this == nil
	} else if this == nil {
		return false
	}
	if this.VotePeriod != that1.VotePeriod {
		return false
	}
	if !this.VoteThreshold.Equal(that1.VoteThreshold) {
		return false
	}
	if !this.MaxDeviation.Equal(that1.MaxDeviation) {
		return false
	}
	if len(this.RequiredDenoms) != len(that1.RequiredDenoms) {
		return false
	}
	for i := range this.RequiredDenoms {
		if this.RequiredDenoms[i] != that1.RequiredDenoms[i] {
			return false
		}
	}
	if !this.SlashFraction.Equal(that1.SlashFraction) {
		return false
	}
	if this.SlashWindow != that1.SlashWindow {
		return false
	}
	if !this.MinValidPerWindow.Equal(that1.MinValidPerWindow) {
		return false
	}
	if !this.RewardBand.Equal(that1.RewardBand) {
		return false
	}
	if len(this.Whitelist) != len(that1.Whitelist) {
		return false
	}
	for i := range this.Whitelist {
		if !this.Whitelist[i].Equal(&that1.Whitelist[i]) {
			return false
		}
	}
	return true
}
func (m *Params) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Params) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Params) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Whitelist) > 0 {
		for iNdEx := len(m.Whitelist) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Whitelist[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintOracle(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x4a
		}
	}
	{
		size := m.RewardBand.Size()
		i -= size
		if _, err := m.RewardBand.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOracle(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x42
	{
		size := m.MinValidPerWindow.Size()
		i -= size
		if _, err := m.MinValidPerWindow.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOracle(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x3a
	if m.SlashWindow != 0 {
		i = encodeVarintOracle(dAtA, i, uint64(m.SlashWindow))
		i--
		dAtA[i] = 0x30
	}
	{
		size := m.SlashFraction.Size()
		i -= size
		if _, err := m.SlashFraction.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOracle(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if len(m.RequiredDenoms) > 0 {
		for iNdEx := len(m.RequiredDenoms) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.RequiredDenoms[iNdEx])
			copy(dAtA[i:], m.RequiredDenoms[iNdEx])
			i = encodeVarintOracle(dAtA, i, uint64(len(m.RequiredDenoms[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	{
		size := m.MaxDeviation.Size()
		i -= size
		if _, err := m.MaxDeviation.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOracle(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x1a
	{
		size := m.VoteThreshold.Size()
		i -= size
		if _, err := m.VoteThreshold.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOracle(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if m.VotePeriod != 0 {
		i = encodeVarintOracle(dAtA, i, uint64(m.VotePeriod))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *Denom) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *Denom) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *Denom) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Name) > 0 {
		i -= len(m.Name)
		copy(dAtA[i:], m.Name)
		i = encodeVarintOracle(dAtA, i, uint64(len(m.Name)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *ExchangeRateTuple) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *ExchangeRateTuple) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *ExchangeRateTuple) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size := m.ExchangeRate.Size()
		i -= size
		if _, err := m.ExchangeRate.MarshalTo(dAtA[i:]); err != nil {
			return 0, err
		}
		i = encodeVarintOracle(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x12
	if len(m.Denom) > 0 {
		i -= len(m.Denom)
		copy(dAtA[i:], m.Denom)
		i = encodeVarintOracle(dAtA, i, uint64(len(m.Denom)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintOracle(dAtA []byte, offset int, v uint64) int {
	offset -= sovOracle(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *Params) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.VotePeriod != 0 {
		n += 1 + sovOracle(uint64(m.VotePeriod))
	}
	l = m.VoteThreshold.Size()
	n += 1 + l + sovOracle(uint64(l))
	l = m.MaxDeviation.Size()
	n += 1 + l + sovOracle(uint64(l))
	if len(m.RequiredDenoms) > 0 {
		for _, s := range m.RequiredDenoms {
			l = len(s)
			n += 1 + l + sovOracle(uint64(l))
		}
	}
	l = m.SlashFraction.Size()
	n += 1 + l + sovOracle(uint64(l))
	if m.SlashWindow != 0 {
		n += 1 + sovOracle(uint64(m.SlashWindow))
	}
	l = m.MinValidPerWindow.Size()
	n += 1 + l + sovOracle(uint64(l))
	l = m.RewardBand.Size()
	n += 1 + l + sovOracle(uint64(l))
	if len(m.Whitelist) > 0 {
		for _, e := range m.Whitelist {
			l = e.Size()
			n += 1 + l + sovOracle(uint64(l))
		}
	}
	return n
}

func (m *Denom) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Name)
	if l > 0 {
		n += 1 + l + sovOracle(uint64(l))
	}
	return n
}

func (m *ExchangeRateTuple) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Denom)
	if l > 0 {
		n += 1 + l + sovOracle(uint64(l))
	}
	l = m.ExchangeRate.Size()
	n += 1 + l + sovOracle(uint64(l))
	return n
}

func sovOracle(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozOracle(x uint64) (n int) {
	return sovOracle(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOracle
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Params: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Params: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VotePeriod", wireType)
			}
			m.VotePeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VotePeriod |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field VoteThreshold", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.VoteThreshold.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaxDeviation", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MaxDeviation.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RequiredDenoms", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.RequiredDenoms = append(m.RequiredDenoms, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashFraction", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.SlashFraction.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SlashWindow", wireType)
			}
			m.SlashWindow = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SlashWindow |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MinValidPerWindow", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.MinValidPerWindow.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field RewardBand", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.RewardBand.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 9:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Whitelist", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				msglen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if msglen < 0 {
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Whitelist = append(m.Whitelist, Denom{})
			if err := m.Whitelist[len(m.Whitelist)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipOracle(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthOracle
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *Denom) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOracle
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: Denom: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: Denom: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Name", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Name = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipOracle(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthOracle
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func (m *ExchangeRateTuple) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowOracle
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: ExchangeRateTuple: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: ExchangeRateTuple: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Denom", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Denom = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExchangeRate", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				stringLen |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			intStringLen := int(stringLen)
			if intStringLen < 0 {
				return ErrInvalidLengthOracle
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthOracle
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			if err := m.ExchangeRate.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipOracle(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthOracle
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipOracle(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowOracle
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowOracle
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthOracle
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupOracle
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthOracle
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthOracle        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowOracle          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupOracle = fmt.Errorf("proto: unexpected end of group")
)
