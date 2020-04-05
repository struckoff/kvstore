// Code generated by protoc-gen-go. DO NOT EDIT.
// source: rpcapi.proto

package rpcapi

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type KeyValue struct {
	Key                  string   `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	Value                []byte   `protobuf:"bytes,2,opt,name=Value,proto3" json:"Value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KeyValue) Reset()         { *m = KeyValue{} }
func (m *KeyValue) String() string { return proto.CompactTextString(m) }
func (*KeyValue) ProtoMessage()    {}
func (*KeyValue) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2fac6d73d0553fa, []int{0}
}

func (m *KeyValue) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeyValue.Unmarshal(m, b)
}
func (m *KeyValue) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeyValue.Marshal(b, m, deterministic)
}
func (m *KeyValue) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeyValue.Merge(m, src)
}
func (m *KeyValue) XXX_Size() int {
	return xxx_messageInfo_KeyValue.Size(m)
}
func (m *KeyValue) XXX_DiscardUnknown() {
	xxx_messageInfo_KeyValue.DiscardUnknown(m)
}

var xxx_messageInfo_KeyValue proto.InternalMessageInfo

func (m *KeyValue) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

func (m *KeyValue) GetValue() []byte {
	if m != nil {
		return m.Value
	}
	return nil
}

type Empty struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Empty) Reset()         { *m = Empty{} }
func (m *Empty) String() string { return proto.CompactTextString(m) }
func (*Empty) ProtoMessage()    {}
func (*Empty) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2fac6d73d0553fa, []int{1}
}

func (m *Empty) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Empty.Unmarshal(m, b)
}
func (m *Empty) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Empty.Marshal(b, m, deterministic)
}
func (m *Empty) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Empty.Merge(m, src)
}
func (m *Empty) XXX_Size() int {
	return xxx_messageInfo_Empty.Size(m)
}
func (m *Empty) XXX_DiscardUnknown() {
	xxx_messageInfo_Empty.DiscardUnknown(m)
}

var xxx_messageInfo_Empty proto.InternalMessageInfo

type KeyReq struct {
	Key                  string   `protobuf:"bytes,1,opt,name=Key,proto3" json:"Key,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *KeyReq) Reset()         { *m = KeyReq{} }
func (m *KeyReq) String() string { return proto.CompactTextString(m) }
func (*KeyReq) ProtoMessage()    {}
func (*KeyReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2fac6d73d0553fa, []int{2}
}

func (m *KeyReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeyReq.Unmarshal(m, b)
}
func (m *KeyReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeyReq.Marshal(b, m, deterministic)
}
func (m *KeyReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeyReq.Merge(m, src)
}
func (m *KeyReq) XXX_Size() int {
	return xxx_messageInfo_KeyReq.Size(m)
}
func (m *KeyReq) XXX_DiscardUnknown() {
	xxx_messageInfo_KeyReq.DiscardUnknown(m)
}

var xxx_messageInfo_KeyReq proto.InternalMessageInfo

func (m *KeyReq) GetKey() string {
	if m != nil {
		return m.Key
	}
	return ""
}

type NodeMeta struct {
	ID                   string   `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Address              string   `protobuf:"bytes,2,opt,name=Address,proto3" json:"Address,omitempty"`
	RPCAddress           string   `protobuf:"bytes,3,opt,name=RPCAddress,proto3" json:"RPCAddress,omitempty"`
	Power                float64  `protobuf:"fixed64,4,opt,name=Power,proto3" json:"Power,omitempty"`
	Capacity             float64  `protobuf:"fixed64,5,opt,name=Capacity,proto3" json:"Capacity,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeMeta) Reset()         { *m = NodeMeta{} }
func (m *NodeMeta) String() string { return proto.CompactTextString(m) }
func (*NodeMeta) ProtoMessage()    {}
func (*NodeMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2fac6d73d0553fa, []int{3}
}

func (m *NodeMeta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeMeta.Unmarshal(m, b)
}
func (m *NodeMeta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeMeta.Marshal(b, m, deterministic)
}
func (m *NodeMeta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeMeta.Merge(m, src)
}
func (m *NodeMeta) XXX_Size() int {
	return xxx_messageInfo_NodeMeta.Size(m)
}
func (m *NodeMeta) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeMeta.DiscardUnknown(m)
}

var xxx_messageInfo_NodeMeta proto.InternalMessageInfo

func (m *NodeMeta) GetID() string {
	if m != nil {
		return m.ID
	}
	return ""
}

func (m *NodeMeta) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *NodeMeta) GetRPCAddress() string {
	if m != nil {
		return m.RPCAddress
	}
	return ""
}

func (m *NodeMeta) GetPower() float64 {
	if m != nil {
		return m.Power
	}
	return 0
}

func (m *NodeMeta) GetCapacity() float64 {
	if m != nil {
		return m.Capacity
	}
	return 0
}

type NodeMetas struct {
	Metas                []*NodeMeta `protobuf:"bytes,1,rep,name=Metas,proto3" json:"Metas,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *NodeMetas) Reset()         { *m = NodeMetas{} }
func (m *NodeMetas) String() string { return proto.CompactTextString(m) }
func (*NodeMetas) ProtoMessage()    {}
func (*NodeMetas) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2fac6d73d0553fa, []int{4}
}

func (m *NodeMetas) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeMetas.Unmarshal(m, b)
}
func (m *NodeMetas) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeMetas.Marshal(b, m, deterministic)
}
func (m *NodeMetas) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeMetas.Merge(m, src)
}
func (m *NodeMetas) XXX_Size() int {
	return xxx_messageInfo_NodeMetas.Size(m)
}
func (m *NodeMetas) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeMetas.DiscardUnknown(m)
}

var xxx_messageInfo_NodeMetas proto.InternalMessageInfo

func (m *NodeMetas) GetMetas() []*NodeMeta {
	if m != nil {
		return m.Metas
	}
	return nil
}

type ExploreRes struct {
	Keys                 []string `protobuf:"bytes,1,rep,name=Keys,proto3" json:"Keys,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ExploreRes) Reset()         { *m = ExploreRes{} }
func (m *ExploreRes) String() string { return proto.CompactTextString(m) }
func (*ExploreRes) ProtoMessage()    {}
func (*ExploreRes) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2fac6d73d0553fa, []int{5}
}

func (m *ExploreRes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExploreRes.Unmarshal(m, b)
}
func (m *ExploreRes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExploreRes.Marshal(b, m, deterministic)
}
func (m *ExploreRes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExploreRes.Merge(m, src)
}
func (m *ExploreRes) XXX_Size() int {
	return xxx_messageInfo_ExploreRes.Size(m)
}
func (m *ExploreRes) XXX_DiscardUnknown() {
	xxx_messageInfo_ExploreRes.DiscardUnknown(m)
}

var xxx_messageInfo_ExploreRes proto.InternalMessageInfo

func (m *ExploreRes) GetKeys() []string {
	if m != nil {
		return m.Keys
	}
	return nil
}

type KeyValues struct {
	KVs                  []*KeyValue `protobuf:"bytes,1,rep,name=KVs,proto3" json:"KVs,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *KeyValues) Reset()         { *m = KeyValues{} }
func (m *KeyValues) String() string { return proto.CompactTextString(m) }
func (*KeyValues) ProtoMessage()    {}
func (*KeyValues) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2fac6d73d0553fa, []int{6}
}

func (m *KeyValues) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeyValues.Unmarshal(m, b)
}
func (m *KeyValues) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeyValues.Marshal(b, m, deterministic)
}
func (m *KeyValues) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeyValues.Merge(m, src)
}
func (m *KeyValues) XXX_Size() int {
	return xxx_messageInfo_KeyValues.Size(m)
}
func (m *KeyValues) XXX_DiscardUnknown() {
	xxx_messageInfo_KeyValues.DiscardUnknown(m)
}

var xxx_messageInfo_KeyValues proto.InternalMessageInfo

func (m *KeyValues) GetKVs() []*KeyValue {
	if m != nil {
		return m.KVs
	}
	return nil
}

func init() {
	proto.RegisterType((*KeyValue)(nil), "rpcapi.KeyValue")
	proto.RegisterType((*Empty)(nil), "rpcapi.Empty")
	proto.RegisterType((*KeyReq)(nil), "rpcapi.KeyReq")
	proto.RegisterType((*NodeMeta)(nil), "rpcapi.NodeMeta")
	proto.RegisterType((*NodeMetas)(nil), "rpcapi.NodeMetas")
	proto.RegisterType((*ExploreRes)(nil), "rpcapi.ExploreRes")
	proto.RegisterType((*KeyValues)(nil), "rpcapi.KeyValues")
}

func init() {
	proto.RegisterFile("rpcapi.proto", fileDescriptor_b2fac6d73d0553fa)
}

var fileDescriptor_b2fac6d73d0553fa = []byte{
	// 362 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x92, 0xcf, 0x6e, 0xaa, 0x40,
	0x14, 0xc6, 0x03, 0x88, 0xc2, 0xf1, 0x4f, 0xbc, 0x27, 0x77, 0x41, 0x58, 0xdc, 0x90, 0x59, 0xdc,
	0x10, 0x6f, 0xa2, 0xb9, 0xfa, 0x04, 0x0d, 0xba, 0x30, 0xb4, 0x0d, 0x99, 0x26, 0xee, 0xa9, 0x9e,
	0x05, 0x89, 0x16, 0x9c, 0xa1, 0x7f, 0x78, 0x81, 0x3e, 0x4a, 0x9f, 0xb3, 0x71, 0x70, 0xac, 0x85,
	0xee, 0xce, 0x37, 0xbf, 0xef, 0x3b, 0x4c, 0x3e, 0x06, 0x06, 0xa2, 0xd8, 0xa6, 0x45, 0x36, 0x2d,
	0x44, 0x5e, 0xe6, 0xd8, 0xad, 0x15, 0x9b, 0x83, 0x13, 0x53, 0xb5, 0x49, 0xf7, 0xcf, 0x84, 0x63,
	0xb0, 0x62, 0xaa, 0x3c, 0x23, 0x30, 0x42, 0x97, 0x9f, 0x46, 0xfc, 0x0d, 0xb6, 0x42, 0x9e, 0x19,
	0x18, 0xe1, 0x80, 0xd7, 0x82, 0xf5, 0xc0, 0x5e, 0x1d, 0x8a, 0xb2, 0x62, 0x3e, 0x74, 0x63, 0xaa,
	0x38, 0x1d, 0xdb, 0x51, 0xf6, 0x6e, 0x80, 0x73, 0x9f, 0xef, 0xe8, 0x8e, 0xca, 0x14, 0x47, 0x60,
	0xae, 0x97, 0x67, 0x6a, 0xae, 0x97, 0xe8, 0x41, 0xef, 0x66, 0xb7, 0x13, 0x24, 0xa5, 0xda, 0xec,
	0x72, 0x2d, 0xf1, 0x0f, 0x00, 0x4f, 0x22, 0x0d, 0x2d, 0x05, 0xaf, 0x4e, 0x4e, 0x37, 0x4a, 0xf2,
	0x57, 0x12, 0x5e, 0x27, 0x30, 0x42, 0x83, 0xd7, 0x02, 0x7d, 0x70, 0xa2, 0xb4, 0x48, 0xb7, 0x59,
	0x59, 0x79, 0xb6, 0x02, 0x17, 0xcd, 0x16, 0xe0, 0xea, 0x7b, 0x48, 0xfc, 0x0b, 0xb6, 0x1a, 0x3c,
	0x23, 0xb0, 0xc2, 0xfe, 0x7c, 0x3c, 0x3d, 0x97, 0xa2, 0x1d, 0xbc, 0xc6, 0x2c, 0x00, 0x58, 0xbd,
	0x15, 0xfb, 0x5c, 0x10, 0x27, 0x89, 0x08, 0x9d, 0x98, 0xaa, 0x3a, 0xe4, 0x72, 0x35, 0xb3, 0x19,
	0xb8, 0xba, 0x38, 0x89, 0x0c, 0xac, 0x78, 0xd3, 0x5a, 0xaa, 0x39, 0x3f, 0xc1, 0xf9, 0x87, 0x09,
	0x7d, 0x9e, 0x44, 0xb7, 0x99, 0x2c, 0xe9, 0x89, 0x04, 0xfe, 0x03, 0x87, 0x27, 0xd1, 0x43, 0x99,
	0x0b, 0xc2, 0x56, 0xc4, 0x1f, 0xea, 0x13, 0xd5, 0x34, 0xfe, 0x87, 0xa1, 0x36, 0x27, 0x69, 0x26,
	0x24, 0xfe, 0x6a, 0x26, 0x64, 0x33, 0x32, 0x55, 0x4d, 0x72, 0xda, 0x52, 0xf6, 0x42, 0x38, 0xba,
	0xf2, 0x73, 0x3a, 0xfa, 0xad, 0x2f, 0xe2, 0x04, 0x5c, 0xe5, 0x3f, 0xe4, 0x3f, 0xd8, 0x1b, 0xbb,
	0x67, 0x6a, 0xf7, 0xb9, 0x21, 0xfc, 0x0e, 0x7d, 0xbc, 0xc8, 0xaf, 0x06, 0x27, 0xd0, 0xe3, 0x49,
	0xa4, 0xde, 0x42, 0xc3, 0xdd, 0xfa, 0x05, 0x8f, 0x5d, 0xf5, 0x42, 0x17, 0x9f, 0x01, 0x00, 0x00,
	0xff, 0xff, 0xd9, 0xe8, 0x41, 0xe7, 0xb1, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RPCListenerClient is the client API for RPCListener service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RPCListenerClient interface {
	RPCStore(ctx context.Context, in *KeyValue, opts ...grpc.CallOption) (*Empty, error)
	RPCStorePairs(ctx context.Context, in *KeyValues, opts ...grpc.CallOption) (*Empty, error)
	RPCReceive(ctx context.Context, in *KeyReq, opts ...grpc.CallOption) (*KeyValue, error)
	RPCRemove(ctx context.Context, in *KeyReq, opts ...grpc.CallOption) (*Empty, error)
	RPCExplore(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ExploreRes, error)
	RPCMeta(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*NodeMeta, error)
}

type rPCListenerClient struct {
	cc grpc.ClientConnInterface
}

func NewRPCListenerClient(cc grpc.ClientConnInterface) RPCListenerClient {
	return &rPCListenerClient{cc}
}

func (c *rPCListenerClient) RPCStore(ctx context.Context, in *KeyValue, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCListener/RPCStore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCListenerClient) RPCStorePairs(ctx context.Context, in *KeyValues, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCListener/RPCStorePairs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCListenerClient) RPCReceive(ctx context.Context, in *KeyReq, opts ...grpc.CallOption) (*KeyValue, error) {
	out := new(KeyValue)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCListener/RPCReceive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCListenerClient) RPCRemove(ctx context.Context, in *KeyReq, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCListener/RPCRemove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCListenerClient) RPCExplore(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ExploreRes, error) {
	out := new(ExploreRes)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCListener/RPCExplore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCListenerClient) RPCMeta(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*NodeMeta, error) {
	out := new(NodeMeta)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCListener/RPCMeta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RPCListenerServer is the server API for RPCListener service.
type RPCListenerServer interface {
	RPCStore(context.Context, *KeyValue) (*Empty, error)
	RPCStorePairs(context.Context, *KeyValues) (*Empty, error)
	RPCReceive(context.Context, *KeyReq) (*KeyValue, error)
	RPCRemove(context.Context, *KeyReq) (*Empty, error)
	RPCExplore(context.Context, *Empty) (*ExploreRes, error)
	RPCMeta(context.Context, *Empty) (*NodeMeta, error)
}

// UnimplementedRPCListenerServer can be embedded to have forward compatible implementations.
type UnimplementedRPCListenerServer struct {
}

func (*UnimplementedRPCListenerServer) RPCStore(ctx context.Context, req *KeyValue) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCStore not implemented")
}
func (*UnimplementedRPCListenerServer) RPCStorePairs(ctx context.Context, req *KeyValues) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCStorePairs not implemented")
}
func (*UnimplementedRPCListenerServer) RPCReceive(ctx context.Context, req *KeyReq) (*KeyValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCReceive not implemented")
}
func (*UnimplementedRPCListenerServer) RPCRemove(ctx context.Context, req *KeyReq) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCRemove not implemented")
}
func (*UnimplementedRPCListenerServer) RPCExplore(ctx context.Context, req *Empty) (*ExploreRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCExplore not implemented")
}
func (*UnimplementedRPCListenerServer) RPCMeta(ctx context.Context, req *Empty) (*NodeMeta, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCMeta not implemented")
}

func RegisterRPCListenerServer(s *grpc.Server, srv RPCListenerServer) {
	s.RegisterService(&_RPCListener_serviceDesc, srv)
}

func _RPCListener_RPCStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCListenerServer).RPCStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCListener/RPCStore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCListenerServer).RPCStore(ctx, req.(*KeyValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCListener_RPCStorePairs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValues)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCListenerServer).RPCStorePairs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCListener/RPCStorePairs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCListenerServer).RPCStorePairs(ctx, req.(*KeyValues))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCListener_RPCReceive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCListenerServer).RPCReceive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCListener/RPCReceive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCListenerServer).RPCReceive(ctx, req.(*KeyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCListener_RPCRemove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCListenerServer).RPCRemove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCListener/RPCRemove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCListenerServer).RPCRemove(ctx, req.(*KeyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCListener_RPCExplore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCListenerServer).RPCExplore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCListener/RPCExplore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCListenerServer).RPCExplore(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCListener_RPCMeta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCListenerServer).RPCMeta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCListener/RPCMeta",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCListenerServer).RPCMeta(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

var _RPCListener_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpcapi.RPCListener",
	HandlerType: (*RPCListenerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RPCStore",
			Handler:    _RPCListener_RPCStore_Handler,
		},
		{
			MethodName: "RPCStorePairs",
			Handler:    _RPCListener_RPCStorePairs_Handler,
		},
		{
			MethodName: "RPCReceive",
			Handler:    _RPCListener_RPCReceive_Handler,
		},
		{
			MethodName: "RPCRemove",
			Handler:    _RPCListener_RPCRemove_Handler,
		},
		{
			MethodName: "RPCExplore",
			Handler:    _RPCListener_RPCExplore_Handler,
		},
		{
			MethodName: "RPCMeta",
			Handler:    _RPCListener_RPCMeta_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpcapi.proto",
}