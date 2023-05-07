// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ojo/airdrop/v1/airdrop.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	github_com_cosmos_cosmos_sdk_types "github.com/cosmos/cosmos-sdk/types"
	_ "github.com/gogo/protobuf/gogoproto"
	proto "github.com/gogo/protobuf/proto"
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
	// Flag to determine if the origin vesting accounts have been created yet
	OriginAccountsCreated bool `protobuf:"varint,1,opt,name=origin_accounts_created,json=originAccountsCreated,proto3" json:"origin_accounts_created,omitempty"`
	// The block at which all unclaimed AirdropAccounts will instead mint tokens
	// into the community pool. After this block, all unclaimed airdrop accounts
	// will no longer be able to be claimed.
	ExpiryBlock uint64 `protobuf:"varint,2,opt,name=expiry_block,json=expiryBlock,proto3" json:"expiry_block,omitempty"`
	// The percentage of the initial airdrop that users must delegate in order to
	// receive their second portion.
	// E.g., if we want to require users to stake their entire initial airdrop, this will be 1.
	// cosmos.base.v1beta1.Dec delegation_requirement = 1;
	DelegationRequirement *github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,3,opt,name=delegation_requirement,json=delegationRequirement,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"delegation_requirement,omitempty"`
	// The multiplier for the amount of tokens users will receive once they claim their airdrop.
	// E.g., if we want users to receive an equal second half, this will be 2.
	AirdropFactor *github_com_cosmos_cosmos_sdk_types.Dec `protobuf:"bytes,4,opt,name=airdrop_factor,json=airdropFactor,proto3,customtype=github.com/cosmos/cosmos-sdk/types.Dec" json:"airdrop_factor,omitempty"`
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
	// The current state of the airdrop account
	State string `protobuf:"bytes,4,opt,name=state,proto3" json:"state,omitempty"`
	// The address of the account that the user has claimed the 2nd half of their airdrop to.
	ClaimAddress string `protobuf:"bytes,5,opt,name=claim_address,json=claimAddress,proto3" json:"claim_address,omitempty"`
	// The amount of tokens claimed in the 2nd half of the airdrop.
	ClaimAmount uint64 `protobuf:"varint,6,opt,name=claim_amount,json=claimAmount,proto3" json:"claim_amount,omitempty"`
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
	// 461 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x52, 0x41, 0x6f, 0xd3, 0x30,
	0x18, 0x6d, 0xba, 0xae, 0x62, 0xa6, 0x8d, 0x90, 0xd5, 0x41, 0x98, 0x50, 0xe8, 0x86, 0x84, 0x2a,
	0xa4, 0x36, 0x9a, 0x90, 0x10, 0x17, 0x84, 0x5a, 0x06, 0xe2, 0x88, 0x02, 0x27, 0x2e, 0x91, 0xeb,
	0x7c, 0x04, 0xb7, 0x8d, 0xbf, 0x62, 0xbb, 0x65, 0xfb, 0x17, 0x1c, 0xf9, 0x21, 0xe3, 0x3f, 0xec,
	0x38, 0xed, 0x84, 0x38, 0x20, 0x68, 0xff, 0x08, 0xaa, 0xed, 0xac, 0xdc, 0x38, 0x70, 0xca, 0xe7,
	0xf7, 0x5e, 0xde, 0xb3, 0x9f, 0x4d, 0xee, 0xe1, 0x04, 0x13, 0x26, 0x54, 0xae, 0x70, 0x9e, 0x2c,
	0x8f, 0xab, 0x71, 0x30, 0x57, 0x68, 0x90, 0x86, 0x38, 0xc1, 0x41, 0x05, 0x2d, 0x8f, 0x0f, 0x3a,
	0x05, 0x16, 0x68, 0xa9, 0x64, 0x33, 0x39, 0xd5, 0xc1, 0x5d, 0x8e, 0xba, 0x44, 0x9d, 0x39, 0xc2,
	0x2d, 0x1c, 0x75, 0xf4, 0xad, 0x4e, 0x9a, 0x6f, 0x98, 0x62, 0xa5, 0xa6, 0x4f, 0xc8, 0x1d, 0x54,
	0xa2, 0x10, 0x32, 0x63, 0x9c, 0xe3, 0x42, 0x1a, 0x9d, 0x71, 0x05, 0xcc, 0x40, 0x1e, 0x05, 0xdd,
	0xa0, 0x77, 0x23, 0xdd, 0x77, 0xf4, 0xd0, 0xb3, 0x2f, 0x1c, 0x49, 0x0f, 0x49, 0x0b, 0x4e, 0xe7,
	0x42, 0x9d, 0x65, 0xe3, 0x19, 0xf2, 0x69, 0x54, 0xef, 0x06, 0xbd, 0x46, 0x7a, 0xd3, 0x61, 0xa3,
	0x0d, 0x44, 0x91, 0xdc, 0xce, 0x61, 0x06, 0x05, 0x33, 0x02, 0x65, 0xa6, 0xe0, 0xd3, 0x42, 0x28,
	0x28, 0x41, 0x9a, 0x68, 0xa7, 0x1b, 0xf4, 0xf6, 0x46, 0x4f, 0x7f, 0xfc, 0xbc, 0xff, 0xb0, 0x10,
	0xe6, 0xe3, 0x62, 0x3c, 0xe0, 0x58, 0xfa, 0x2d, 0xfa, 0x4f, 0x5f, 0xe7, 0xd3, 0xc4, 0x9c, 0xcd,
	0x41, 0x0f, 0x4e, 0x80, 0x5f, 0x9d, 0xf7, 0x89, 0x3f, 0xc1, 0x09, 0xf0, 0x74, 0x7f, 0xeb, 0x9b,
	0x6e, 0x6d, 0x69, 0x46, 0x42, 0xdf, 0x4a, 0xf6, 0x81, 0x71, 0x83, 0x2a, 0x6a, 0xfc, 0x67, 0x50,
	0xdb, 0xfb, 0xbd, 0xb2, 0x76, 0x47, 0x5f, 0xeb, 0x24, 0x1c, 0x3a, 0xc4, 0xf7, 0x41, 0x7b, 0xe4,
	0xd6, 0x12, 0xb4, 0x11, 0xb2, 0xc8, 0x40, 0xe6, 0x99, 0x11, 0x25, 0xd8, 0xe2, 0x76, 0xd2, 0xd0,
	0xe3, 0x2f, 0x65, 0xfe, 0x4e, 0x94, 0x40, 0x9f, 0x93, 0xb0, 0x6a, 0x3a, 0xcf, 0x15, 0x68, 0x6d,
	0x3b, 0xdb, 0x1b, 0x45, 0x57, 0xe7, 0xfd, 0x8e, 0xcf, 0x1c, 0x3a, 0xe6, 0xad, 0x51, 0x42, 0x16,
	0x69, 0xdb, 0x57, 0xef, 0x40, 0xfa, 0x80, 0xb4, 0x2b, 0x83, 0x72, 0x93, 0x6d, 0x6b, 0x6c, 0xa4,
	0x2d, 0xaf, 0xb2, 0x18, 0xed, 0x90, 0x5d, 0x6d, 0x98, 0x01, 0x77, 0xf4, 0xd4, 0x2d, 0xe8, 0x33,
	0xd2, 0xe6, 0x33, 0x26, 0xca, 0xeb, 0xe8, 0xdd, 0x7f, 0x44, 0xb7, 0xac, 0xbc, 0x4a, 0x3e, 0x24,
	0x2d, 0xff, 0xbb, 0x0b, 0x6e, 0xba, 0xcb, 0x76, 0x1a, 0x0b, 0x8d, 0x5e, 0x5f, 0xfc, 0x8e, 0x6b,
	0x17, 0xab, 0x38, 0xb8, 0x5c, 0xc5, 0xc1, 0xaf, 0x55, 0x1c, 0x7c, 0x59, 0xc7, 0xb5, 0xcb, 0x75,
	0x5c, 0xfb, 0xbe, 0x8e, 0x6b, 0xef, 0x1f, 0xfd, 0xd5, 0x3e, 0x4e, 0xb0, 0x2f, 0xc1, 0x7c, 0x46,
	0x35, 0xdd, 0xcc, 0xc9, 0xe9, 0xf5, 0x43, 0xb7, 0xb7, 0x30, 0x6e, 0xda, 0x37, 0xfa, 0xf8, 0x4f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x80, 0xc6, 0x76, 0x1f, 0x04, 0x03, 0x00, 0x00,
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
		dAtA[i] = 0x22
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
		dAtA[i] = 0x1a
	}
	if m.ExpiryBlock != 0 {
		i = encodeVarintAirdrop(dAtA, i, uint64(m.ExpiryBlock))
		i--
		dAtA[i] = 0x10
	}
	if m.OriginAccountsCreated {
		i--
		if m.OriginAccountsCreated {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
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
		dAtA[i] = 0x30
	}
	if len(m.ClaimAddress) > 0 {
		i -= len(m.ClaimAddress)
		copy(dAtA[i:], m.ClaimAddress)
		i = encodeVarintAirdrop(dAtA, i, uint64(len(m.ClaimAddress)))
		i--
		dAtA[i] = 0x2a
	}
	if len(m.State) > 0 {
		i -= len(m.State)
		copy(dAtA[i:], m.State)
		i = encodeVarintAirdrop(dAtA, i, uint64(len(m.State)))
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
	if m.OriginAccountsCreated {
		n += 2
	}
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
	l = len(m.State)
	if l > 0 {
		n += 1 + l + sovAirdrop(uint64(l))
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
				return fmt.Errorf("proto: wrong wireType = %d for field OriginAccountsCreated", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowAirdrop
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.OriginAccountsCreated = bool(v != 0)
		case 2:
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
		case 3:
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
		case 4:
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
				return fmt.Errorf("proto: wrong wireType = %d for field State", wireType)
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
			m.State = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 5:
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
		case 6:
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
