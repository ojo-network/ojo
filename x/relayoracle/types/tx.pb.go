// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: ojo/relayoracle/v1/tx.proto

package types

import (
	context "context"
	fmt "fmt"
	_ "github.com/cosmos/cosmos-proto"
	_ "github.com/cosmos/cosmos-sdk/types/msgservice"
	grpc1 "github.com/cosmos/gogoproto/grpc"
	proto "github.com/cosmos/gogoproto/proto"
	_ "github.com/gogo/protobuf/gogoproto"
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

// MsgGovUpdateParams defines the Msg/GovUpdateParams request type.
type MsgGovUpdateParams struct {
	// authority is the address of the governance account.
	Authority   string   `protobuf:"bytes,1,opt,name=authority,proto3" json:"authority,omitempty"`
	Title       string   `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Description string   `protobuf:"bytes,3,opt,name=description,proto3" json:"description,omitempty"`
	Keys        []string `protobuf:"bytes,4,rep,name=keys,proto3" json:"keys,omitempty"`
	Changes     Params   `protobuf:"bytes,5,opt,name=changes,proto3" json:"changes"`
}

func (m *MsgGovUpdateParams) Reset()         { *m = MsgGovUpdateParams{} }
func (m *MsgGovUpdateParams) String() string { return proto.CompactTextString(m) }
func (*MsgGovUpdateParams) ProtoMessage()    {}
func (*MsgGovUpdateParams) Descriptor() ([]byte, []int) {
	return fileDescriptor_c6c83bb585b4bc3b, []int{0}
}
func (m *MsgGovUpdateParams) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgGovUpdateParams) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgGovUpdateParams.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgGovUpdateParams) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgGovUpdateParams.Merge(m, src)
}
func (m *MsgGovUpdateParams) XXX_Size() int {
	return m.Size()
}
func (m *MsgGovUpdateParams) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgGovUpdateParams.DiscardUnknown(m)
}

var xxx_messageInfo_MsgGovUpdateParams proto.InternalMessageInfo

// MsgGovUpdateParamsResponse defines the Msg/GovUpdateParams response type.
type MsgGovUpdateParamsResponse struct {
}

func (m *MsgGovUpdateParamsResponse) Reset()         { *m = MsgGovUpdateParamsResponse{} }
func (m *MsgGovUpdateParamsResponse) String() string { return proto.CompactTextString(m) }
func (*MsgGovUpdateParamsResponse) ProtoMessage()    {}
func (*MsgGovUpdateParamsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_c6c83bb585b4bc3b, []int{1}
}
func (m *MsgGovUpdateParamsResponse) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *MsgGovUpdateParamsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_MsgGovUpdateParamsResponse.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *MsgGovUpdateParamsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MsgGovUpdateParamsResponse.Merge(m, src)
}
func (m *MsgGovUpdateParamsResponse) XXX_Size() int {
	return m.Size()
}
func (m *MsgGovUpdateParamsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_MsgGovUpdateParamsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_MsgGovUpdateParamsResponse proto.InternalMessageInfo

func init() {
	proto.RegisterType((*MsgGovUpdateParams)(nil), "ojo.relayoracle.v1.MsgGovUpdateParams")
	proto.RegisterType((*MsgGovUpdateParamsResponse)(nil), "ojo.relayoracle.v1.MsgGovUpdateParamsResponse")
}

func init() { proto.RegisterFile("ojo/relayoracle/v1/tx.proto", fileDescriptor_c6c83bb585b4bc3b) }

var fileDescriptor_c6c83bb585b4bc3b = []byte{
	// 395 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x91, 0xbf, 0x8e, 0xda, 0x40,
	0x10, 0xc6, 0xbd, 0x01, 0x12, 0xb1, 0x14, 0x91, 0x56, 0x48, 0x71, 0x9c, 0xc8, 0xb6, 0x28, 0x22,
	0x14, 0x89, 0xb5, 0x20, 0x52, 0x0a, 0xba, 0xd0, 0x24, 0x0d, 0x52, 0xe4, 0x28, 0x4d, 0x9a, 0xc8,
	0xd8, 0xab, 0xc5, 0x80, 0x3d, 0xd6, 0xee, 0x42, 0x70, 0x9b, 0x2a, 0x65, 0xca, 0x94, 0x3c, 0x42,
	0x8a, 0x3c, 0x04, 0x25, 0x4a, 0x75, 0xd5, 0xe9, 0x04, 0xc5, 0x5d, 0x79, 0x8f, 0x70, 0xf2, 0x1f,
	0x04, 0x77, 0x50, 0x5c, 0x37, 0xb3, 0xbf, 0xf9, 0xe6, 0x9b, 0x9d, 0xc1, 0xaf, 0x60, 0x02, 0x8e,
	0x60, 0x33, 0x2f, 0x05, 0xe1, 0xf9, 0x33, 0xe6, 0x2c, 0xba, 0x8e, 0x5a, 0xd2, 0x44, 0x80, 0x02,
	0x42, 0x60, 0x02, 0xf4, 0x08, 0xd2, 0x45, 0xd7, 0x68, 0x72, 0xe0, 0x90, 0x63, 0x27, 0x8b, 0x8a,
	0x4a, 0xe3, 0xa5, 0x0f, 0x32, 0x02, 0xf9, 0xbd, 0x00, 0x45, 0x52, 0xa2, 0x17, 0x45, 0xe6, 0x44,
	0x92, 0x67, 0xcd, 0x23, 0xc9, 0x4b, 0x60, 0x9d, 0xb1, 0x4e, 0x3c, 0xe1, 0x45, 0xa5, 0xb2, 0x75,
	0x8b, 0x30, 0x19, 0x4a, 0xfe, 0x11, 0x16, 0x5f, 0x93, 0xc0, 0x53, 0xec, 0x73, 0x0e, 0xc9, 0x7b,
	0x5c, 0xf7, 0xe6, 0x6a, 0x0c, 0x22, 0x54, 0xa9, 0x8e, 0x6c, 0xd4, 0xae, 0x0f, 0xf4, 0xff, 0xff,
	0x3a, 0xcd, 0xd2, 0xf5, 0x43, 0x10, 0x08, 0x26, 0xe5, 0x17, 0x25, 0xc2, 0x98, 0xbb, 0x87, 0x52,
	0xd2, 0xc4, 0x35, 0x15, 0xaa, 0x19, 0xd3, 0x9f, 0x64, 0x1a, 0xb7, 0x48, 0x88, 0x8d, 0x1b, 0x01,
	0x93, 0xbe, 0x08, 0x13, 0x15, 0x42, 0xac, 0x57, 0x72, 0x76, 0xfc, 0x44, 0x08, 0xae, 0x4e, 0x59,
	0x2a, 0xf5, 0xaa, 0x5d, 0x69, 0xd7, 0xdd, 0x3c, 0x26, 0x7d, 0xfc, 0xcc, 0x1f, 0x7b, 0x31, 0x67,
	0x52, 0xaf, 0xd9, 0xa8, 0xdd, 0xe8, 0x19, 0xf4, 0x74, 0x57, 0xb4, 0x18, 0x78, 0x50, 0x5d, 0x5f,
	0x5a, 0x9a, 0xbb, 0x17, 0xf4, 0x8d, 0x5f, 0x2b, 0x4b, 0xfb, 0xb3, 0xb2, 0xd0, 0xcd, 0xca, 0x42,
	0x3f, 0xaf, 0xff, 0xbe, 0x3d, 0xcc, 0xd8, 0x7a, 0x8d, 0x8d, 0xd3, 0x1f, 0xbb, 0x4c, 0x26, 0x10,
	0x4b, 0xd6, 0x4b, 0x70, 0x65, 0x28, 0x39, 0x09, 0xf1, 0xf3, 0x87, 0x3b, 0x79, 0x73, 0xce, 0xfe,
	0xb4, 0x93, 0x41, 0x1f, 0x57, 0xb7, 0x77, 0x1c, 0x7c, 0x5a, 0x6f, 0x4d, 0xb4, 0xd9, 0x9a, 0xe8,
	0x6a, 0x6b, 0xa2, 0xdf, 0x3b, 0x53, 0xdb, 0xec, 0x4c, 0xed, 0x62, 0x67, 0x6a, 0xdf, 0x28, 0x0f,
	0xd5, 0x78, 0x3e, 0xa2, 0x3e, 0x44, 0x0e, 0x4c, 0xa0, 0x13, 0x33, 0xf5, 0x03, 0xc4, 0x34, 0x8b,
	0x9d, 0xe5, 0xbd, 0xb3, 0xaa, 0x34, 0x61, 0x72, 0xf4, 0x34, 0xbf, 0xe9, 0xbb, 0xbb, 0x00, 0x00,
	0x00, 0xff, 0xff, 0xf0, 0x0d, 0x80, 0x82, 0x71, 0x02, 0x00, 0x00,
}

func (this *MsgGovUpdateParams) Equal(that interface{}) bool {
	if that == nil {
		return this == nil
	}

	that1, ok := that.(*MsgGovUpdateParams)
	if !ok {
		that2, ok := that.(MsgGovUpdateParams)
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
	if this.Authority != that1.Authority {
		return false
	}
	if this.Title != that1.Title {
		return false
	}
	if this.Description != that1.Description {
		return false
	}
	if len(this.Keys) != len(that1.Keys) {
		return false
	}
	for i := range this.Keys {
		if this.Keys[i] != that1.Keys[i] {
			return false
		}
	}
	if !this.Changes.Equal(&that1.Changes) {
		return false
	}
	return true
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
	GovUpdateParams(ctx context.Context, in *MsgGovUpdateParams, opts ...grpc.CallOption) (*MsgGovUpdateParamsResponse, error)
}

type msgClient struct {
	cc grpc1.ClientConn
}

func NewMsgClient(cc grpc1.ClientConn) MsgClient {
	return &msgClient{cc}
}

func (c *msgClient) GovUpdateParams(ctx context.Context, in *MsgGovUpdateParams, opts ...grpc.CallOption) (*MsgGovUpdateParamsResponse, error) {
	out := new(MsgGovUpdateParamsResponse)
	err := c.cc.Invoke(ctx, "/ojo.relayoracle.v1.Msg/GovUpdateParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MsgServer is the server API for Msg service.
type MsgServer interface {
	GovUpdateParams(context.Context, *MsgGovUpdateParams) (*MsgGovUpdateParamsResponse, error)
}

// UnimplementedMsgServer can be embedded to have forward compatible implementations.
type UnimplementedMsgServer struct {
}

func (*UnimplementedMsgServer) GovUpdateParams(ctx context.Context, req *MsgGovUpdateParams) (*MsgGovUpdateParamsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GovUpdateParams not implemented")
}

func RegisterMsgServer(s grpc1.Server, srv MsgServer) {
	s.RegisterService(&_Msg_serviceDesc, srv)
}

func _Msg_GovUpdateParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MsgGovUpdateParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MsgServer).GovUpdateParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ojo.relayoracle.v1.Msg/GovUpdateParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MsgServer).GovUpdateParams(ctx, req.(*MsgGovUpdateParams))
	}
	return interceptor(ctx, in, info, handler)
}

var _Msg_serviceDesc = grpc.ServiceDesc{
	ServiceName: "ojo.relayoracle.v1.Msg",
	HandlerType: (*MsgServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GovUpdateParams",
			Handler:    _Msg_GovUpdateParams_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ojo/relayoracle/v1/tx.proto",
}

func (m *MsgGovUpdateParams) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgGovUpdateParams) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgGovUpdateParams) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	{
		size, err := m.Changes.MarshalToSizedBuffer(dAtA[:i])
		if err != nil {
			return 0, err
		}
		i -= size
		i = encodeVarintTx(dAtA, i, uint64(size))
	}
	i--
	dAtA[i] = 0x2a
	if len(m.Keys) > 0 {
		for iNdEx := len(m.Keys) - 1; iNdEx >= 0; iNdEx-- {
			i -= len(m.Keys[iNdEx])
			copy(dAtA[i:], m.Keys[iNdEx])
			i = encodeVarintTx(dAtA, i, uint64(len(m.Keys[iNdEx])))
			i--
			dAtA[i] = 0x22
		}
	}
	if len(m.Description) > 0 {
		i -= len(m.Description)
		copy(dAtA[i:], m.Description)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Description)))
		i--
		dAtA[i] = 0x1a
	}
	if len(m.Title) > 0 {
		i -= len(m.Title)
		copy(dAtA[i:], m.Title)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Title)))
		i--
		dAtA[i] = 0x12
	}
	if len(m.Authority) > 0 {
		i -= len(m.Authority)
		copy(dAtA[i:], m.Authority)
		i = encodeVarintTx(dAtA, i, uint64(len(m.Authority)))
		i--
		dAtA[i] = 0xa
	}
	return len(dAtA) - i, nil
}

func (m *MsgGovUpdateParamsResponse) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *MsgGovUpdateParamsResponse) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *MsgGovUpdateParamsResponse) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
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
func (m *MsgGovUpdateParams) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	l = len(m.Authority)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Title)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	l = len(m.Description)
	if l > 0 {
		n += 1 + l + sovTx(uint64(l))
	}
	if len(m.Keys) > 0 {
		for _, s := range m.Keys {
			l = len(s)
			n += 1 + l + sovTx(uint64(l))
		}
	}
	l = m.Changes.Size()
	n += 1 + l + sovTx(uint64(l))
	return n
}

func (m *MsgGovUpdateParamsResponse) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	return n
}

func sovTx(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozTx(x uint64) (n int) {
	return sovTx(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *MsgGovUpdateParams) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgGovUpdateParams: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgGovUpdateParams: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Authority", wireType)
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
			m.Authority = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 2:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Title", wireType)
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
			m.Title = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 3:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Description", wireType)
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
			m.Description = string(dAtA[iNdEx:postIndex])
			iNdEx = postIndex
		case 4:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Keys", wireType)
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
			m.Keys = append(m.Keys, string(dAtA[iNdEx:postIndex]))
			iNdEx = postIndex
		case 5:
			if wireType != 2 {
				return fmt.Errorf("proto: wrong wireType = %d for field Changes", wireType)
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
			if err := m.Changes.Unmarshal(dAtA[iNdEx:postIndex]); err != nil {
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
func (m *MsgGovUpdateParamsResponse) Unmarshal(dAtA []byte) error {
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
			return fmt.Errorf("proto: MsgGovUpdateParamsResponse: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: MsgGovUpdateParamsResponse: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
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
