// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: distrib/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	_ "github.com/cosmos/gogoproto/gogoproto"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

// MsgWithdrawAllDelegatorRewards represents delegation withdrawal to a delegator
// from all staked validators.
type MsgWithdrawAllDelegatorRewards struct {
	DelegatorAddress string `protobuf:"bytes,1,opt,name=delegator_address,json=delegatorAddress,proto3" json:"delegator_address,omitempty"`
}

func (m *MsgWithdrawAllDelegatorRewards) Reset()         { *m = MsgWithdrawAllDelegatorRewards{} }
func (m *MsgWithdrawAllDelegatorRewards) String() string { return proto.CompactTextString(m) }
func (*MsgWithdrawAllDelegatorRewards) ProtoMessage()    {}
func (*MsgWithdrawAllDelegatorRewards) Descriptor() ([]byte, []int) {
	return fileDescriptor_94981f9c0db8dca3, []int{0}
}
func (m *MsgWithdrawAllDelegatorRewards) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgWithdrawAllDelegatorRewards) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgWithdrawAllDelegatorRewards.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgWithdrawAllDelegatorRewards) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdrawAllDelegatorRewards.Merge(m, src)
}
func (m *MsgWithdrawAllDelegatorRewards) XXX_Size() int {
	return m.Size()
}
func (m *MsgWithdrawAllDelegatorRewards) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdrawAllDelegatorRewards.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdrawAllDelegatorRewards proto.InternalMessageInfo

// MsgWithdrawAllDelegatorRewardsResponse defines the Msg/WithdrawAllDelegatorRewards response type.
type MsgWithdrawAllDelegatorRewardsResponse struct {
	Amount github_com_cosmos_cosmos_sdk_types.Coins `protobuf:"bytes,1,rep,name=amount,proto3,castrepeated=github.com/cosmos/cosmos-sdk/types.Coins" json:"amount"`
}

func (m *MsgWithdrawAllDelegatorRewardsResponse) Reset() {
	*m = MsgWithdrawAllDelegatorRewardsResponse{}
}
func (m *MsgWithdrawAllDelegatorRewardsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgWithdrawAllDelegatorRewardsResponse) ProtoMessage()    {}
func (*MsgWithdrawAllDelegatorRewardsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_94981f9c0db8dca3, []int{1}
}
func (m *MsgWithdrawAllDelegatorRewardsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgWithdrawAllDelegatorRewardsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgWithdrawAllDelegatorRewardsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgWithdrawAllDelegatorRewardsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgWithdrawAllDelegatorRewardsResponse.Merge(m, src)
}
func (m *MsgWithdrawAllDelegatorRewardsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgWithdrawAllDelegatorRewardsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgWithdrawAllDelegatorRewardsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgWithdrawAllDelegatorRewardsResponse proto.InternalMessageInfo

func (m *MsgWithdrawAllDelegatorRewardsResponse) GetAmount() github_com_cosmos_cosmos_sdk_types.Coins {
	if m != nil {
		return m.Amount
	}
	return nil
}

func init() {
	proto.RegisterType((*MsgWithdrawAllDelegatorRewards)(nil), "distrib.MsgWithdrawAllDelegatorRewards")
	proto.RegisterType((*MsgWithdrawAllDelegatorRewardsResponse)(nil), "distrib.MsgWithdrawAllDelegatorRewardsResponse")
}

func init() { proto.RegisterFile("distrib/tx.proto", fileDescriptor_94981f9c0db8dca3) }

var fileDescriptor_94981f9c0db8dca3 = []byte{
	// 381 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x48, 0xc9, 0x2c, 0x2e,
	0x29, 0xca, 0x4c, 0xd2, 0x2f, 0xa9, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x87, 0x8a,
	0x48, 0x89, 0xa4, 0xe7, 0xa7, 0xe7, 0x83, 0xc5, 0xf4, 0x41, 0x2c, 0x88, 0xb4, 0x94, 0x5c, 0x72,
	0x7e, 0x71, 0x6e, 0x7e, 0xb1, 0x7e, 0x52, 0x62, 0x71, 0xaa, 0x7e, 0x99, 0x61, 0x52, 0x6a, 0x49,
	0xa2, 0xa1, 0x7e, 0x72, 0x7e, 0x66, 0x1e, 0x54, 0x5e, 0x12, 0x22, 0x1f, 0x0f, 0xd1, 0x08, 0xe1,
	0x40, 0xa5, 0xc4, 0xa1, 0x5a, 0x73, 0x8b, 0xd3, 0xf5, 0xcb, 0x0c, 0x41, 0x14, 0x44, 0x42, 0xa9,
	0x9d, 0x91, 0x4b, 0xce, 0xb7, 0x38, 0x3d, 0x3c, 0xb3, 0x24, 0x23, 0xa5, 0x28, 0xb1, 0xdc, 0x31,
	0x27, 0xc7, 0x25, 0x35, 0x27, 0x35, 0x3d, 0xb1, 0x24, 0xbf, 0x28, 0x28, 0xb5, 0x3c, 0xb1, 0x28,
	0xa5, 0x58, 0xc8, 0x95, 0x4b, 0x30, 0x05, 0x26, 0x16, 0x9f, 0x98, 0x92, 0x52, 0x94, 0x5a, 0x5c,
	0x2c, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0xe9, 0x24, 0x71, 0x69, 0x8b, 0xae, 0x08, 0xd4, 0x22, 0x47,
	0x88, 0x4c, 0x70, 0x49, 0x51, 0x66, 0x5e, 0x7a, 0x90, 0x00, 0x5c, 0x0b, 0x54, 0xdc, 0x4a, 0xae,
	0x63, 0x81, 0x3c, 0xc3, 0x8b, 0x05, 0xf2, 0x0c, 0x4d, 0xcf, 0x37, 0x68, 0x61, 0x9a, 0xa8, 0xd4,
	0xcb, 0xc8, 0xa5, 0x86, 0xdf, 0x25, 0x41, 0xa9, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x42, 0xc9,
	0x5c, 0x6c, 0x89, 0xb9, 0xf9, 0xa5, 0x79, 0x25, 0x12, 0x8c, 0x0a, 0xcc, 0x1a, 0xdc, 0x46, 0x92,
	0x7a, 0x50, 0x37, 0x80, 0x42, 0x46, 0x0f, 0x1a, 0x32, 0x7a, 0xce, 0xf9, 0x99, 0x79, 0x4e, 0x06,
	0x27, 0xee, 0xc9, 0x33, 0xac, 0xba, 0x2f, 0xaf, 0x91, 0x9e, 0x59, 0x92, 0x51, 0x9a, 0xa4, 0x97,
	0x9c, 0x9f, 0x0b, 0x0d, 0x19, 0x28, 0xa5, 0x5b, 0x9c, 0x92, 0xad, 0x5f, 0x52, 0x59, 0x90, 0x5a,
	0x0c, 0xd6, 0x50, 0x1c, 0x04, 0x35, 0xda, 0xa8, 0x8e, 0x8b, 0xd9, 0xb7, 0x38, 0x5d, 0xa8, 0x9c,
	0x4b, 0x1a, 0x5f, 0xe0, 0xa8, 0xeb, 0x41, 0xe3, 0x4c, 0x0f, 0xbf, 0xdb, 0xa5, 0xf4, 0x89, 0x54,
	0x08, 0xf3, 0xa4, 0x93, 0xeb, 0x89, 0x47, 0x72, 0x8c, 0x17, 0x1e, 0xc9, 0x31, 0x3e, 0x78, 0x24,
	0xc7, 0x38, 0xe1, 0xb1, 0x1c, 0xc3, 0x85, 0xc7, 0x72, 0x0c, 0x37, 0x1e, 0xcb, 0x31, 0x44, 0x69,
	0x23, 0xf9, 0x25, 0x24, 0x35, 0x31, 0x57, 0xd7, 0xbb, 0x34, 0x2b, 0xb3, 0x28, 0x51, 0x3f, 0x39,
	0xbf, 0x28, 0x55, 0xbf, 0x42, 0x1f, 0x9e, 0xac, 0x40, 0x9e, 0x4a, 0x62, 0x03, 0xc7, 0xb3, 0x31,
	0x20, 0x00, 0x00, 0xff, 0xff, 0x77, 0xfb, 0x12, 0x5b, 0x6e, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// MsgClient is the client API for Msg service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type MsgClient interface {
	// WithdrawAllDelegatorRewards defines a method to withdraw rewards of delegator
	// from all staked validators.
	WithdrawAllDelegatorRewards(ctx context.Context, in *MsgWithdrawAllDelegatorRewards, opts ...grpc.CallOption) (*MsgWithdrawAllDelegatorRewardsResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) WithdrawAllDelegatorRewards(ctx context.Context, in *MsgWithdrawAllDelegatorRewards, opts ...grpc.CallOption) (*MsgWithdrawAllDelegatorRewardsResponse, error) {
	out := new(MsgWithdrawAllDelegatorRewardsResponse)
	err := c.cc.Invoke(ctx, "/distrib.Msg/WithdrawAllDelegatorRewards", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	// WithdrawAllDelegatorRewards defines a method to withdraw rewards of delegator
	// from all staked validators.
	WithdrawAllDelegatorRewards(context.Context, *MsgWithdrawAllDelegatorRewards) (*MsgWithdrawAllDelegatorRewardsResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) WithdrawAllDelegatorRewards(ctx context.Context, req *MsgWithdrawAllDelegatorRewards) (*MsgWithdrawAllDelegatorRewardsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WithdrawAllDelegatorRewards not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_WithdrawAllDelegatorRewards_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgWithdrawAllDelegatorRewards)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).WithdrawAllDelegatorRewards(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distrib.Msg/WithdrawAllDelegatorRewards",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).WithdrawAllDelegatorRewards(ctx, req.(*MsgWithdrawAllDelegatorRewards))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "distrib.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "WithdrawAllDelegatorRewards",
			Handler:    _Msg_WithdrawAllDelegatorRewards_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "distrib/tx.proto",
}

func (m *MsgWithdrawAllDelegatorRewards) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdrawAllDelegatorRewards) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdrawAllDelegatorRewards) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.DelegatorAddress) > 0 {
		i -= len(m.DelegatorAddress)
		copy(dAtA[i:], m.DelegatorAddress)
		i = encodeVarintTx(dAtA, i, uint64(len(m.DelegatorAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgWithdrawAllDelegatorRewardsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgWithdrawAllDelegatorRewardsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgWithdrawAllDelegatorRewardsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for iNdEx := len(m.Amount) - 1; iNdEx >= 0; iNdEx-- {
			{
				size, err := m.Amount[iNdEx].MarshalToSizedBuffer(dAtA[:i])
				if err != nil {
					return 0, err
				}
				i -= size
				i = encodeVarintTx(dAtA, i, uint64(size))
			}
			i--
			dAtA[i] = 0xa
		}
	}
	return len(dAtA) - i, nil
}

func encodeVarintTx(dAtA []byte, offset int, v uint64) int {
	offset -= sovTx(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *MsgWithdrawAllDelegatorRewards) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.DelegatorAddress)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	return n
}

func (m *MsgWithdrawAllDelegatorRewardsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if len(m.Amount) > 0 {
		for _, e := range m.Amount {
			l = e.Size()
			n += 1 + l + sovTx(uint64(l))
		}
	}
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgWithdrawAllDelegatorRewards) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgWithdrawAllDelegatorRewards: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdrawAllDelegatorRewards: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DelegatorAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.DelegatorAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func (m *MsgWithdrawAllDelegatorRewardsResponse) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowTx
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
			return fmt.Errorf("proto: MsgWithdrawAllDelegatorRewardsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgWithdrawAllDelegatorRewardsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Amount", wireType)
			}
			var msglen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowTx
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
				return ErrInvalidLengthTx
			}
			postIndex := iNdEx + msglen
			if postIndex < 0 {
				return ErrInvalidLengthTx
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Amount = append(m.Amount, types.Coin{})
			if err := m.Amount[len(m.Amount)-1].Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipTx(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthTx
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
func skipTx(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
					return 0, ErrIntOverflowTx
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
				return 0, ErrInvalidLengthTx
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupTx
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthTx
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthTx        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowTx          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupTx = fmt.Errorf("proto: unexpected end of group")
)