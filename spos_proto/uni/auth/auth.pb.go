// Code generated by protoc-gen-go. DO NOT EDIT.
// source: uni/auth.proto

package uni_auth

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

// The request message containing the text lang gender audio_code
type LoginReq struct {
	Username             string   `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Passwd               string   `protobuf:"bytes,2,opt,name=passwd,proto3" json:"passwd,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginReq) Reset()         { *m = LoginReq{} }
func (m *LoginReq) String() string { return proto.CompactTextString(m) }
func (*LoginReq) ProtoMessage()    {}
func (*LoginReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_auth_ec0d7c7920b062ae, []int{0}
}
func (m *LoginReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginReq.Unmarshal(m, b)
}
func (m *LoginReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginReq.Marshal(b, m, deterministic)
}
func (dst *LoginReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginReq.Merge(dst, src)
}
func (m *LoginReq) XXX_Size() int {
	return xxx_messageInfo_LoginReq.Size(m)
}
func (m *LoginReq) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginReq.DiscardUnknown(m)
}

var xxx_messageInfo_LoginReq proto.InternalMessageInfo

func (m *LoginReq) GetUsername() string {
	if m != nil {
		return m.Username
	}
	return ""
}

func (m *LoginReq) GetPasswd() string {
	if m != nil {
		return m.Passwd
	}
	return ""
}

// The response message containing the audio data
type LoginReply struct {
	Ret                  bool     `protobuf:"varint,1,opt,name=ret,proto3" json:"ret,omitempty"`
	Token                string   `protobuf:"bytes,2,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *LoginReply) Reset()         { *m = LoginReply{} }
func (m *LoginReply) String() string { return proto.CompactTextString(m) }
func (*LoginReply) ProtoMessage()    {}
func (*LoginReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_auth_ec0d7c7920b062ae, []int{1}
}
func (m *LoginReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_LoginReply.Unmarshal(m, b)
}
func (m *LoginReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_LoginReply.Marshal(b, m, deterministic)
}
func (dst *LoginReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_LoginReply.Merge(dst, src)
}
func (m *LoginReply) XXX_Size() int {
	return xxx_messageInfo_LoginReply.Size(m)
}
func (m *LoginReply) XXX_DiscardUnknown() {
	xxx_messageInfo_LoginReply.DiscardUnknown(m)
}

var xxx_messageInfo_LoginReply proto.InternalMessageInfo

func (m *LoginReply) GetRet() bool {
	if m != nil {
		return m.Ret
	}
	return false
}

func (m *LoginReply) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func init() {
	proto.RegisterType((*LoginReq)(nil), "uni_auth.LoginReq")
	proto.RegisterType((*LoginReply)(nil), "uni_auth.LoginReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AuthClient interface {
	// tts req
	LoginRequest(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (Auth_LoginRequestClient, error)
}

type authClient struct {
	cc *grpc.ClientConn
}

func NewAuthClient(cc *grpc.ClientConn) AuthClient {
	return &authClient{cc}
}

func (c *authClient) LoginRequest(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (Auth_LoginRequestClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Auth_serviceDesc.Streams[0], "/uni_auth.Auth/LoginRequest", opts...)
	if err != nil {
		return nil, err
	}
	x := &authLoginRequestClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Auth_LoginRequestClient interface {
	Recv() (*LoginReply, error)
	grpc.ClientStream
}

type authLoginRequestClient struct {
	grpc.ClientStream
}

func (x *authLoginRequestClient) Recv() (*LoginReply, error) {
	m := new(LoginReply)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// AuthServer is the server API for Auth service.
type AuthServer interface {
	// tts req
	LoginRequest(*LoginReq, Auth_LoginRequestServer) error
}

func RegisterAuthServer(s *grpc.Server, srv AuthServer) {
	s.RegisterService(&_Auth_serviceDesc, srv)
}

func _Auth_LoginRequest_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(LoginReq)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(AuthServer).LoginRequest(m, &authLoginRequestServer{stream})
}

type Auth_LoginRequestServer interface {
	Send(*LoginReply) error
	grpc.ServerStream
}

type authLoginRequestServer struct {
	grpc.ServerStream
}

func (x *authLoginRequestServer) Send(m *LoginReply) error {
	return x.ServerStream.SendMsg(m)
}

var _Auth_serviceDesc = grpc.ServiceDesc{
	ServiceName: "uni_auth.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "LoginRequest",
			Handler:       _Auth_LoginRequest_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "uni/auth.proto",
}

func init() { proto.RegisterFile("uni/auth.proto", fileDescriptor_auth_ec0d7c7920b062ae) }

var fileDescriptor_auth_ec0d7c7920b062ae = []byte{
	// 204 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2b, 0xcd, 0xcb, 0xd4,
	0x4f, 0x2c, 0x2d, 0xc9, 0xd0, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0xe2, 0x28, 0xcd, 0xcb, 0x8c,
	0x07, 0xf1, 0x95, 0xec, 0xb8, 0x38, 0x7c, 0xf2, 0xd3, 0x33, 0xf3, 0x82, 0x52, 0x0b, 0x85, 0xa4,
	0xb8, 0x38, 0x4a, 0x8b, 0x53, 0x8b, 0xf2, 0x12, 0x73, 0x53, 0x25, 0x18, 0x15, 0x18, 0x35, 0x38,
	0x83, 0xe0, 0x7c, 0x21, 0x31, 0x2e, 0xb6, 0x82, 0xc4, 0xe2, 0xe2, 0xf2, 0x14, 0x09, 0x26, 0xb0,
	0x0c, 0x94, 0xa7, 0x64, 0xc2, 0xc5, 0x05, 0xd5, 0x5f, 0x90, 0x53, 0x29, 0x24, 0xc0, 0xc5, 0x5c,
	0x94, 0x5a, 0x02, 0xd6, 0xcc, 0x11, 0x04, 0x62, 0x0a, 0x89, 0x70, 0xb1, 0x96, 0xe4, 0x67, 0xa7,
	0xe6, 0x41, 0xb5, 0x41, 0x38, 0x46, 0x2e, 0x5c, 0x2c, 0x8e, 0xa5, 0x25, 0x19, 0x42, 0x36, 0x5c,
	0x3c, 0x30, 0xdb, 0x4b, 0x53, 0x8b, 0x4b, 0x84, 0x84, 0xf4, 0x60, 0x0e, 0xd3, 0x83, 0x89, 0x4b,
	0x89, 0x60, 0x88, 0x15, 0xe4, 0x54, 0x2a, 0x31, 0x18, 0x30, 0x3a, 0x29, 0x71, 0x09, 0x64, 0xe6,
	0xeb, 0xa5, 0x17, 0x15, 0x24, 0x83, 0x94, 0xe8, 0x15, 0x17, 0xe4, 0x17, 0x3b, 0xf1, 0x94, 0xe6,
	0x65, 0x82, 0x8c, 0x0e, 0x00, 0xf9, 0x33, 0x80, 0x31, 0x89, 0x0d, 0xec, 0x61, 0x63, 0x40, 0x00,
	0x00, 0x00, 0xff, 0xff, 0x79, 0xa4, 0x97, 0xdc, 0x02, 0x01, 0x00, 0x00,
}
