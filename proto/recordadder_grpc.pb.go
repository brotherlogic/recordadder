// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: recordadder.proto

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

// AddRecordServiceClient is the client API for AddRecordService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AddRecordServiceClient interface {
	AddRecord(ctx context.Context, in *AddRecordRequest, opts ...grpc.CallOption) (*AddRecordResponse, error)
	ListQueue(ctx context.Context, in *ListQueueRequest, opts ...grpc.CallOption) (*ListQueueResponse, error)
	UpdateRecord(ctx context.Context, in *UpdateRecordRequest, opts ...grpc.CallOption) (*UpdateRecordResponse, error)
	DeleteRecord(ctx context.Context, in *DeleteRecordRequest, opts ...grpc.CallOption) (*DeleteRecordResponse, error)
	ProcAdded(ctx context.Context, in *ProcAddedRequest, opts ...grpc.CallOption) (*ProcAddedResponse, error)
}

type addRecordServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAddRecordServiceClient(cc grpc.ClientConnInterface) AddRecordServiceClient {
	return &addRecordServiceClient{cc}
}

func (c *addRecordServiceClient) AddRecord(ctx context.Context, in *AddRecordRequest, opts ...grpc.CallOption) (*AddRecordResponse, error) {
	out := new(AddRecordResponse)
	err := c.cc.Invoke(ctx, "/recordadder.AddRecordService/AddRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addRecordServiceClient) ListQueue(ctx context.Context, in *ListQueueRequest, opts ...grpc.CallOption) (*ListQueueResponse, error) {
	out := new(ListQueueResponse)
	err := c.cc.Invoke(ctx, "/recordadder.AddRecordService/ListQueue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addRecordServiceClient) UpdateRecord(ctx context.Context, in *UpdateRecordRequest, opts ...grpc.CallOption) (*UpdateRecordResponse, error) {
	out := new(UpdateRecordResponse)
	err := c.cc.Invoke(ctx, "/recordadder.AddRecordService/UpdateRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addRecordServiceClient) DeleteRecord(ctx context.Context, in *DeleteRecordRequest, opts ...grpc.CallOption) (*DeleteRecordResponse, error) {
	out := new(DeleteRecordResponse)
	err := c.cc.Invoke(ctx, "/recordadder.AddRecordService/DeleteRecord", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addRecordServiceClient) ProcAdded(ctx context.Context, in *ProcAddedRequest, opts ...grpc.CallOption) (*ProcAddedResponse, error) {
	out := new(ProcAddedResponse)
	err := c.cc.Invoke(ctx, "/recordadder.AddRecordService/ProcAdded", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AddRecordServiceServer is the server API for AddRecordService service.
// All implementations should embed UnimplementedAddRecordServiceServer
// for forward compatibility
type AddRecordServiceServer interface {
	AddRecord(context.Context, *AddRecordRequest) (*AddRecordResponse, error)
	ListQueue(context.Context, *ListQueueRequest) (*ListQueueResponse, error)
	UpdateRecord(context.Context, *UpdateRecordRequest) (*UpdateRecordResponse, error)
	DeleteRecord(context.Context, *DeleteRecordRequest) (*DeleteRecordResponse, error)
	ProcAdded(context.Context, *ProcAddedRequest) (*ProcAddedResponse, error)
}

// UnimplementedAddRecordServiceServer should be embedded to have forward compatible implementations.
type UnimplementedAddRecordServiceServer struct {
}

func (UnimplementedAddRecordServiceServer) AddRecord(context.Context, *AddRecordRequest) (*AddRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddRecord not implemented")
}
func (UnimplementedAddRecordServiceServer) ListQueue(context.Context, *ListQueueRequest) (*ListQueueResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListQueue not implemented")
}
func (UnimplementedAddRecordServiceServer) UpdateRecord(context.Context, *UpdateRecordRequest) (*UpdateRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateRecord not implemented")
}
func (UnimplementedAddRecordServiceServer) DeleteRecord(context.Context, *DeleteRecordRequest) (*DeleteRecordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteRecord not implemented")
}
func (UnimplementedAddRecordServiceServer) ProcAdded(context.Context, *ProcAddedRequest) (*ProcAddedResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ProcAdded not implemented")
}

// UnsafeAddRecordServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AddRecordServiceServer will
// result in compilation errors.
type UnsafeAddRecordServiceServer interface {
	mustEmbedUnimplementedAddRecordServiceServer()
}

func RegisterAddRecordServiceServer(s grpc.ServiceRegistrar, srv AddRecordServiceServer) {
	s.RegisterService(&AddRecordService_ServiceDesc, srv)
}

func _AddRecordService_AddRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddRecordServiceServer).AddRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordadder.AddRecordService/AddRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddRecordServiceServer).AddRecord(ctx, req.(*AddRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AddRecordService_ListQueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddRecordServiceServer).ListQueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordadder.AddRecordService/ListQueue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddRecordServiceServer).ListQueue(ctx, req.(*ListQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AddRecordService_UpdateRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddRecordServiceServer).UpdateRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordadder.AddRecordService/UpdateRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddRecordServiceServer).UpdateRecord(ctx, req.(*UpdateRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AddRecordService_DeleteRecord_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRecordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddRecordServiceServer).DeleteRecord(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordadder.AddRecordService/DeleteRecord",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddRecordServiceServer).DeleteRecord(ctx, req.(*DeleteRecordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AddRecordService_ProcAdded_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcAddedRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddRecordServiceServer).ProcAdded(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordadder.AddRecordService/ProcAdded",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddRecordServiceServer).ProcAdded(ctx, req.(*ProcAddedRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AddRecordService_ServiceDesc is the grpc.ServiceDesc for AddRecordService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AddRecordService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "recordadder.AddRecordService",
	HandlerType: (*AddRecordServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddRecord",
			Handler:    _AddRecordService_AddRecord_Handler,
		},
		{
			MethodName: "ListQueue",
			Handler:    _AddRecordService_ListQueue_Handler,
		},
		{
			MethodName: "UpdateRecord",
			Handler:    _AddRecordService_UpdateRecord_Handler,
		},
		{
			MethodName: "DeleteRecord",
			Handler:    _AddRecordService_DeleteRecord_Handler,
		},
		{
			MethodName: "ProcAdded",
			Handler:    _AddRecordService_ProcAdded_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "recordadder.proto",
}

// ClientAddUpdateServiceClient is the client API for ClientAddUpdateService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ClientAddUpdateServiceClient interface {
	ClientAddUpdate(ctx context.Context, in *ClientAddUpdateRequest, opts ...grpc.CallOption) (*ClientAddUpdateResponse, error)
}

type clientAddUpdateServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewClientAddUpdateServiceClient(cc grpc.ClientConnInterface) ClientAddUpdateServiceClient {
	return &clientAddUpdateServiceClient{cc}
}

func (c *clientAddUpdateServiceClient) ClientAddUpdate(ctx context.Context, in *ClientAddUpdateRequest, opts ...grpc.CallOption) (*ClientAddUpdateResponse, error) {
	out := new(ClientAddUpdateResponse)
	err := c.cc.Invoke(ctx, "/recordadder.ClientAddUpdateService/ClientAddUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ClientAddUpdateServiceServer is the server API for ClientAddUpdateService service.
// All implementations should embed UnimplementedClientAddUpdateServiceServer
// for forward compatibility
type ClientAddUpdateServiceServer interface {
	ClientAddUpdate(context.Context, *ClientAddUpdateRequest) (*ClientAddUpdateResponse, error)
}

// UnimplementedClientAddUpdateServiceServer should be embedded to have forward compatible implementations.
type UnimplementedClientAddUpdateServiceServer struct {
}

func (UnimplementedClientAddUpdateServiceServer) ClientAddUpdate(context.Context, *ClientAddUpdateRequest) (*ClientAddUpdateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ClientAddUpdate not implemented")
}

// UnsafeClientAddUpdateServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ClientAddUpdateServiceServer will
// result in compilation errors.
type UnsafeClientAddUpdateServiceServer interface {
	mustEmbedUnimplementedClientAddUpdateServiceServer()
}

func RegisterClientAddUpdateServiceServer(s grpc.ServiceRegistrar, srv ClientAddUpdateServiceServer) {
	s.RegisterService(&ClientAddUpdateService_ServiceDesc, srv)
}

func _ClientAddUpdateService_ClientAddUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientAddUpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ClientAddUpdateServiceServer).ClientAddUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/recordadder.ClientAddUpdateService/ClientAddUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ClientAddUpdateServiceServer).ClientAddUpdate(ctx, req.(*ClientAddUpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// ClientAddUpdateService_ServiceDesc is the grpc.ServiceDesc for ClientAddUpdateService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ClientAddUpdateService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "recordadder.ClientAddUpdateService",
	HandlerType: (*ClientAddUpdateServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ClientAddUpdate",
			Handler:    _ClientAddUpdateService_ClientAddUpdate_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "recordadder.proto",
}
