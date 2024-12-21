// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ojo/symbiotic/v1/symbiotic.proto

package types

import (
	fmt "fmt"
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

// Params defines the parameters for the symbiotic module.
type Params struct {
	// address of the ojo middleware contract
	MiddlewareAddress string `protobuf:"bytes,1,opt,name=middleware_address,json=middlewareAddress,proto3" json:"middleware_address,omitempty"`
	// block period for syncing with the symbiotic network on Ethereum
	SymbioticSyncPeriod int64 `protobuf:"varint,2,opt,name=symbiotic_sync_period,json=symbioticSyncPeriod,proto3" json:"symbiotic_sync_period,omitempty"`
	// maximum amount of cached block hashes in the store
	MaximumCachedBlockHashes uint64 `protobuf:"varint,3,opt,name=maximum_cached_block_hashes,json=maximumCachedBlockHashes,proto3" json:"maximum_cached_block_hashes,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_4008c10d1c664ed5, []int{0}
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

// Cached block hash and height of block hash from the chain the ojo middleware is on.
type CachedBlockHash struct {
	// Block hash of cached block on chain.
	BlockHash string `protobuf:"bytes,1,opt,name=block_hash,json=blockHash,proto3" json:"block_hash,omitempty"`
	// Block height of block hash.
	Height int64 `protobuf:"varint,2,opt,name=height,proto3" json:"height,omitempty"`
}

func (m *CachedBlockHash) Reset()         { *m = CachedBlockHash{} }
func (m *CachedBlockHash) String() string { return proto.CompactTextString(m) }
func (*CachedBlockHash) ProtoMessage()    {}
func (*CachedBlockHash) Descriptor() ([]byte, []int) {
	return fileDescriptor_4008c10d1c664ed5, []int{1}
}
func (m *CachedBlockHash) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *CachedBlockHash) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_CachedBlockHash.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *CachedBlockHash) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CachedBlockHash.Merge(m, src)
}
func (m *CachedBlockHash) XXX_Size() int {
	return m.Size()
}
func (m *CachedBlockHash) XXX_DiscardUnknown() {
	xxx_messageInfo_CachedBlockHash.DiscardUnknown(m)
}

var xxx_messageInfo_CachedBlockHash proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Params)(nil), "ojo.symbiotic.v1.Params")
	proto.RegisterType((*CachedBlockHash)(nil), "ojo.symbiotic.v1.CachedBlockHash")
}

func init() { proto.RegisterFile("ojo/symbiotic/v1/symbiotic.proto", fileDescriptor_4008c10d1c664ed5) }

var fileDescriptor_4008c10d1c664ed5 = []byte{
	// 302 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x5c, 0x90, 0xcd, 0x4e, 0x32, 0x31,
	0x14, 0x86, 0xa7, 0x1f, 0x5f, 0x48, 0xe8, 0x46, 0x1d, 0x7f, 0x32, 0xd1, 0xd8, 0x4c, 0x58, 0xb1,
	0x10, 0x26, 0xe8, 0xda, 0x85, 0xb8, 0x21, 0xae, 0x08, 0xee, 0xdc, 0x34, 0x9d, 0xb6, 0x99, 0x16,
	0x28, 0x87, 0xb4, 0xe5, 0x67, 0xee, 0xc2, 0xeb, 0xf0, 0x4a, 0x58, 0xb2, 0x74, 0xa9, 0x70, 0x23,
	0x86, 0x61, 0xc2, 0x18, 0x77, 0xe7, 0x9c, 0xe7, 0xc9, 0xc9, 0x9b, 0x17, 0xc7, 0x30, 0x82, 0xc4,
	0xe5, 0x26, 0xd5, 0xe0, 0x35, 0x4f, 0x16, 0xdd, 0x6a, 0xe9, 0xcc, 0x2c, 0x78, 0x08, 0x4f, 0x61,
	0x04, 0x9d, 0xea, 0xb8, 0xe8, 0x5e, 0x5f, 0x64, 0x90, 0x41, 0x01, 0x93, 0xfd, 0x74, 0xf0, 0x9a,
	0x1f, 0x08, 0xd7, 0x07, 0xcc, 0x32, 0xe3, 0xc2, 0x36, 0x0e, 0x8d, 0x16, 0x62, 0x22, 0x97, 0xcc,
	0x4a, 0xca, 0x84, 0xb0, 0xd2, 0xb9, 0x08, 0xc5, 0xa8, 0xd5, 0x18, 0x9e, 0x55, 0xe4, 0xe9, 0x00,
	0xc2, 0x7b, 0x7c, 0x79, 0xfc, 0x4f, 0x5d, 0x3e, 0xe5, 0x74, 0x26, 0xad, 0x06, 0x11, 0xfd, 0x8b,
	0x51, 0xab, 0x36, 0x3c, 0x3f, 0xc2, 0xd7, 0x7c, 0xca, 0x07, 0x05, 0x0a, 0x1f, 0xf1, 0x8d, 0x61,
	0x2b, 0x6d, 0xe6, 0x86, 0x72, 0xc6, 0x95, 0x14, 0x34, 0x9d, 0x00, 0x1f, 0x53, 0xc5, 0x9c, 0x92,
	0x2e, 0xaa, 0xc5, 0xa8, 0xf5, 0x7f, 0x18, 0x95, 0xca, 0x73, 0x61, 0xf4, 0xf6, 0x42, 0xbf, 0xe0,
	0xcd, 0x3e, 0x3e, 0xf9, 0x73, 0x0c, 0x6f, 0x31, 0xae, 0x5e, 0x94, 0x61, 0x1b, 0xe9, 0x11, 0x5f,
	0xe1, 0xba, 0x92, 0x3a, 0x53, 0xbe, 0x4c, 0x55, 0x6e, 0xbd, 0x97, 0xf5, 0x37, 0x09, 0xd6, 0x5b,
	0x82, 0x36, 0x5b, 0x82, 0xbe, 0xb6, 0x04, 0xbd, 0xef, 0x48, 0xb0, 0xd9, 0x91, 0xe0, 0x73, 0x47,
	0x82, 0xb7, 0xbb, 0x4c, 0x7b, 0x35, 0x4f, 0x3b, 0x1c, 0x4c, 0x02, 0x23, 0x68, 0x4f, 0xa5, 0x5f,
	0x82, 0x1d, 0xef, 0xe7, 0x64, 0xf5, 0xab, 0x77, 0x9f, 0xcf, 0xa4, 0x4b, 0xeb, 0x45, 0x93, 0x0f,
	0x3f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x01, 0x76, 0x64, 0x1d, 0x95, 0x01, 0x00, 0x00,
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
	if m.MaximumCachedBlockHashes != 0 {
		i = encodeVarintSymbiotic(dAtA, i, uint64(m.MaximumCachedBlockHashes))
		i--
		dAtA[i] = 0x18
	}
	if m.SymbioticSyncPeriod != 0 {
		i = encodeVarintSymbiotic(dAtA, i, uint64(m.SymbioticSyncPeriod))
		i--
		dAtA[i] = 0x10
	}
	if len(m.MiddlewareAddress) > 0 {
		i -= len(m.MiddlewareAddress)
		copy(dAtA[i:], m.MiddlewareAddress)
		i = encodeVarintSymbiotic(dAtA, i, uint64(len(m.MiddlewareAddress)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *CachedBlockHash) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *CachedBlockHash) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *CachedBlockHash) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.Height != 0 {
		i = encodeVarintSymbiotic(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x10
	}
	if len(m.BlockHash) > 0 {
		i -= len(m.BlockHash)
		copy(dAtA[i:], m.BlockHash)
		i = encodeVarintSymbiotic(dAtA, i, uint64(len(m.BlockHash)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func encodeVarintSymbiotic(dAtA []byte, offset int, v uint64) int {
	offset -= sovSymbiotic(v)
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
	l = len(m.MiddlewareAddress)
	if l > 0 {
		n += 1 + l + sovSymbiotic(uint64(l))
	}
	if m.SymbioticSyncPeriod != 0 {
		n += 1 + sovSymbiotic(uint64(m.SymbioticSyncPeriod))
	}
	if m.MaximumCachedBlockHashes != 0 {
		n += 1 + sovSymbiotic(uint64(m.MaximumCachedBlockHashes))
	}
	return n
}

func (m *CachedBlockHash) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.BlockHash)
	if l > 0 {
		n += 1 + l + sovSymbiotic(uint64(l))
	}
	if m.Height != 0 {
		n += 1 + sovSymbiotic(uint64(m.Height))
	}
	return n
}

func sovSymbiotic(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozSymbiotic(x uint64) (n int) {
	return sovSymbiotic(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSymbiotic
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
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field MiddlewareAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSymbiotic
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
				return ErrInvalidLengthSymbiotic
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSymbiotic
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.MiddlewareAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field SymbioticSyncPeriod", wireType)
			}
			m.SymbioticSyncPeriod = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSymbiotic
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.SymbioticSyncPeriod |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MaximumCachedBlockHashes", wireType)
			}
			m.MaximumCachedBlockHashes = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSymbiotic
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MaximumCachedBlockHashes |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipSymbiotic(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSymbiotic
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
func (m *CachedBlockHash) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowSymbiotic
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
			return fmt.Errorf("proto: CachedBlockHash: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: CachedBlockHash: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field BlockHash", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSymbiotic
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
				return ErrInvalidLengthSymbiotic
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthSymbiotic
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.BlockHash = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowSymbiotic
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
		default:
			iNdEx = preIndex
			skippy, err := skipSymbiotic(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthSymbiotic
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
func skipSymbiotic(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowSymbiotic
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
					return 0, ErrIntOverflowSymbiotic
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
					return 0, ErrIntOverflowSymbiotic
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
				return 0, ErrInvalidLengthSymbiotic
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupSymbiotic
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthSymbiotic
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthSymbiotic        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowSymbiotic          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupSymbiotic = fmt.Errorf("proto: unexpected end of group")
)