// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

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

// APIServiceClient is the client API for APIService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type APIServiceClient interface {
	DescribeRAN(ctx context.Context, in *DescribeRANRequest, opts ...grpc.CallOption) (*DescribeRANResponse, error)
	GetUEs(ctx context.Context, in *GetUEsRequest, opts ...grpc.CallOption) (*GetUEsResponse, error)
	DescribeUE(ctx context.Context, in *DescribeUERequest, opts ...grpc.CallOption) (*DescribeUEResponse, error)
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	ServiceRequestProc(ctx context.Context, in *ServiceRequest, opts ...grpc.CallOption) (*ServiceRequestResult, error)
	Deregister(ctx context.Context, in *DeregisterRequest, opts ...grpc.CallOption) (*DeregisterResponse, error)
	SubscribeLog(ctx context.Context, in *LogStreamingRequest, opts ...grpc.CallOption) (APIService_SubscribeLogClient, error)
	ConnectAMF(ctx context.Context, in *ConnectAMFRequest, opts ...grpc.CallOption) (*ConnectAMFResponse, error)
}

type aPIServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAPIServiceClient(cc grpc.ClientConnInterface) APIServiceClient {
	return &aPIServiceClient{cc}
}

func (c *aPIServiceClient) DescribeRAN(ctx context.Context, in *DescribeRANRequest, opts ...grpc.CallOption) (*DescribeRANResponse, error) {
	out := new(DescribeRANResponse)
	err := c.cc.Invoke(ctx, "/APIService/DescribeRAN", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIServiceClient) GetUEs(ctx context.Context, in *GetUEsRequest, opts ...grpc.CallOption) (*GetUEsResponse, error) {
	out := new(GetUEsResponse)
	err := c.cc.Invoke(ctx, "/APIService/GetUEs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIServiceClient) DescribeUE(ctx context.Context, in *DescribeUERequest, opts ...grpc.CallOption) (*DescribeUEResponse, error) {
	out := new(DescribeUEResponse)
	err := c.cc.Invoke(ctx, "/APIService/DescribeUE", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/APIService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIServiceClient) ServiceRequestProc(ctx context.Context, in *ServiceRequest, opts ...grpc.CallOption) (*ServiceRequestResult, error) {
	out := new(ServiceRequestResult)
	err := c.cc.Invoke(ctx, "/APIService/ServiceRequestProc", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIServiceClient) Deregister(ctx context.Context, in *DeregisterRequest, opts ...grpc.CallOption) (*DeregisterResponse, error) {
	out := new(DeregisterResponse)
	err := c.cc.Invoke(ctx, "/APIService/Deregister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *aPIServiceClient) SubscribeLog(ctx context.Context, in *LogStreamingRequest, opts ...grpc.CallOption) (APIService_SubscribeLogClient, error) {
	stream, err := c.cc.NewStream(ctx, &APIService_ServiceDesc.Streams[0], "/APIService/SubscribeLog", opts...)
	if err != nil {
		return nil, err
	}
	x := &aPIServiceSubscribeLogClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type APIService_SubscribeLogClient interface {
	Recv() (*LogStreamingResponse, error)
	grpc.ClientStream
}

type aPIServiceSubscribeLogClient struct {
	grpc.ClientStream
}

func (x *aPIServiceSubscribeLogClient) Recv() (*LogStreamingResponse, error) {
	m := new(LogStreamingResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *aPIServiceClient) ConnectAMF(ctx context.Context, in *ConnectAMFRequest, opts ...grpc.CallOption) (*ConnectAMFResponse, error) {
	out := new(ConnectAMFResponse)
	err := c.cc.Invoke(ctx, "/APIService/ConnectAMF", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// APIServiceServer is the server API for APIService service.
// All implementations must embed UnimplementedAPIServiceServer
// for forward compatibility
type APIServiceServer interface {
	DescribeRAN(context.Context, *DescribeRANRequest) (*DescribeRANResponse, error)
	GetUEs(context.Context, *GetUEsRequest) (*GetUEsResponse, error)
	DescribeUE(context.Context, *DescribeUERequest) (*DescribeUEResponse, error)
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	ServiceRequestProc(context.Context, *ServiceRequest) (*ServiceRequestResult, error)
	Deregister(context.Context, *DeregisterRequest) (*DeregisterResponse, error)
	SubscribeLog(*LogStreamingRequest, APIService_SubscribeLogServer) error
	ConnectAMF(context.Context, *ConnectAMFRequest) (*ConnectAMFResponse, error)
	mustEmbedUnimplementedAPIServiceServer()
}

// UnimplementedAPIServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAPIServiceServer struct {
}

func (UnimplementedAPIServiceServer) DescribeRAN(context.Context, *DescribeRANRequest) (*DescribeRANResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeRAN not implemented")
}
func (UnimplementedAPIServiceServer) GetUEs(context.Context, *GetUEsRequest) (*GetUEsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUEs not implemented")
}
func (UnimplementedAPIServiceServer) DescribeUE(context.Context, *DescribeUERequest) (*DescribeUEResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeUE not implemented")
}
func (UnimplementedAPIServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedAPIServiceServer) ServiceRequestProc(context.Context, *ServiceRequest) (*ServiceRequestResult, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ServiceRequestProc not implemented")
}
func (UnimplementedAPIServiceServer) Deregister(context.Context, *DeregisterRequest) (*DeregisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Deregister not implemented")
}
func (UnimplementedAPIServiceServer) SubscribeLog(*LogStreamingRequest, APIService_SubscribeLogServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeLog not implemented")
}
func (UnimplementedAPIServiceServer) ConnectAMF(context.Context, *ConnectAMFRequest) (*ConnectAMFResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConnectAMF not implemented")
}
func (UnimplementedAPIServiceServer) mustEmbedUnimplementedAPIServiceServer() {}

// UnsafeAPIServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to APIServiceServer will
// result in compilation errors.
type UnsafeAPIServiceServer interface {
	mustEmbedUnimplementedAPIServiceServer()
}

func RegisterAPIServiceServer(s grpc.ServiceRegistrar, srv APIServiceServer) {
	s.RegisterService(&APIService_ServiceDesc, srv)
}

func _APIService_DescribeRAN_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribeRANRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServiceServer).DescribeRAN(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/APIService/DescribeRAN",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServiceServer).DescribeRAN(ctx, req.(*DescribeRANRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _APIService_GetUEs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUEsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServiceServer).GetUEs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/APIService/GetUEs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServiceServer).GetUEs(ctx, req.(*GetUEsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _APIService_DescribeUE_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribeUERequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServiceServer).DescribeUE(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/APIService/DescribeUE",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServiceServer).DescribeUE(ctx, req.(*DescribeUERequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _APIService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/APIService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _APIService_ServiceRequestProc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServiceRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServiceServer).ServiceRequestProc(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/APIService/ServiceRequestProc",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServiceServer).ServiceRequestProc(ctx, req.(*ServiceRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _APIService_Deregister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeregisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServiceServer).Deregister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/APIService/Deregister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServiceServer).Deregister(ctx, req.(*DeregisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _APIService_SubscribeLog_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(LogStreamingRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(APIServiceServer).SubscribeLog(m, &aPIServiceSubscribeLogServer{stream})
}

type APIService_SubscribeLogServer interface {
	Send(*LogStreamingResponse) error
	grpc.ServerStream
}

type aPIServiceSubscribeLogServer struct {
	grpc.ServerStream
}

func (x *aPIServiceSubscribeLogServer) Send(m *LogStreamingResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _APIService_ConnectAMF_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConnectAMFRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(APIServiceServer).ConnectAMF(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/APIService/ConnectAMF",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(APIServiceServer).ConnectAMF(ctx, req.(*ConnectAMFRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// APIService_ServiceDesc is the grpc.ServiceDesc for APIService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var APIService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "APIService",
	HandlerType: (*APIServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DescribeRAN",
			Handler:    _APIService_DescribeRAN_Handler,
		},
		{
			MethodName: "GetUEs",
			Handler:    _APIService_GetUEs_Handler,
		},
		{
			MethodName: "DescribeUE",
			Handler:    _APIService_DescribeUE_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _APIService_Register_Handler,
		},
		{
			MethodName: "ServiceRequestProc",
			Handler:    _APIService_ServiceRequestProc_Handler,
		},
		{
			MethodName: "Deregister",
			Handler:    _APIService_Deregister_Handler,
		},
		{
			MethodName: "ConnectAMF",
			Handler:    _APIService_ConnectAMF_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeLog",
			Handler:       _APIService_SubscribeLog_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "pkg/api/api.proto",
}
