// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ojo/airdrop/v1/airdrop.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
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

// Params defines the parameters for the airdrop module.
type Params struct {
	// The block at which all unclaimed AirdropAccounts will instead mint tokens
	// into the community pool. After this block, all unclaimed airdrop accounts
	// will no longer be able to be claimed.
	ExpiryBlock uint64 `protobuf:"varint,1,opt,name=expiry_block,json=expiryBlock,proto3" json:"expiry_block,omitempty"`
	// The percentage of the initial airdrop that users must delegate in order to
	// receive their second portion.
	// E.g., if we want to require users to stake their entire initial airdrop, this will be 1.
	// cosmos.base.v1beta1.Dec delegation_requirement = 1;
	DelegationRequirement *github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,2,opt,name=delegation_requirement,json=delegationRequirement,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"delegation_requirement,omitempty"`
	// The multiplier for the amount of tokens users will receive once they claim their airdrop.
	// E.g., if we want users to receive an equal second half, this will be 2.
	AirdropFactor *github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=airdrop_factor,json=airdropFactor,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"airdrop_factor,omitempty"`
}

func (m *Params) Reset()         { *m = Params{} }
func (m *Params) String() string { return proto.CompactTextString(m) }
func (*Params) ProtoMessage()    {}
func (*Params) Descriptor() ([]byte, []int) {
	return fileDescriptor_295074f7d14bf8dc, []int{0}
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

// AirDropAccount defines an account that was created at genesis with an initial airdrop.
type AirdropAccount struct {
	VestingEndTime int64 `protobuf:"varint,1,opt,name=vesting_end_time,json=vestingEndTime,proto3" json:"vesting_end_time,omitempty"`
	// The address of the account that was created at genesis with the initial airdrop.
	OriginAddress string `protobuf:"bytes,2,opt,name=origin_address,json=originAddress,proto3" json:"origin_address,omitempty"`
	// The amount of tokens that were airdropped to the genesis account.
	OriginAmount uint64 `protobuf:"varint,3,opt,name=origin_amount,json=originAmount,proto3" json:"origin_amount,omitempty"`
	// The address of the account that the user has claimed the 2nd half of their airdrop to.
	ClaimAddress string `protobuf:"bytes,4,opt,name=claim_address,json=claimAddress,proto3" json:"claim_address,omitempty"`
	// The amount of tokens claimed in the 2nd half of the airdrop.
	ClaimAmount uint64 `protobuf:"varint,5,opt,name=claim_amount,json=claimAmount,proto3" json:"claim_amount,omitempty"`
}

func (m *AirdropAccount) Reset()         { *m = AirdropAccount{} }
func (m *AirdropAccount) String() string { return proto.CompactTextString(m) }
func (*AirdropAccount) ProtoMessage()    {}
func (*AirdropAccount) Descriptor() ([]byte, []int) {
	return fileDescriptor_295074f7d14bf8dc, []int{1}
}
func (m *AirdropAccount) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *AirdropAccount) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_AirdropAccount.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *AirdropAccount) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AirdropAccount.Merge(m, src)
}
func (m *AirdropAccount) XXX_Size() int {
	return m.Size()
}
func (m *AirdropAccount) XXX_DiscardUnknown() {
	xxx_messageInfo_AirdropAccount.DiscardUnknown(m)
}

var xxx_messageInfo_AirdropAccount proto.InternalMessageInfo

func init() {
	proto.RegisterType((*Params)(nil), "ojo.airdrop.v1.Params")
	proto.RegisterType((*AirdropAccount)(nil), "ojo.airdrop.v1.AirdropAccount")
}

func init() { proto.RegisterFile("ojo/airdrop/v1/airdrop.proto", fileDescriptor_295074f7d14bf8dc) }

var fileDescriptor_295074f7d14bf8dc = []byte{
	// 423 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x92, 0x41, 0x6f, 0xd3, 0x30,
	0x14, 0xc7, 0x9b, 0xad, 0x4c, 0xc2, 0x24, 0x11, 0xb2, 0x06, 0x0a, 0x13, 0x0a, 0xdb, 0x90, 0x50,
	0x85, 0x94, 0x44, 0x13, 0x17, 0x2e, 0x08, 0xb5, 0x1a, 0x88, 0x23, 0x0a, 0x9c, 0xb8, 0x58, 0xa9,
	0x63, 0x82, 0xdb, 0xda, 0x2f, 0xd8, 0x6e, 0xd9, 0xae, 0x7c, 0x02, 0x3e, 0xcc, 0x3e, 0xc4, 0x8e,
	0xd3, 0x4e, 0x88, 0x03, 0x82, 0xf6, 0x73, 0x20, 0xa1, 0xd8, 0xce, 0xc6, 0x8d, 0xc3, 0x4e, 0x79,
	0xf9, 0xbf, 0x7f, 0xfe, 0xbf, 0xf8, 0xf9, 0xa1, 0x87, 0x30, 0x83, 0xa2, 0xe2, 0xaa, 0x56, 0xd0,
	0x16, 0xab, 0xa3, 0xbe, 0xcc, 0x5b, 0x05, 0x06, 0x70, 0x0c, 0x33, 0xc8, 0x7b, 0x69, 0x75, 0xb4,
	0xb7, 0xdb, 0x40, 0x03, 0xb6, 0x55, 0x74, 0x95, 0x73, 0xed, 0x3d, 0xa0, 0xa0, 0x05, 0x68, 0xe2,
	0x1a, 0xee, 0xc5, 0xb5, 0x0e, 0xff, 0x04, 0x68, 0xe7, 0x6d, 0xa5, 0x2a, 0xa1, 0xf1, 0x01, 0x0a,
	0xd9, 0x49, 0xcb, 0xd5, 0x29, 0x99, 0x2e, 0x80, 0xce, 0x93, 0x60, 0x3f, 0x18, 0x0d, 0xcb, 0x3b,
	0x4e, 0x9b, 0x74, 0x12, 0x06, 0x74, 0xbf, 0x66, 0x0b, 0xd6, 0x54, 0x86, 0x83, 0x24, 0x8a, 0x7d,
	0x5e, 0x72, 0xc5, 0x04, 0x93, 0x26, 0xd9, 0xda, 0x0f, 0x46, 0xb7, 0x27, 0xcf, 0x7f, 0xfc, 0x7c,
	0xf4, 0xa4, 0xe1, 0xe6, 0xd3, 0x72, 0x9a, 0x53, 0x10, 0x1e, 0xe5, 0x1f, 0x99, 0xae, 0xe7, 0x85,
	0x39, 0x6d, 0x99, 0xce, 0x8f, 0x19, 0xbd, 0x3c, 0xcb, 0x90, 0xff, 0x93, 0x63, 0x46, 0xcb, 0x7b,
	0xd7, 0xb9, 0xe5, 0x75, 0x2c, 0x26, 0x28, 0xf6, 0xa7, 0x23, 0x1f, 0x2b, 0x6a, 0x40, 0x25, 0xdb,
	0x37, 0x04, 0x45, 0x3e, 0xef, 0xb5, 0x8d, 0x3b, 0xfc, 0xba, 0x85, 0xe2, 0xb1, 0x53, 0xc6, 0x94,
	0xc2, 0x52, 0x1a, 0x3c, 0x42, 0x77, 0x57, 0x4c, 0x1b, 0x2e, 0x1b, 0xc2, 0x64, 0x4d, 0x0c, 0x17,
	0xcc, 0xce, 0x62, 0xbb, 0x8c, 0xbd, 0xfe, 0x4a, 0xd6, 0xef, 0xb9, 0x60, 0xf8, 0x25, 0x8a, 0x41,
	0xf1, 0x86, 0x4b, 0x52, 0xd5, 0xb5, 0x62, 0x5a, 0xfb, 0x31, 0x24, 0x97, 0x67, 0xd9, 0xae, 0x67,
	0x8e, 0x5d, 0xe7, 0x9d, 0x51, 0x5c, 0x36, 0x65, 0xe4, 0xfc, 0x5e, 0xc4, 0x8f, 0x51, 0xd4, 0x07,
	0x88, 0x8e, 0x6d, 0x4f, 0x37, 0x2c, 0x43, 0xef, 0xb2, 0x1a, 0x7e, 0x81, 0x22, 0xba, 0xa8, 0xb8,
	0xb8, 0x82, 0x0c, 0xff, 0x03, 0x09, 0xad, 0xbd, 0x67, 0x1c, 0xa0, 0xd0, 0x7f, 0xee, 0x10, 0xb7,
	0xdc, 0xb5, 0x3a, 0x8f, 0x95, 0x26, 0x6f, 0xce, 0x7f, 0xa7, 0x83, 0xf3, 0x75, 0x1a, 0x5c, 0xac,
	0xd3, 0xe0, 0xd7, 0x3a, 0x0d, 0xbe, 0x6d, 0xd2, 0xc1, 0xc5, 0x26, 0x1d, 0x7c, 0xdf, 0xa4, 0x83,
	0x0f, 0x4f, 0xff, 0x99, 0x33, 0xcc, 0x20, 0x93, 0xcc, 0x7c, 0x01, 0x35, 0xef, 0xea, 0xe2, 0xe4,
	0x6a, 0x35, 0xed, 0xbc, 0xa7, 0x3b, 0x76, 0xab, 0x9e, 0xfd, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x2f,
	0xcd, 0x05, 0x00, 0xb6, 0x02, 0x00, 0x00,
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
	if m.AirdropFactor != nil {
		{
			size := m.AirdropFactor.Size()
			i -= size
			if _, err := m.AirdropFactor.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintAirdrop(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x1a
	}
	if m.DelegationRequirement != nil {
		{
			size := m.DelegationRequirement.Size()
			i -= size
			if _, err := m.DelegationRequirement.MarshalTo(dAtA[i:]); err != nil {
				return 0, err
			}
			i = encodeVarintAirdrop(dAtA, i, uint64(size))
		}
		i--
		dAtA[i] = 0x12
	}
	if m.ExpiryBlock != 0 {
		i = encodeVarintAirdrop(dAtA, i, uint64(m.ExpiryBlock))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func (m *AirdropAccount) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *AirdropAccount) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *AirdropAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.ClaimAmount != 0 {
		i = encodeVarintAirdrop(dAtA, i, uint64(m.ClaimAmount))
		i--
		dAtA[i] = 0x28
	}
	if len(m.ClaimAddress) > 0 {
		i -= len(m.ClaimAddress)
		copy(dAtA[i:], m.ClaimAddress)
		i = encodeVarintAirdrop(dAtA, i, uint64(len(m.ClaimAddress)))
		i--
		dAtA[i] = 0x22
	}
	if m.OriginAmount != 0 {
		i = encodeVarintAirdrop(dAtA, i, uint64(m.OriginAmount))
		i--
		dAtA[i] = 0x18
	}
	if len(m.OriginAddress) > 0 {
		i -= len(m.OriginAddress)
		copy(dAtA[i:], m.OriginAddress)
		i = encodeVarintAirdrop(dAtA, i, uint64(len(m.OriginAddress)))
		i--
		dAtA[i] = 0x12
	}
	if m.VestingEndTime != 0 {
		i = encodeVarintAirdrop(dAtA, i, uint64(m.VestingEndTime))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintAirdrop(dAtA []byte, offset int, v uint64) int {
	offset -= sovAirdrop(v)
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
	if m.ExpiryBlock != 0 {
		n += 1 + sovAirdrop(uint64(m.ExpiryBlock))
	}
	if m.DelegationRequirement != nil {
		l = m.DelegationRequirement.Size()
		n += 1 + l + sovAirdrop(uint64(l))
	}
	if m.AirdropFactor != nil {
		l = m.AirdropFactor.Size()
		n += 1 + l + sovAirdrop(uint64(l))
	}
	return n
}

func (m *AirdropAccount) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.VestingEndTime != 0 {
		n += 1 + sovAirdrop(uint64(m.VestingEndTime))
	}
	l = len(m.OriginAddress)
	if l > 0 {
		n += 1 + l + sovAirdrop(uint64(l))
	}
	if m.OriginAmount != 0 {
		n += 1 + sovAirdrop(uint64(m.OriginAmount))
	}
	l = len(m.ClaimAddress)
	if l > 0 {
		n += 1 + l + sovAirdrop(uint64(l))
	}
	if m.ClaimAmount != 0 {
		n += 1 + sovAirdrop(uint64(m.ClaimAmount))
	}
	return n
}

func sovAirdrop(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozAirdrop(x uint64) (n int) {
	return sovAirdrop(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *Params) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAirdrop
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
				return fmt.Errorf("proto: wrong wireType = %d for field ExpiryBlock", wireType)
			}
			m.ExpiryBlock = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAirdrop
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ExpiryBlock |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field DelegationRequirement", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAirdrop
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
				return ErrInvalidLengthAirdrop
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAirdrop
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Dec
			m.DelegationRequirement = &v
			if err := m.DelegationRequirement.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field AirdropFactor", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAirdrop
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
				return ErrInvalidLengthAirdrop
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAirdrop
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			var v github_com_cosmos_cosmos_sdk_types.Dec
			m.AirdropFactor = &v
			if err := m.AirdropFactor.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
				return err
			}
			iNdEx = postIndex
		default:
			iNdEx = preIndex
			skippy, err := skipAirdrop(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAirdrop
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
func (m *AirdropAccount) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowAirdrop
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
			return fmt.Errorf("proto: AirdropAccount: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: AirdropAccount: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field VestingEndTime", wireType)
			}
			m.VestingEndTime = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAirdrop
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.VestingEndTime |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field OriginAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAirdrop
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
				return ErrInvalidLengthAirdrop
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAirdrop
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.OriginAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field OriginAmount", wireType)
			}
			m.OriginAmount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAirdrop
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.OriginAmount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClaimAddress", wireType)
			}
			var stringLen uint64
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAirdrop
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
				return ErrInvalidLengthAirdrop
			}
			postIndex := iNdEx + intStringLen
			if postIndex < 0 {
				return ErrInvalidLengthAirdrop
			}
			if postIndex > l {
				return io.ErrUnexpectedEOF
			}
			m.ClaimAddress = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field ClaimAmount", wireType)
			}
			m.ClaimAmount = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAirdrop
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.ClaimAmount |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipAirdrop(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthAirdrop
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
func skipAirdrop(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowAirdrop
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
					return 0, ErrIntOverflowAirdrop
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
					return 0, ErrIntOverflowAirdrop
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
				return 0, ErrInvalidLengthAirdrop
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupAirdrop
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthAirdrop
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthAirdrop        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowAirdrop          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupAirdrop = fmt.Errorf("proto: unexpected end of group")
)
