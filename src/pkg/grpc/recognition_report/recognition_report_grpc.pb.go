// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.6
// source: recognition_report.proto

package recognition_report

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

// RecognitionReportClient is the client API for RecognitionReport service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RecognitionReportClient interface {
	Create(ctx context.Context, in *RecognitionReportRequest, opts ...grpc.CallOption) (*RecognitionReportResponse, error)
}

type recognitionReportClient struct {
	cc grpc.ClientConnInterface
}

func NewRecognitionReportClient(cc grpc.ClientConnInterface) RecognitionReportClient {
	return &recognitionReportClient{cc}
}

func (c *recognitionReportClient) Create(ctx context.Context, in *RecognitionReportRequest, opts ...grpc.CallOption) (*RecognitionReportResponse, error) {
	out := new(RecognitionReportResponse)
	err := c.cc.Invoke(ctx, "/RecognitionReport/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RecognitionReportServer is the server API for RecognitionReport service.
// All implementations should embed UnimplementedRecognitionReportServer
// for forward compatibility
type RecognitionReportServer interface {
	Create(context.Context, *RecognitionReportRequest) (*RecognitionReportResponse, error)
}

// UnimplementedRecognitionReportServer should be embedded to have forward compatible implementations.
type UnimplementedRecognitionReportServer struct {
}

func (UnimplementedRecognitionReportServer) Create(context.Context, *RecognitionReportRequest) (*RecognitionReportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}

// UnsafeRecognitionReportServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RecognitionReportServer will
// result in compilation errors.
type UnsafeRecognitionReportServer interface {
	mustEmbedUnimplementedRecognitionReportServer()
}

func RegisterRecognitionReportServer(s grpc.ServiceRegistrar, srv RecognitionReportServer) {
	s.RegisterService(&RecognitionReport_ServiceDesc, srv)
}

func _RecognitionReport_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecognitionReportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RecognitionReportServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/RecognitionReport/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RecognitionReportServer).Create(ctx, req.(*RecognitionReportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RecognitionReport_ServiceDesc is the grpc.ServiceDesc for RecognitionReport service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RecognitionReport_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "RecognitionReport",
	HandlerType: (*RecognitionReportServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _RecognitionReport_Create_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "recognition_report.proto",
}
