// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kujira/scheduler/proposal.proto

package types

import (
	fmt "fmt"
	github_com_CosmWasm_wasmd_x_wasm_types "github.com/CosmWasm/wasmd/x/wasm/types"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
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

type CreateHookProposal struct {
	// Title is a short summary
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// Description is a human readable text
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	// The account that will execute the msg on the schedule
	Executor string `protobuf:"bytes,3,opt,name=executor,proto3" json:"executor,omitempty"`
	// The contract that the msg is called on
	Contract  string                                                    `protobuf:"bytes,4,opt,name=contract,proto3" json:"contract,omitempty"`
	Msg       github_com_CosmWasm_wasmd_x_wasm_types.RawContractMessage `protobuf:"bytes,5,opt,name=msg,proto3,casttype=github.com/CosmWasm/wasmd/x/wasm/types.RawContractMessage" json:"msg,omitempty"`
	Frequency int64                                                     `protobuf:"varint,6,opt,name=frequency,proto3" json:"frequency,omitempty"`
	Funds     github_com_cosmos_cosmos_sdk_types.Coins                  `protobuf:"bytes,7,rep,name=funds,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"funds"`
}

func (m *CreateHookProposal) Reset()      { *m = CreateHookProposal{} }
func (*CreateHookProposal) ProtoMessage() {}
func (*CreateHookProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad97a40d5d538195, []int{0}
}
func (m *CreateHookProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CreateHookProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CreateHookProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CreateHookProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateHookProposal.Merge(m, src)
}
func (m *CreateHookProposal) XXX_Size() int {
	return m.Size()
}
func (m *CreateHookProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateHookProposal.DiscardUnknown(m)
}

var xxx_messageInfo_CreateHookProposal proto.InternalMessageInfo

func (m *CreateHookProposal) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *CreateHookProposal) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *CreateHookProposal) GetExecutor() string {
	if m != nil {
		return m.Executor
	}
	return ""
}

func (m *CreateHookProposal) GetContract() string {
	if m != nil {
		return m.Contract
	}
	return ""
}

func (m *CreateHookProposal) GetMsg() github_com_CosmWasm_wasmd_x_wasm_types.RawContractMessage {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (m *CreateHookProposal) GetFrequency() int64 {
	if m != nil {
		return m.Frequency
	}
	return 0
}

func (m *CreateHookProposal) GetFunds() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Funds
	}
	return nil
}

type UpdateHookProposal struct {
	// Title is a short summary
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// Description is a human readable text
	Description string                                                    `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Id          uint64                                                    `protobuf:"varint,3,opt,name=id,proto3" json:"id,omitempty"`
	Executor    string                                                    `protobuf:"bytes,4,opt,name=executor,proto3" json:"executor,omitempty"`
	Contract    string                                                    `protobuf:"bytes,5,opt,name=contract,proto3" json:"contract,omitempty"`
	Msg         github_com_CosmWasm_wasmd_x_wasm_types.RawContractMessage `protobuf:"bytes,6,opt,name=msg,proto3,casttype=github.com/CosmWasm/wasmd/x/wasm/types.RawContractMessage" json:"msg,omitempty"`
	Frequency   int64                                                     `protobuf:"varint,7,opt,name=frequency,proto3" json:"frequency,omitempty"`
	Funds       github_com_cosmos_cosmos_sdk_types.Coins                  `protobuf:"bytes,8,rep,name=funds,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"funds"`
}

func (m *UpdateHookProposal) Reset()      { *m = UpdateHookProposal{} }
func (*UpdateHookProposal) ProtoMessage() {}
func (*UpdateHookProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad97a40d5d538195, []int{1}
}
func (m *UpdateHookProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *UpdateHookProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_UpdateHookProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *UpdateHookProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateHookProposal.Merge(m, src)
}
func (m *UpdateHookProposal) XXX_Size() int {
	return m.Size()
}
func (m *UpdateHookProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateHookProposal.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateHookProposal proto.InternalMessageInfo

func (m *UpdateHookProposal) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *UpdateHookProposal) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *UpdateHookProposal) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *UpdateHookProposal) GetExecutor() string {
	if m != nil {
		return m.Executor
	}
	return ""
}

func (m *UpdateHookProposal) GetContract() string {
	if m != nil {
		return m.Contract
	}
	return ""
}

func (m *UpdateHookProposal) GetMsg() github_com_CosmWasm_wasmd_x_wasm_types.RawContractMessage {
	if m != nil {
		return m.Msg
	}
	return nil
}

func (m *UpdateHookProposal) GetFrequency() int64 {
	if m != nil {
		return m.Frequency
	}
	return 0
}

func (m *UpdateHookProposal) GetFunds() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Funds
	}
	return nil
}

type DeleteHookProposal struct {
	// Title is a short summary
	Title string `protobuf:"bytes,1,opt,name=title,proto3" json:"title,omitempty"`
	// Description is a human readable text
	Description string `protobuf:"bytes,2,opt,name=description,proto3" json:"description,omitempty"`
	Id          uint64 `protobuf:"varint,3,opt,name=id,proto3" json:"id,omitempty"`
}

func (m *DeleteHookProposal) Reset()      { *m = DeleteHookProposal{} }
func (*DeleteHookProposal) ProtoMessage() {}
func (*DeleteHookProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad97a40d5d538195, []int{2}
}
func (m *DeleteHookProposal) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *DeleteHookProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_DeleteHookProposal.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *DeleteHookProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteHookProposal.Merge(m, src)
}
func (m *DeleteHookProposal) XXX_Size() int {
	return m.Size()
}
func (m *DeleteHookProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteHookProposal.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteHookProposal proto.InternalMessageInfo

func (m *DeleteHookProposal) GetTitle() string {
	if m != nil {
		return m.Title
	}
	return ""
}

func (m *DeleteHookProposal) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *DeleteHookProposal) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func init() {
	proto.RegisterType((*CreateHookProposal)(nil), "kujira.scheduler.CreateHookProposal")
	proto.RegisterType((*UpdateHookProposal)(nil), "kujira.scheduler.UpdateHookProposal")
	proto.RegisterType((*DeleteHookProposal)(nil), "kujira.scheduler.DeleteHookProposal")
}

func init() { proto.RegisterFile("kujira/scheduler/proposal.proto", fileDescriptor_ad97a40d5d538195) }

var fileDescriptor_ad97a40d5d538195 = []byte{
	// 455 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x93, 0x4d, 0x6f, 0xd3, 0x4e,
	0x10, 0xc6, 0x6d, 0xe7, 0xa5, 0xed, 0xf6, 0xaf, 0xbf, 0xd0, 0xaa, 0x07, 0x13, 0x90, 0x63, 0xf5,
	0xe4, 0x4b, 0xbc, 0x14, 0x4e, 0x1c, 0xb8, 0x24, 0x1c, 0x90, 0x2a, 0x04, 0xb2, 0x40, 0x48, 0x88,
	0xcb, 0x66, 0x3d, 0x75, 0x96, 0xc4, 0x1e, 0xb3, 0xbb, 0xa6, 0xe9, 0xb7, 0xe0, 0x73, 0xf0, 0x49,
	0x7a, 0xec, 0x09, 0xf5, 0x54, 0x20, 0x91, 0xf8, 0x10, 0x9c, 0x90, 0xbd, 0xa6, 0x04, 0xf1, 0x72,
	0x22, 0xa7, 0xd9, 0x99, 0x67, 0xf6, 0xd1, 0xa3, 0x9f, 0x34, 0x64, 0x38, 0xaf, 0x5e, 0x4b, 0xc5,
	0x99, 0x16, 0x33, 0x48, 0xab, 0x05, 0x28, 0x56, 0x2a, 0x2c, 0x51, 0xf3, 0x45, 0x5c, 0x2a, 0x34,
	0x48, 0x6f, 0xd8, 0x85, 0xf8, 0x7a, 0x61, 0x70, 0x90, 0x61, 0x86, 0x8d, 0xc8, 0xea, 0x97, 0xdd,
	0x1b, 0xdc, 0xfa, 0xc5, 0x68, 0x86, 0x38, 0x6f, 0xc5, 0x40, 0xa0, 0xce, 0x51, 0xb3, 0x29, 0xd7,
	0xc0, 0xde, 0x1e, 0x4d, 0xc1, 0xf0, 0x23, 0x26, 0x50, 0x16, 0x56, 0x3f, 0xfc, 0xe0, 0x11, 0x3a,
	0x51, 0xc0, 0x0d, 0x3c, 0x42, 0x9c, 0x3f, 0x6d, 0x13, 0xd0, 0x03, 0xd2, 0x33, 0xd2, 0x2c, 0xc0,
	0x77, 0x43, 0x37, 0xda, 0x4b, 0x6c, 0x43, 0x43, 0xb2, 0x9f, 0x82, 0x16, 0x4a, 0x96, 0x46, 0x62,
	0xe1, 0x7b, 0x8d, 0xb6, 0x39, 0xa2, 0x03, 0xb2, 0x0b, 0x4b, 0x10, 0x95, 0x41, 0xe5, 0x77, 0x1a,
	0xf9, 0xba, 0xaf, 0x35, 0x81, 0x85, 0x51, 0x5c, 0x18, 0xbf, 0x6b, 0xb5, 0xef, 0x3d, 0x7d, 0x42,
	0x3a, 0xb9, 0xce, 0xfc, 0x5e, 0xe8, 0x46, 0xff, 0x8d, 0x1f, 0x7c, 0xbd, 0x1a, 0xde, 0xcf, 0xa4,
	0x99, 0x55, 0xd3, 0x58, 0x60, 0xce, 0x26, 0xa8, 0xf3, 0x17, 0x5c, 0xe7, 0xec, 0x94, 0xeb, 0x3c,
	0x65, 0xcb, 0xa6, 0x32, 0x73, 0x56, 0x82, 0x8e, 0x13, 0x7e, 0x3a, 0x69, 0x4d, 0x1e, 0x83, 0xd6,
	0x3c, 0x83, 0xa4, 0x76, 0xa2, 0xb7, 0xc9, 0xde, 0x89, 0x82, 0x37, 0x15, 0x14, 0xe2, 0xcc, 0xef,
	0x87, 0x6e, 0xd4, 0x49, 0x7e, 0x0c, 0x28, 0x27, 0xbd, 0x93, 0xaa, 0x48, 0xb5, 0xbf, 0x13, 0x76,
	0xa2, 0xfd, 0xbb, 0x37, 0x63, 0x4b, 0x29, 0xae, 0x29, 0xc5, 0x2d, 0xa5, 0x78, 0x82, 0xb2, 0x18,
	0xdf, 0x39, 0xbf, 0x1a, 0x3a, 0xef, 0x3f, 0x0e, 0xa3, 0x8d, 0x3c, 0x2d, 0x52, 0x5b, 0x46, 0x3a,
	0x9d, 0xb7, 0x59, 0xea, 0x0f, 0x3a, 0xb1, 0xce, 0x87, 0x5f, 0x3c, 0x42, 0x9f, 0x97, 0xe9, 0xbf,
	0x02, 0xfb, 0x3f, 0xf1, 0x64, 0xda, 0x20, 0xed, 0x26, 0x9e, 0x4c, 0x7f, 0x02, 0xdd, 0xfd, 0x0b,
	0xe8, 0xde, 0xef, 0x41, 0xf7, 0xb7, 0x03, 0x7a, 0xe7, 0x8f, 0xa0, 0x77, 0xb7, 0x06, 0xfa, 0x15,
	0xa1, 0x0f, 0x61, 0x01, 0xdb, 0xe1, 0x3c, 0x3e, 0xbe, 0xfc, 0x1c, 0x38, 0xe7, 0xab, 0xc0, 0xbd,
	0x58, 0x05, 0xee, 0xa7, 0x55, 0xe0, 0xbe, 0x5b, 0x07, 0xce, 0xc5, 0x3a, 0x70, 0x2e, 0xd7, 0x81,
	0xf3, 0x72, 0xb4, 0x11, 0xf6, 0x19, 0xf0, 0x7c, 0x74, 0x6c, 0x4f, 0x51, 0xa0, 0x02, 0xb6, 0xdc,
	0xb8, 0xc8, 0x26, 0xf7, 0xb4, 0xdf, 0xdc, 0xdc, 0xbd, 0x6f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x21,
	0xb6, 0x20, 0x7c, 0xfb, 0x03, 0x00, 0x00,
}

func (m *CreateHookProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CreateHookProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CreateHookProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Funds) > 0 {
		for iNdEx := len(m.Funds) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Funds[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintProposal(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x3a
		}
	}
	if m.Frequency != 0 {
		i = encodeVarintProposal(dAtA, i, uint64(m.Frequency))
		i--
		dAtA[i] = 0x30
	}
	if len(m.Msg) > 0 {
		i -= len(m.Msg)
		copy(dAtA[i:], m.Msg)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Msg)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Contract) > 0 {
		i -= len(m.Contract)
		copy(dAtA[i:], m.Contract)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Contract)))
		i--
		dAtA[i] = 0x22
	}
	if len(m.Executor) > 0 {
		i -= len(m.Executor)
		copy(dAtA[i:], m.Executor)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Executor)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *UpdateHookProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *UpdateHookProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *UpdateHookProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Funds) > 0 {
		for iNdEx := len(m.Funds) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Funds[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintProposal(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0x42
		}
	}
	if m.Frequency != 0 {
		i = encodeVarintProposal(dAtA, i, uint64(m.Frequency))
		i--
		dAtA[i] = 0x38
	}
	if len(m.Msg) > 0 {
		i -= len(m.Msg)
		copy(dAtA[i:], m.Msg)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Msg)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.Contract) > 0 {
		i -= len(m.Contract)
		copy(dAtA[i:], m.Contract)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Contract)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Executor) > 0 {
		i -= len(m.Executor)
		copy(dAtA[i:], m.Executor)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Executor)))
		i--
		dAtA[i] = 0x22
	}
	if m.Id != 0 {
		i = encodeVarintProposal(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *DeleteHookProposal) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *DeleteHookProposal) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *DeleteHookProposal) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Id != 0 {
		i = encodeVarintProposal(dAtA, i, uint64(m.Id))
		i--
		dAtA[i] = 0x18
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintProposal(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintProposal(dAtA []byte, offset int, v uint64) int {
	offset -= sovProposal(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *CreateHookProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Executor)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Contract)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Msg)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	if m.Frequency != 0 {
		n += 1 + sovProposal(uint64(m.Frequency))
	}
	if len(m.Funds) > 0 {
		for _, e := range m.Funds {
			l = e.Size()
			n += 1 + l + sovProposal(uint64(l))
		}
	}
	return n
}

func (m *UpdateHookProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	if m.Id != 0 {
		n += 1 + sovProposal(uint64(m.Id))
	}
	l = len(m.Executor)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Contract)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Msg)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	if m.Frequency != 0 {
		n += 1 + sovProposal(uint64(m.Frequency))
	}
	if len(m.Funds) > 0 {
		for _, e := range m.Funds {
			l = e.Size()
			n += 1 + l + sovProposal(uint64(l))
		}
	}
	return n
}

func (m *DeleteHookProposal) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovProposal(uint64(l))
	}
	if m.Id != 0 {
		n += 1 + sovProposal(uint64(m.Id))
	}
	return n
}

func sovProposal(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozProposal(x uint64) (n int) {
	return sovProposal(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *CreateHookProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProposal
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
			return fmt.Errorf("proto: CreateHookProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CreateHookProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Executor", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Executor = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Contract", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Contract = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msg = append(m.Msg[:0], dAtA[iNdEx:postIndex]...)
			if m.Msg == nil {
				m.Msg = []byte{}
			}
			iNdEx = postIndex
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Frequency", wireType)
			}
			m.Frequency = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Frequency |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Funds", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Funds = append(m.Funds, types.Coin{})
			if err := m.Funds[len(m.Funds)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProposal
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
func (m *UpdateHookProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProposal
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
			return fmt.Errorf("proto: UpdateHookProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: UpdateHookProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Executor", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Executor = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Contract", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Contract = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Msg", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				byteLen |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if byteLen < 0 {
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Msg = append(m.Msg[:0], dAtA[iNdEx:postIndex]...)
			if m.Msg == nil {
				m.Msg = []byte{}
			}
			iNdEx = postIndex
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Frequency", wireType)
			}
			m.Frequency = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Frequency |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 8:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Funds", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Funds = append(m.Funds, types.Coin{})
			if err := m.Funds[len(m.Funds)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProposal
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
func (m *DeleteHookProposal) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowProposal
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
			return fmt.Errorf("proto: DeleteHookProposal: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: DeleteHookProposal: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
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
				return ErrInvalidLengthProposal
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthProposal
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Id", wireType)
			}
			m.Id = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowProposal
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Id |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipProposal(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthProposal
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
func skipProposal(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowProposal
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
					return 0, ErrIntOverflowProposal
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
					return 0, ErrIntOverflowProposal
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
				return 0, ErrInvalidLengthProposal
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupProposal
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthProposal
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthProposal        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowProposal          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupProposal = fmt.Errorf("proto: unexpected end of group")
)
