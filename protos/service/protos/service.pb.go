// Code generated by protoc-gen-go. DO NOT EDIT.
// source: protos/service.proto

/*
Package service is a generated protocol buffer package.

It is generated from these files:
	protos/service.proto

It has these top-level messages:
	StatusRequest
	StatusResponse
*/
package service

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type StatusRequest struct {
}

func (m *StatusRequest) Reset()                    { *m = StatusRequest{} }
func (m *StatusRequest) String() string            { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()               {}
func (*StatusRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type StatusResponse struct {
	Status bool   `protobuf:"varint,1,opt,name=Status,json=status" json:"Status,omitempty"`
	Time   string `protobuf:"bytes,2,opt,name=Time,json=time" json:"Time,omitempty"`
}

func (m *StatusResponse) Reset()                    { *m = StatusResponse{} }
func (m *StatusResponse) String() string            { return proto.CompactTextString(m) }
func (*StatusResponse) ProtoMessage()               {}
func (*StatusResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *StatusResponse) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

func (m *StatusResponse) GetTime() string {
	if m != nil {
		return m.Time
	}
	return ""
}

func init() {
	proto.RegisterType((*StatusRequest)(nil), "StatusRequest")
	proto.RegisterType((*StatusResponse)(nil), "StatusResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Service service

type ServiceClient interface {
	HeartBeatStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error)
}

type serviceClient struct {
	cc *grpc.ClientConn
}

func NewServiceClient(cc *grpc.ClientConn) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) HeartBeatStatus(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := grpc.Invoke(ctx, "/Service/HeartBeatStatus", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Service service

type ServiceServer interface {
	HeartBeatStatus(context.Context, *StatusRequest) (*StatusResponse, error)
}

func RegisterServiceServer(s *grpc.Server, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_HeartBeatStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).HeartBeatStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Service/HeartBeatStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).HeartBeatStatus(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HeartBeatStatus",
			Handler:    _Service_HeartBeatStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/service.proto",
}

func init() { proto.RegisterFile("protos/service.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 146 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x29, 0x28, 0xca, 0x2f,
	0xc9, 0x2f, 0xd6, 0x2f, 0x4e, 0x2d, 0x2a, 0xcb, 0x4c, 0x4e, 0xd5, 0x03, 0x73, 0x95, 0xf8, 0xb9,
	0x78, 0x83, 0x4b, 0x12, 0x4b, 0x4a, 0x8b, 0x83, 0x52, 0x0b, 0x4b, 0x53, 0x8b, 0x4b, 0x94, 0x6c,
	0xb8, 0xf8, 0x60, 0x02, 0xc5, 0x05, 0xf9, 0x79, 0xc5, 0xa9, 0x42, 0x62, 0x5c, 0x6c, 0x10, 0x11,
	0x09, 0x46, 0x05, 0x46, 0x0d, 0x8e, 0x20, 0xb6, 0x62, 0x30, 0x4f, 0x48, 0x88, 0x8b, 0x25, 0x24,
	0x33, 0x37, 0x55, 0x82, 0x49, 0x81, 0x51, 0x83, 0x33, 0x88, 0xa5, 0x24, 0x33, 0x37, 0xd5, 0xc8,
	0x96, 0x8b, 0x3d, 0x18, 0x62, 0xbe, 0x90, 0x11, 0x17, 0xbf, 0x47, 0x6a, 0x62, 0x51, 0x89, 0x53,
	0x6a, 0x62, 0x09, 0x44, 0xbf, 0x10, 0x9f, 0x1e, 0x8a, 0x5d, 0x52, 0xfc, 0x7a, 0xa8, 0x56, 0x25,
	0xb1, 0x81, 0x1d, 0x65, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x20, 0x4c, 0xb7, 0xd1, 0xac, 0x00,
	0x00, 0x00,
}
