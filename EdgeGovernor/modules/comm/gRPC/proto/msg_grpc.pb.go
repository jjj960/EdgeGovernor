// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v4.22.3
// source: msg.proto

package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NodeCommClient is the client API for NodeComm service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NodeCommClient interface {
	Messaging(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error)
}

type nodeCommClient struct {
	cc grpc.ClientConnInterface
}

func NewNodeCommClient(cc grpc.ClientConnInterface) NodeCommClient {
	return &nodeCommClient{cc}
}

func (c *nodeCommClient) Messaging(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.NodeComm/Messaging", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NodeCommServer is the server API for NodeComm service.
// All implementations must embed UnimplementedNodeCommServer
// for forward compatibility
type NodeCommServer interface {
	Messaging(context.Context, *Request) (*Response, error)
	mustEmbedUnimplementedNodeCommServer()
}

// UnimplementedNodeCommServer must be embedded to have forward compatible implementations.
type UnimplementedNodeCommServer struct {
}

func (UnimplementedNodeCommServer) Messaging(context.Context, *Request) (*Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Messaging not implemented")
}
func (UnimplementedNodeCommServer) mustEmbedUnimplementedNodeCommServer() {}

// UnsafeNodeCommServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NodeCommServer will
// result in compilation errors.
type UnsafeNodeCommServer interface {
	mustEmbedUnimplementedNodeCommServer()
}

func RegisterNodeCommServer(s grpc.ServiceRegistrar, srv NodeCommServer) {
	s.RegisterService(&NodeComm_ServiceDesc, srv)
}

func _NodeComm_Messaging_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NodeCommServer).Messaging(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.NodeComm/Messaging",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NodeCommServer).Messaging(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// NodeComm_ServiceDesc is the grpc.ServiceDesc for NodeComm service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NodeComm_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.NodeComm",
	HandlerType: (*NodeCommServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Messaging",
			Handler:    _NodeComm_Messaging_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "msg.proto",
}