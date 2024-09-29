// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ojo/gmp/v1/abci.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	proto "github.com/cosmos/gogoproto/proto"
	_ "github.com/gogo/protobuf/gogoproto"
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

// GmpVoteExtension defines the vote extension structure used by the gmp
// module.
type GmpVoteExtension struct {
	Height        int64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	GasEstimation int64 `protobuf:"varint,2,opt,name=gas_estimation,json=gasEstimation,proto3" json:"gas_estimation,omitempty"`
}

func (m *GmpVoteExtension) Reset()         { *m = GmpVoteExtension{} }
func (m *GmpVoteExtension) String() string { return proto.CompactTextString(m) }
func (*GmpVoteExtension) ProtoMessage()    {}
func (*GmpVoteExtension) Descriptor() ([]byte, []int) {
	return fileDescriptor_85641ed06050038a, []int{0}
}
func (m *GmpVoteExtension) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *GmpVoteExtension) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_GmpVoteExtension.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *GmpVoteExtension) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GmpVoteExtension.Merge(m, src)
}
func (m *GmpVoteExtension) XXX_Size() int {
	return m.Size()
}
func (m *GmpVoteExtension) XXX_DiscardUnknown() {
	xxx_messageInfo_GmpVoteExtension.DiscardUnknown(m)
}

var xxx_messageInfo_GmpVoteExtension proto.InternalMessageInfo

// InjectedVoteExtensionTx defines the vote extension tx injected by the prepare
// proposal handler.
type InjectedVoteExtensionTx struct {
	MedianGasEstimation int64  `protobuf:"varint,1,opt,name=median_gas_estimation,json=medianGasEstimation,proto3" json:"median_gas_estimation,omitempty"`
	ExtendedCommitInfo  []byte `protobuf:"bytes,2,opt,name=extended_commit_info,json=extendedCommitInfo,proto3" json:"extended_commit_info,omitempty"`
}

func (m *InjectedVoteExtensionTx) Reset()         { *m = InjectedVoteExtensionTx{} }
func (m *InjectedVoteExtensionTx) String() string { return proto.CompactTextString(m) }
func (*InjectedVoteExtensionTx) ProtoMessage()    {}
func (*InjectedVoteExtensionTx) Descriptor() ([]byte, []int) {
	return fileDescriptor_85641ed06050038a, []int{1}
}
func (m *InjectedVoteExtensionTx) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *InjectedVoteExtensionTx) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_InjectedVoteExtensionTx.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *InjectedVoteExtensionTx) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InjectedVoteExtensionTx.Merge(m, src)
}
func (m *InjectedVoteExtensionTx) XXX_Size() int {
	return m.Size()
}
func (m *InjectedVoteExtensionTx) XXX_DiscardUnknown() {
	xxx_messageInfo_InjectedVoteExtensionTx.DiscardUnknown(m)
}

var xxx_messageInfo_InjectedVoteExtensionTx proto.InternalMessageInfo

func init() {
	proto.RegisterType((*GmpVoteExtension)(nil), "ojo.gmp.v1.GmpVoteExtension")
	proto.RegisterType((*InjectedVoteExtensionTx)(nil), "ojo.gmp.v1.InjectedVoteExtensionTx")
}

func init() { proto.RegisterFile("ojo/gmp/v1/abci.proto", fileDescriptor_85641ed06050038a) }

var fileDescriptor_85641ed06050038a = []byte{
	// 292 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xb1, 0x4e, 0xf3, 0x30,
	0x10, 0xc7, 0x93, 0xef, 0x93, 0x3a, 0x58, 0x80, 0x90, 0x69, 0xa1, 0x74, 0xb0, 0x50, 0x25, 0x10,
	0x0b, 0x31, 0x85, 0x37, 0x00, 0xaa, 0xaa, 0x23, 0x15, 0x62, 0x60, 0x89, 0xd2, 0xc4, 0x75, 0x1d,
	0x64, 0x5f, 0x54, 0x1f, 0x25, 0x4c, 0xbc, 0x02, 0x8f, 0xd5, 0xb1, 0x23, 0x23, 0x24, 0x2f, 0x82,
	0x62, 0x03, 0x02, 0xb6, 0xbb, 0xfb, 0x9d, 0x7e, 0x7f, 0xe9, 0x4f, 0x3a, 0x90, 0x03, 0x97, 0xba,
	0xe0, 0xcb, 0x01, 0x4f, 0xa6, 0xa9, 0x8a, 0x8a, 0x05, 0x20, 0x50, 0x02, 0x39, 0x44, 0x52, 0x17,
	0xd1, 0x72, 0xd0, 0x6b, 0x4b, 0x90, 0xe0, 0xce, 0xbc, 0x99, 0xfc, 0x47, 0x6f, 0x3f, 0x05, 0xab,
	0xc1, 0xc6, 0x1e, 0xf8, 0xc5, 0xa3, 0xfe, 0x35, 0xd9, 0x1e, 0xe9, 0xe2, 0x16, 0x50, 0x0c, 0x4b,
	0x14, 0xc6, 0x2a, 0x30, 0x74, 0x97, 0xb4, 0xe6, 0x42, 0xc9, 0x39, 0x76, 0xc3, 0x83, 0xf0, 0xf8,
	0xff, 0xe4, 0x73, 0xa3, 0x87, 0x64, 0x4b, 0x26, 0x36, 0x16, 0x16, 0x95, 0x4e, 0x50, 0x81, 0xe9,
	0xfe, 0x73, 0x7c, 0x53, 0x26, 0x76, 0xf8, 0x7d, 0xec, 0x3f, 0x93, 0xbd, 0xb1, 0xc9, 0x45, 0x8a,
	0x22, 0xfb, 0xe5, 0xbd, 0x29, 0xe9, 0x19, 0xe9, 0x68, 0x91, 0xa9, 0xc4, 0xc4, 0x7f, 0x44, 0x3e,
	0x68, 0xc7, 0xc3, 0xd1, 0x4f, 0x1d, 0x3d, 0x25, 0x6d, 0xd1, 0x28, 0x32, 0x91, 0xc5, 0x29, 0x68,
	0xad, 0x30, 0x56, 0x66, 0x06, 0x2e, 0x7b, 0x63, 0x42, 0xbf, 0xd8, 0xa5, 0x43, 0x63, 0x33, 0x83,
	0x8b, 0xab, 0xd5, 0x3b, 0x0b, 0x56, 0x15, 0x0b, 0xd7, 0x15, 0x0b, 0xdf, 0x2a, 0x16, 0xbe, 0xd4,
	0x2c, 0x58, 0xd7, 0x2c, 0x78, 0xad, 0x59, 0x70, 0x77, 0x24, 0x15, 0xce, 0x1f, 0xa6, 0x51, 0x0a,
	0x9a, 0x43, 0x0e, 0x27, 0x46, 0xe0, 0x23, 0x2c, 0xee, 0x9b, 0x99, 0x97, 0xae, 0x5e, 0x7c, 0x2a,
	0x84, 0x9d, 0xb6, 0x5c, 0x41, 0xe7, 0x1f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xc9, 0xdc, 0xdc, 0xa9,
	0x76, 0x01, 0x00, 0x00,
}

func (m *GmpVoteExtension) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *GmpVoteExtension) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *GmpVoteExtension) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.GasEstimation != 0 {
		i = encodeVarintAbci(dAtA, i, uint64(m.GasEstimation))
		i--
		dAtA[i] = 0x10
	}
	if m.Height != 0 {
		i = encodeVarintAbci(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *InjectedVoteExtensionTx) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *InjectedVoteExtensionTx) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *InjectedVoteExtensionTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if len(m.ExtendedCommitInfo) > 0 {
		i -= len(m.ExtendedCommitInfo)
		copy(dAtA[i:], m.ExtendedCommitInfo)
		i = encodeVarintAbci(dAtA, i, uint64(len(m.ExtendedCommitInfo)))
		i--
		dAtA[i] = 0x12
	}
	if m.MedianGasEstimation != 0 {
		i = encodeVarintAbci(dAtA, i, uint64(m.MedianGasEstimation))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintAbci(dAtA []byte, offset int, v uint64) int {
	offset -= sovAbci(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *GmpVoteExtension) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovAbci(uint64(m.Height))
	}
	if m.GasEstimation != 0 {
		n += 1 + sovAbci(uint64(m.GasEstimation))
	}
	return n
}

func (m *InjectedVoteExtensionTx) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.MedianGasEstimation != 0 {
		n += 1 + sovAbci(uint64(m.MedianGasEstimation))
	}
	l = len(m.ExtendedCommitInfo)
	if l > 0 {
		n += 1 + l + sovAbci(uint64(l))
	}
	return n
}

func sovAbci(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAbci(x uint64) (n int) {
	return sovAbci(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *GmpVoteExtension) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAbci
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
			return fmt.Errorf("proto: GmpVoteExtension: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: GmpVoteExtension: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAbci
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasEstimation", wireType)
			}
			m.GasEstimation = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAbci
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GasEstimation |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipAbci(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAbci
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
func (m *InjectedVoteExtensionTx) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAbci
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
			return fmt.Errorf("proto: InjectedVoteExtensionTx: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: InjectedVoteExtensionTx: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MedianGasEstimation", wireType)
			}
			m.MedianGasEstimation = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAbci
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MedianGasEstimation |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ExtendedCommitInfo", wireType)
			}
			var byteLen int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAbci
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
				return ErrInvalidLengthAbci
			}
			postIndex := iNdEx + byteLen
			if postIndex < 0 {
				return ErrInvalidLengthAbci
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ExtendedCommitInfo = append(m.ExtendedCommitInfo[:0], dAtA[iNdEx:postIndex]...)
			if m.ExtendedCommitInfo == nil {
				m.ExtendedCommitInfo = []byte{}
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAbci(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAbci
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
func skipAbci(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAbci
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
					return 0, ErrIntOverflowAbci
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
					return 0, ErrIntOverflowAbci
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
				return 0, ErrInvalidLengthAbci
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAbci
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAbci
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAbci        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAbci          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAbci = fmt.Errorf("proto: unexpected end of group")
)
