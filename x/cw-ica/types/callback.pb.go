// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: kujira/cwica/callback.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/cosmos/gogoproto/proto"
	_ "google.golang.org/genproto/googleapis/api/annotations"
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

type CallbackData struct {
	PortId       string `protobuf:"bytes,1,opt,name=port_id,json=portId,proto3" json:"port_id,omitempty"`
	ChannelId    string `protobuf:"bytes,2,opt,name=channel_id,json=channelId,proto3" json:"channel_id,omitempty"`
	Sequence     uint64 `protobuf:"varint,3,opt,name=sequence,proto3" json:"sequence,omitempty"`
	Contract     string `protobuf:"bytes,4,opt,name=contract,proto3" json:"contract,omitempty"`
	ConnectionId string `protobuf:"bytes,5,opt,name=connection_id,json=connectionId,proto3" json:"connection_id,omitempty"`
	AccountId    string `protobuf:"bytes,6,opt,name=account_id,json=accountId,proto3" json:"account_id,omitempty"`
	Callback     []byte `protobuf:"bytes,7,opt,name=callback,proto3" json:"callback,omitempty"`
}

func (m *CallbackData) Reset()         { *m = CallbackData{} }
func (m *CallbackData) String() string { return proto.CompactTextString(m) }
func (*CallbackData) ProtoMessage()    {}
func (*CallbackData) Descriptor() ([]byte, []int) {
	return fileDescriptor_c1de6cedd893c366, []int{0}
}
func (m *CallbackData) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CallbackData) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CallbackData.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CallbackData) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CallbackData.Merge(m, src)
}
func (m *CallbackData) XXX_Size() int {
	return m.Size()
}
func (m *CallbackData) XXX_DiscardUnknown() {
	xxx_messageInfo_CallbackData.DiscardUnknown(m)
}

var xxx_messageInfo_CallbackData proto.InternalMessageInfo

func (m *CallbackData) GetPortId() string {
	if m != nil {
		return m.PortId
	}
	return ""
}

func (m *CallbackData) GetChannelId() string {
	if m != nil {
		return m.ChannelId
	}
	return ""
}

func (m *CallbackData) GetSequence() uint64 {
	if m != nil {
		return m.Sequence
	}
	return 0
}

func (m *CallbackData) GetContract() string {
	if m != nil {
		return m.Contract
	}
	return ""
}

func (m *CallbackData) GetConnectionId() string {
	if m != nil {
		return m.ConnectionId
	}
	return ""
}

func (m *CallbackData) GetAccountId() string {
	if m != nil {
		return m.AccountId
	}
	return ""
}

func (m *CallbackData) GetCallback() []byte {
	if m != nil {
		return m.Callback
	}
	return nil
}

func init() {
	proto.RegisterType((*CallbackData)(nil), "kujira.cwica.CallbackData")
}

func init() { proto.RegisterFile("kujira/cwica/callback.proto", fileDescriptor_c1de6cedd893c366) }

var fileDescriptor_c1de6cedd893c366 = []byte{
	// 296 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0xc1, 0x4e, 0x02, 0x31,
	0x10, 0x86, 0xa9, 0x22, 0x48, 0xb3, 0x5e, 0x36, 0x26, 0x6e, 0x50, 0x1b, 0xa2, 0x17, 0x62, 0x02,
	0x3d, 0xf8, 0x06, 0xca, 0x85, 0x78, 0x23, 0x9e, 0xbc, 0x98, 0x61, 0x68, 0x96, 0xca, 0xd2, 0x59,
	0x97, 0x6e, 0xd0, 0xb7, 0xf0, 0xb1, 0x3c, 0x72, 0xf4, 0x64, 0x0c, 0xfb, 0x22, 0xa6, 0xed, 0xba,
	0xde, 0xfa, 0xff, 0xdf, 0x74, 0xfe, 0xc9, 0xcf, 0xcf, 0x57, 0xe5, 0x8b, 0x2e, 0x40, 0xe2, 0x56,
	0x23, 0x48, 0x84, 0x2c, 0x9b, 0x03, 0xae, 0xc6, 0x79, 0x41, 0x96, 0xe2, 0x28, 0xc0, 0xb1, 0x87,
	0xfd, 0xd3, 0x94, 0x52, 0xf2, 0x40, 0xba, 0x57, 0x98, 0xe9, 0x5f, 0xa4, 0x44, 0x69, 0xa6, 0x24,
	0xe4, 0x5a, 0x82, 0x31, 0x64, 0xc1, 0x6a, 0x32, 0x9b, 0x40, 0xaf, 0xbe, 0x19, 0x8f, 0xee, 0xeb,
	0xa5, 0x13, 0xb0, 0x10, 0x9f, 0xf1, 0x6e, 0x4e, 0x85, 0x7d, 0xd6, 0x8b, 0x84, 0x0d, 0xd8, 0xb0,
	0x37, 0xeb, 0x38, 0x39, 0x5d, 0xc4, 0x97, 0x9c, 0xe3, 0x12, 0x8c, 0x51, 0x99, 0x63, 0x07, 0x9e,
	0xf5, 0x6a, 0x67, 0xba, 0x88, 0xfb, 0xfc, 0x78, 0xa3, 0x5e, 0x4b, 0x65, 0x50, 0x25, 0x87, 0x03,
	0x36, 0x6c, 0xcf, 0x1a, 0xed, 0x18, 0x92, 0xb1, 0x05, 0xa0, 0x4d, 0xda, 0xfe, 0x63, 0xa3, 0xe3,
	0x6b, 0x7e, 0x82, 0x64, 0x8c, 0x42, 0x77, 0x95, 0xdb, 0x7c, 0xe4, 0x07, 0xa2, 0x7f, 0x33, 0x64,
	0x03, 0x22, 0x95, 0xc6, 0xdf, 0xd5, 0x09, 0xd9, 0xb5, 0x13, 0xb2, 0xff, 0x8a, 0x49, 0xba, 0x03,
	0x36, 0x8c, 0x66, 0x8d, 0xbe, 0x9b, 0x7c, 0xee, 0x05, 0xdb, 0xed, 0x05, 0xfb, 0xd9, 0x0b, 0xf6,
	0x51, 0x89, 0xd6, 0xae, 0x12, 0xad, 0xaf, 0x4a, 0xb4, 0x9e, 0x6e, 0x52, 0x6d, 0x97, 0xe5, 0x7c,
	0x8c, 0xb4, 0x96, 0x8f, 0x0a, 0xd6, 0xa3, 0x87, 0xba, 0x69, 0x2a, 0x94, 0x7c, 0x93, 0xb8, 0x1d,
	0xb9, 0xc6, 0xed, 0x7b, 0xae, 0x36, 0xf3, 0x8e, 0x6f, 0xeb, 0xf6, 0x37, 0x00, 0x00, 0xff, 0xff,
	0x6f, 0x8a, 0xcf, 0xa4, 0x8e, 0x01, 0x00, 0x00,
}

func (m *CallbackData) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CallbackData) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CallbackData) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.Callback) > 0 {
		i -= len(m.Callback)
		copy(dAtA[i:], m.Callback)
		i = encodeVarintCallback(dAtA, i, uint64(len(m.Callback)))
		i--
		dAtA[i] = 0x3a
	}
	if len(m.AccountId) > 0 {
		i -= len(m.AccountId)
		copy(dAtA[i:], m.AccountId)
		i = encodeVarintCallback(dAtA, i, uint64(len(m.AccountId)))
		i--
		dAtA[i] = 0x32
	}
	if len(m.ConnectionId) > 0 {
		i -= len(m.ConnectionId)
		copy(dAtA[i:], m.ConnectionId)
		i = encodeVarintCallback(dAtA, i, uint64(len(m.ConnectionId)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.Contract) > 0 {
		i -= len(m.Contract)
		copy(dAtA[i:], m.Contract)
		i = encodeVarintCallback(dAtA, i, uint64(len(m.Contract)))
		i--
		dAtA[i] = 0x22
	}
	if m.Sequence != 0 {
		i = encodeVarintCallback(dAtA, i, uint64(m.Sequence))
		i--
		dAtA[i] = 0x18
	}
	if len(m.ChannelId) > 0 {
		i -= len(m.ChannelId)
		copy(dAtA[i:], m.ChannelId)
		i = encodeVarintCallback(dAtA, i, uint64(len(m.ChannelId)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.PortId) > 0 {
		i -= len(m.PortId)
		copy(dAtA[i:], m.PortId)
		i = encodeVarintCallback(dAtA, i, uint64(len(m.PortId)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintCallback(dAtA []byte, offset int, v uint64) int {
	offset -= sovCallback(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *CallbackData) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.PortId)
	if l > 0 {
		n += 1 + l + sovCallback(uint64(l))
	}
	l = len(m.ChannelId)
	if l > 0 {
		n += 1 + l + sovCallback(uint64(l))
	}
	if m.Sequence != 0 {
		n += 1 + sovCallback(uint64(m.Sequence))
	}
	l = len(m.Contract)
	if l > 0 {
		n += 1 + l + sovCallback(uint64(l))
	}
	l = len(m.ConnectionId)
	if l > 0 {
		n += 1 + l + sovCallback(uint64(l))
	}
	l = len(m.AccountId)
	if l > 0 {
		n += 1 + l + sovCallback(uint64(l))
	}
	l = len(m.Callback)
	if l > 0 {
		n += 1 + l + sovCallback(uint64(l))
	}
	return n
}

func sovCallback(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozCallback(x uint64) (n int) {
	return sovCallback(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *CallbackData) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowCallback
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
			return fmt.Errorf("proto: CallbackData: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CallbackData: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field PortId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCallback
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
				return ErrInvalidLengthCallback
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCallback
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.PortId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ChannelId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCallback
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
				return ErrInvalidLengthCallback
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCallback
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ChannelId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Sequence", wireType)
			}
			m.Sequence = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCallback
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Sequence |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Contract", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCallback
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
				return ErrInvalidLengthCallback
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCallback
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Contract = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ConnectionId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCallback
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
				return ErrInvalidLengthCallback
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCallback
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ConnectionId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 6:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AccountId", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCallback
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
				return ErrInvalidLengthCallback
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthCallback
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.AccountId = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 7:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Callback", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowCallback
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
				return ErrInvalidLengthCallback
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthCallback
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.Callback = append(m.Callback[:0], dAtA[iNdEx:postIndex]...)
			if m.Callback == nil {
				m.Callback = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipCallback(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthCallback
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
func skipCallback(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowCallback
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
					return 0, ErrIntOverflowCallback
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
					return 0, ErrIntOverflowCallback
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
				return 0, ErrInvalidLengthCallback
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupCallback
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthCallback
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthCallback        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowCallback          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupCallback = fmt.Errorf("proto: unexpected end of group")
)
