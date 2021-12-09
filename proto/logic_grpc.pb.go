// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// RpcLogicServiceClient is the client API for RpcLogicService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RpcLogicServiceClient interface {
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	TokenLogin(ctx context.Context, in *TokenLoginRequest, opts ...grpc.CallOption) (*TokenLoginResponse, error)
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	CheckAuth(ctx context.Context, in *CheckAuthRequest, opts ...grpc.CallOption) (*CheckAuthResponse, error)
}

type rpcLogicServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewRpcLogicServiceClient(cc grpc.ClientConnInterface) RpcLogicServiceClient {
	return &rpcLogicServiceClient{cc}
}

func (c *rpcLogicServiceClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, "/proto.RpcLogicService/Login", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rpcLogicServiceClient) TokenLogin(ctx context.Context, in *TokenLoginRequest, opts ...grpc.CallOption) (*TokenLoginResponse, error) {
	out := new(TokenLoginResponse)
	err := c.cc.Invoke(ctx, "/proto.RpcLogicService/TokenLogin", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rpcLogicServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, "/proto.RpcLogicService/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rpcLogicServiceClient) CheckAuth(ctx context.Context, in *CheckAuthRequest, opts ...grpc.CallOption) (*CheckAuthResponse, error) {
	out := new(CheckAuthResponse)
	err := c.cc.Invoke(ctx, "/proto.RpcLogicService/CheckAuth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RpcLogicServiceServer is the server API for RpcLogicService service.
// All implementations must embed UnimplementedRpcLogicServiceServer
// for forward compatibility
type RpcLogicServiceServer interface {
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	TokenLogin(context.Context, *TokenLoginRequest) (*TokenLoginResponse, error)
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	CheckAuth(context.Context, *CheckAuthRequest) (*CheckAuthResponse, error)
	mustEmbedUnimplementedRpcLogicServiceServer()
}

// UnimplementedRpcLogicServiceServer must be embedded to have forward compatible implementations.
type UnimplementedRpcLogicServiceServer struct {
}

func (UnimplementedRpcLogicServiceServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedRpcLogicServiceServer) TokenLogin(context.Context, *TokenLoginRequest) (*TokenLoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method TokenLogin not implemented")
}
func (UnimplementedRpcLogicServiceServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedRpcLogicServiceServer) CheckAuth(context.Context, *CheckAuthRequest) (*CheckAuthResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckAuth not implemented")
}
func (UnimplementedRpcLogicServiceServer) mustEmbedUnimplementedRpcLogicServiceServer() {}

// UnsafeRpcLogicServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RpcLogicServiceServer will
// result in compilation errors.
type UnsafeRpcLogicServiceServer interface {
	mustEmbedUnimplementedRpcLogicServiceServer()
}

func RegisterRpcLogicServiceServer(s grpc.ServiceRegistrar, srv RpcLogicServiceServer) {
	s.RegisterService(&RpcLogicService_ServiceDesc, srv)
}

func _RpcLogicService_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcLogicServiceServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.RpcLogicService/Login",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcLogicServiceServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RpcLogicService_TokenLogin_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TokenLoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcLogicServiceServer).TokenLogin(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.RpcLogicService/TokenLogin",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcLogicServiceServer).TokenLogin(ctx, req.(*TokenLoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RpcLogicService_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcLogicServiceServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.RpcLogicService/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcLogicServiceServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _RpcLogicService_CheckAuth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckAuthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RpcLogicServiceServer).CheckAuth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.RpcLogicService/CheckAuth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RpcLogicServiceServer).CheckAuth(ctx, req.(*CheckAuthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RpcLogicService_ServiceDesc is the grpc.ServiceDesc for RpcLogicService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RpcLogicService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.RpcLogicService",
	HandlerType: (*RpcLogicServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Login",
			Handler:    _RpcLogicService_Login_Handler,
		},
		{
			MethodName: "TokenLogin",
			Handler:    _RpcLogicService_TokenLogin_Handler,
		},
		{
			MethodName: "Register",
			Handler:    _RpcLogicService_Register_Handler,
		},
		{
			MethodName: "CheckAuth",
			Handler:    _RpcLogicService_CheckAuth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/logic.proto",
}
