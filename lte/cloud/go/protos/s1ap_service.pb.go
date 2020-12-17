// Code generated by protoc-gen-go. DO NOT EDIT.
// source: lte/protos/s1ap_service.proto

package protos

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protos "magma/orc8r/lib/go/protos"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type EnbStateResult struct {
	// enb ids -> num of UEs connected
	EnbStateMap          map[uint32]uint32 `protobuf:"bytes,1,rep,name=enb_state_map,json=enbStateMap,proto3" json:"enb_state_map,omitempty" protobuf_key:"varint,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *EnbStateResult) Reset()         { *m = EnbStateResult{} }
func (m *EnbStateResult) String() string { return proto.CompactTextString(m) }
func (*EnbStateResult) ProtoMessage()    {}
func (*EnbStateResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_b760ca1045684df1, []int{0}
}

func (m *EnbStateResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnbStateResult.Unmarshal(m, b)
}
func (m *EnbStateResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnbStateResult.Marshal(b, m, deterministic)
}
func (m *EnbStateResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnbStateResult.Merge(m, src)
}
func (m *EnbStateResult) XXX_Size() int {
	return xxx_messageInfo_EnbStateResult.Size(m)
}
func (m *EnbStateResult) XXX_DiscardUnknown() {
	xxx_messageInfo_EnbStateResult.DiscardUnknown(m)
}

var xxx_messageInfo_EnbStateResult proto.InternalMessageInfo

func (m *EnbStateResult) GetEnbStateMap() map[uint32]uint32 {
	if m != nil {
		return m.EnbStateMap
	}
	return nil
}

func init() {
	proto.RegisterType((*EnbStateResult)(nil), "magma.lte.EnbStateResult")
	proto.RegisterMapType((map[uint32]uint32)(nil), "magma.lte.EnbStateResult.EnbStateMapEntry")
}

func init() { proto.RegisterFile("lte/protos/s1ap_service.proto", fileDescriptor_b760ca1045684df1) }

var fileDescriptor_b760ca1045684df1 = []byte{
	// 246 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x92, 0xcd, 0x29, 0x49, 0xd5,
	0x2f, 0x28, 0xca, 0x2f, 0xc9, 0x2f, 0xd6, 0x2f, 0x36, 0x4c, 0x2c, 0x88, 0x2f, 0x4e, 0x2d, 0x2a,
	0xcb, 0x4c, 0x4e, 0xd5, 0x03, 0x8b, 0x09, 0x71, 0xe6, 0x26, 0xa6, 0xe7, 0x26, 0xea, 0xe5, 0x94,
	0xa4, 0x4a, 0x49, 0xe6, 0x17, 0x25, 0x5b, 0x14, 0xc1, 0xd4, 0x26, 0xe7, 0xe7, 0xe6, 0xe6, 0xe7,
	0x41, 0x54, 0x29, 0x2d, 0x60, 0xe4, 0xe2, 0x73, 0xcd, 0x4b, 0x0a, 0x2e, 0x49, 0x2c, 0x49, 0x0d,
	0x4a, 0x2d, 0x2e, 0xcd, 0x29, 0x11, 0xf2, 0xe3, 0xe2, 0x4d, 0xcd, 0x4b, 0x8a, 0x2f, 0x06, 0x09,
	0xc5, 0xe7, 0x26, 0x16, 0x48, 0x30, 0x2a, 0x30, 0x6b, 0x70, 0x1b, 0x69, 0xe9, 0xc1, 0x0d, 0xd4,
	0x43, 0xd5, 0x01, 0xe7, 0xfa, 0x26, 0x16, 0xb8, 0xe6, 0x95, 0x14, 0x55, 0x06, 0x71, 0xa7, 0x22,
	0x44, 0xa4, 0xec, 0xb8, 0x04, 0xd0, 0x15, 0x08, 0x09, 0x70, 0x31, 0x67, 0xa7, 0x56, 0x4a, 0x30,
	0x2a, 0x30, 0x6a, 0xf0, 0x06, 0x81, 0x98, 0x42, 0x22, 0x5c, 0xac, 0x65, 0x89, 0x39, 0xa5, 0xa9,
	0x12, 0x4c, 0x60, 0x31, 0x08, 0xc7, 0x8a, 0xc9, 0x82, 0xd1, 0xc8, 0x87, 0x8b, 0x3b, 0xd8, 0x30,
	0xb1, 0x20, 0x18, 0xe2, 0x3b, 0x21, 0x5b, 0x2e, 0x6e, 0xf7, 0xd4, 0x12, 0x98, 0x89, 0x42, 0x82,
	0x50, 0x67, 0x81, 0xbd, 0xa8, 0x17, 0x96, 0x9f, 0x99, 0x22, 0x25, 0x89, 0xd3, 0xa5, 0x4a, 0x0c,
	0x4e, 0xd2, 0x51, 0x92, 0x60, 0x59, 0x7d, 0x50, 0xe8, 0x25, 0xe7, 0xe4, 0x97, 0xa6, 0xe8, 0xa7,
	0xe7, 0x43, 0x83, 0x26, 0x89, 0x0d, 0x4c, 0x1b, 0x03, 0x02, 0x00, 0x00, 0xff, 0xff, 0xce, 0xf1,
	0xd5, 0x60, 0x5b, 0x01, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// S1ApServiceClient is the client API for S1ApService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type S1ApServiceClient interface {
	GetEnbState(ctx context.Context, in *protos.Void, opts ...grpc.CallOption) (*EnbStateResult, error)
}

type s1ApServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewS1ApServiceClient(cc grpc.ClientConnInterface) S1ApServiceClient {
	return &s1ApServiceClient{cc}
}

func (c *s1ApServiceClient) GetEnbState(ctx context.Context, in *protos.Void, opts ...grpc.CallOption) (*EnbStateResult, error) {
	out := new(EnbStateResult)
	err := c.cc.Invoke(ctx, "/magma.lte.S1apService/GetEnbState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// S1ApServiceServer is the server API for S1ApService service.
type S1ApServiceServer interface {
	GetEnbState(context.Context, *protos.Void) (*EnbStateResult, error)
}

// UnimplementedS1ApServiceServer can be embedded to have forward compatible implementations.
type UnimplementedS1ApServiceServer struct {
}

func (*UnimplementedS1ApServiceServer) GetEnbState(ctx context.Context, req *protos.Void) (*EnbStateResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetEnbState not implemented")
}

func RegisterS1ApServiceServer(s *grpc.Server, srv S1ApServiceServer) {
	s.RegisterService(&_S1ApService_serviceDesc, srv)
}

func _S1ApService_GetEnbState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(protos.Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(S1ApServiceServer).GetEnbState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/magma.lte.S1apService/GetEnbState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(S1ApServiceServer).GetEnbState(ctx, req.(*protos.Void))
	}
	return interceptor(ctx, in, info, handler)
}

var _S1ApService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "magma.lte.S1apService",
	HandlerType: (*S1ApServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetEnbState",
			Handler:    _S1ApService_GetEnbState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lte/protos/s1ap_service.proto",
}
