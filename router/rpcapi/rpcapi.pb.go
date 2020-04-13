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

type MoveReq struct {
	KL                   []*KeyList `protobuf:"bytes,1,rep,name=KL,proto3" json:"KL,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *MoveReq) Reset()         { *m = MoveReq{} }
func (m *MoveReq) String() string { return proto.CompactTextString(m) }
func (*MoveReq) ProtoMessage()    {}
func (*MoveReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2fac6d73d0553fa, []int{3}
}

func (m *MoveReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MoveReq.Unmarshal(m, b)
}
func (m *MoveReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MoveReq.Marshal(b, m, deterministic)
}
func (m *MoveReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MoveReq.Merge(m, src)
}
func (m *MoveReq) XXX_Size() int {
	return xxx_messageInfo_MoveReq.Size(m)
}
func (m *MoveReq) XXX_DiscardUnknown() {
	xxx_messageInfo_MoveReq.DiscardUnknown(m)
}

var xxx_messageInfo_MoveReq proto.InternalMessageInfo

func (m *MoveReq) GetKL() []*KeyList {
	if m != nil {
		return m.KL
	}
	return nil
}

type KeyList struct {
	Node                 *NodeMeta `protobuf:"bytes,1,opt,name=Node,proto3" json:"Node,omitempty"`
	Keys                 []string  `protobuf:"bytes,2,rep,name=Keys,proto3" json:"Keys,omitempty"`
	XXX_NoUnkeyedLiteral struct{}  `json:"-"`
	XXX_unrecognized     []byte    `json:"-"`
	XXX_sizecache        int32     `json:"-"`
}

func (m *KeyList) Reset()         { *m = KeyList{} }
func (m *KeyList) String() string { return proto.CompactTextString(m) }
func (*KeyList) ProtoMessage()    {}
func (*KeyList) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2fac6d73d0553fa, []int{4}
}

func (m *KeyList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_KeyList.Unmarshal(m, b)
}
func (m *KeyList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_KeyList.Marshal(b, m, deterministic)
}
func (m *KeyList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_KeyList.Merge(m, src)
}
func (m *KeyList) XXX_Size() int {
	return xxx_messageInfo_KeyList.Size(m)
}
func (m *KeyList) XXX_DiscardUnknown() {
	xxx_messageInfo_KeyList.DiscardUnknown(m)
}

var xxx_messageInfo_KeyList proto.InternalMessageInfo

func (m *KeyList) GetNode() *NodeMeta {
	if m != nil {
		return m.Node
	}
	return nil
}

func (m *KeyList) GetKeys() []string {
	if m != nil {
		return m.Keys
	}
	return nil
}

type NodeMeta struct {
	ID                   string       `protobuf:"bytes,1,opt,name=ID,proto3" json:"ID,omitempty"`
	Address              string       `protobuf:"bytes,2,opt,name=Address,proto3" json:"Address,omitempty"`
	RPCAddress           string       `protobuf:"bytes,3,opt,name=RPCAddress,proto3" json:"RPCAddress,omitempty"`
	Power                float64      `protobuf:"fixed64,4,opt,name=Power,proto3" json:"Power,omitempty"`
	Capacity             float64      `protobuf:"fixed64,5,opt,name=Capacity,proto3" json:"Capacity,omitempty"`
	Check                *HealthCheck `protobuf:"bytes,6,opt,name=Check,proto3" json:"Check,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *NodeMeta) Reset()         { *m = NodeMeta{} }
func (m *NodeMeta) String() string { return proto.CompactTextString(m) }
func (*NodeMeta) ProtoMessage()    {}
func (*NodeMeta) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2fac6d73d0553fa, []int{5}
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

func (m *NodeMeta) GetCheck() *HealthCheck {
	if m != nil {
		return m.Check
	}
	return nil
}

type HealthCheck struct {
	Timeout                        string   `protobuf:"bytes,1,opt,name=Timeout,proto3" json:"Timeout,omitempty"`
	DeregisterCriticalServiceAfter string   `protobuf:"bytes,2,opt,name=DeregisterCriticalServiceAfter,proto3" json:"DeregisterCriticalServiceAfter,omitempty"`
	XXX_NoUnkeyedLiteral           struct{} `json:"-"`
	XXX_unrecognized               []byte   `json:"-"`
	XXX_sizecache                  int32    `json:"-"`
}

func (m *HealthCheck) Reset()         { *m = HealthCheck{} }
func (m *HealthCheck) String() string { return proto.CompactTextString(m) }
func (*HealthCheck) ProtoMessage()    {}
func (*HealthCheck) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2fac6d73d0553fa, []int{6}
}

func (m *HealthCheck) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HealthCheck.Unmarshal(m, b)
}
func (m *HealthCheck) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HealthCheck.Marshal(b, m, deterministic)
}
func (m *HealthCheck) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HealthCheck.Merge(m, src)
}
func (m *HealthCheck) XXX_Size() int {
	return xxx_messageInfo_HealthCheck.Size(m)
}
func (m *HealthCheck) XXX_DiscardUnknown() {
	xxx_messageInfo_HealthCheck.DiscardUnknown(m)
}

var xxx_messageInfo_HealthCheck proto.InternalMessageInfo

func (m *HealthCheck) GetTimeout() string {
	if m != nil {
		return m.Timeout
	}
	return ""
}

func (m *HealthCheck) GetDeregisterCriticalServiceAfter() string {
	if m != nil {
		return m.DeregisterCriticalServiceAfter
	}
	return ""
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
	return fileDescriptor_b2fac6d73d0553fa, []int{7}
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
	return fileDescriptor_b2fac6d73d0553fa, []int{8}
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
	return fileDescriptor_b2fac6d73d0553fa, []int{9}
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

type Ping struct {
	NodeID               string   `protobuf:"bytes,1,opt,name=NodeID,proto3" json:"NodeID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Ping) Reset()         { *m = Ping{} }
func (m *Ping) String() string { return proto.CompactTextString(m) }
func (*Ping) ProtoMessage()    {}
func (*Ping) Descriptor() ([]byte, []int) {
	return fileDescriptor_b2fac6d73d0553fa, []int{10}
}

func (m *Ping) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Ping.Unmarshal(m, b)
}
func (m *Ping) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Ping.Marshal(b, m, deterministic)
}
func (m *Ping) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Ping.Merge(m, src)
}
func (m *Ping) XXX_Size() int {
	return xxx_messageInfo_Ping.Size(m)
}
func (m *Ping) XXX_DiscardUnknown() {
	xxx_messageInfo_Ping.DiscardUnknown(m)
}

var xxx_messageInfo_Ping proto.InternalMessageInfo

func (m *Ping) GetNodeID() string {
	if m != nil {
		return m.NodeID
	}
	return ""
}

func init() {
	proto.RegisterType((*KeyValue)(nil), "rpcapi.KeyValue")
	proto.RegisterType((*Empty)(nil), "rpcapi.Empty")
	proto.RegisterType((*KeyReq)(nil), "rpcapi.KeyReq")
	proto.RegisterType((*MoveReq)(nil), "rpcapi.MoveReq")
	proto.RegisterType((*KeyList)(nil), "rpcapi.KeyList")
	proto.RegisterType((*NodeMeta)(nil), "rpcapi.NodeMeta")
	proto.RegisterType((*HealthCheck)(nil), "rpcapi.HealthCheck")
	proto.RegisterType((*NodeMetas)(nil), "rpcapi.NodeMetas")
	proto.RegisterType((*ExploreRes)(nil), "rpcapi.ExploreRes")
	proto.RegisterType((*KeyValues)(nil), "rpcapi.KeyValues")
	proto.RegisterType((*Ping)(nil), "rpcapi.Ping")
}

func init() {
	proto.RegisterFile("rpcapi.proto", fileDescriptor_b2fac6d73d0553fa)
}

var fileDescriptor_b2fac6d73d0553fa = []byte{
	// 546 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x54, 0xdd, 0x6e, 0xd3, 0x4c,
	0x10, 0x95, 0x9d, 0x5f, 0x4f, 0xd2, 0x7e, 0xfd, 0x06, 0x84, 0xac, 0x5c, 0x04, 0xcb, 0x42, 0x28,
	0x4d, 0xa5, 0x54, 0xa4, 0x4f, 0x50, 0x9c, 0xa2, 0x56, 0x4e, 0x91, 0xb5, 0x45, 0xbd, 0xdf, 0x3a,
	0x43, 0xbb, 0x90, 0xd4, 0xee, 0x7a, 0x1b, 0xf0, 0x63, 0xf1, 0x20, 0xbc, 0x13, 0xf2, 0xda, 0x1b,
	0x82, 0x8d, 0xc4, 0xdd, 0x9e, 0x39, 0x67, 0xc7, 0x33, 0x67, 0x66, 0x0d, 0x43, 0x99, 0xc6, 0x3c,
	0x15, 0xb3, 0x54, 0x26, 0x2a, 0xc1, 0x6e, 0x89, 0xfc, 0x39, 0xf4, 0x43, 0xca, 0x6f, 0xf9, 0xfa,
	0x99, 0xf0, 0x08, 0x5a, 0x21, 0xe5, 0xae, 0xe5, 0x59, 0x13, 0x87, 0x15, 0x47, 0x7c, 0x09, 0x1d,
	0x4d, 0xb9, 0xb6, 0x67, 0x4d, 0x86, 0xac, 0x04, 0x7e, 0x0f, 0x3a, 0x17, 0x9b, 0x54, 0xe5, 0xfe,
	0x08, 0xba, 0x21, 0xe5, 0x8c, 0x9e, 0x9a, 0x57, 0xfd, 0x29, 0xf4, 0xae, 0x93, 0x2d, 0x15, 0xe4,
	0x6b, 0xb0, 0xc3, 0xa5, 0x6b, 0x79, 0xad, 0xc9, 0x60, 0xfe, 0xdf, 0xac, 0x2a, 0x23, 0xa4, 0x7c,
	0x29, 0x32, 0xc5, 0xec, 0x70, 0xe9, 0x07, 0xd0, 0xab, 0x20, 0xbe, 0x81, 0xf6, 0xc7, 0x64, 0x45,
	0x3a, 0xd3, 0x60, 0x7e, 0x64, 0xd4, 0x45, 0xec, 0x9a, 0x14, 0x67, 0x9a, 0x45, 0x84, 0x76, 0x48,
	0x79, 0xe6, 0xda, 0x5e, 0x6b, 0xe2, 0x30, 0x7d, 0xf6, 0x7f, 0x58, 0xd0, 0x37, 0x32, 0x3c, 0x04,
	0xfb, 0x6a, 0x51, 0x95, 0x63, 0x5f, 0x2d, 0xd0, 0x85, 0xde, 0xf9, 0x6a, 0x25, 0x29, 0xcb, 0x74,
	0x2b, 0x0e, 0x33, 0x10, 0xc7, 0x00, 0x2c, 0x0a, 0x0c, 0xd9, 0xd2, 0xe4, 0x5e, 0xa4, 0xb0, 0x20,
	0x4a, 0xbe, 0x91, 0x74, 0xdb, 0x9e, 0x35, 0xb1, 0x58, 0x09, 0x70, 0x04, 0xfd, 0x80, 0xa7, 0x3c,
	0x16, 0x2a, 0x77, 0x3b, 0x9a, 0xd8, 0x61, 0x3c, 0x86, 0x4e, 0xf0, 0x40, 0xf1, 0x57, 0xb7, 0xab,
	0x7b, 0x78, 0x61, 0x7a, 0xb8, 0x24, 0xbe, 0x56, 0x0f, 0x9a, 0x62, 0xa5, 0xc2, 0x4f, 0x60, 0xb0,
	0x17, 0x2d, 0xaa, 0xfc, 0x24, 0x36, 0x94, 0x3c, 0xab, 0xaa, 0x74, 0x03, 0xf1, 0x03, 0x8c, 0x17,
	0x24, 0xe9, 0x5e, 0x64, 0x8a, 0x64, 0x20, 0x85, 0x12, 0x31, 0x5f, 0xdf, 0x90, 0xdc, 0x8a, 0x98,
	0xce, 0x3f, 0x2b, 0x92, 0x55, 0x5b, 0xff, 0x50, 0xf9, 0x67, 0xe0, 0x18, 0x8f, 0x32, 0x7c, 0x0b,
	0x1d, 0x7d, 0xa8, 0x46, 0xd3, 0x34, 0xbb, 0xa4, 0x7d, 0x0f, 0xe0, 0xe2, 0x7b, 0xba, 0x4e, 0x24,
	0x31, 0xca, 0x76, 0xde, 0x5b, 0x7b, 0xde, 0x9f, 0x82, 0x63, 0xb6, 0x28, 0x43, 0x1f, 0x5a, 0xe1,
	0x6d, 0x23, 0xa9, 0xe1, 0x59, 0x41, 0xfa, 0x63, 0x68, 0x47, 0xe2, 0xf1, 0x1e, 0x5f, 0x41, 0xb7,
	0xf8, 0xda, 0x6e, 0x56, 0x15, 0x9a, 0xff, 0xb4, 0xa1, 0xc7, 0xa2, 0x40, 0x0f, 0xfb, 0x04, 0xfa,
	0x2c, 0x0a, 0x6e, 0x54, 0x22, 0x09, 0x1b, 0xe9, 0x46, 0x07, 0x26, 0xa2, 0x57, 0x12, 0xdf, 0xc1,
	0x81, 0x11, 0x47, 0x5c, 0xc8, 0x0c, 0xff, 0xaf, 0xdf, 0xc8, 0xea, 0x57, 0x66, 0x7a, 0x03, 0x18,
	0xc5, 0x24, 0xb6, 0x84, 0x87, 0x7b, 0x7a, 0x46, 0x4f, 0xa3, 0xc6, 0x17, 0x71, 0x0a, 0x8e, 0xd6,
	0x6f, 0x92, 0xbf, 0xc8, 0x6b, 0xb9, 0x4f, 0x75, 0xee, 0xca, 0x3d, 0xfc, 0x93, 0x1c, 0xe1, 0x0e,
	0xfe, 0x76, 0x77, 0xaa, 0xfb, 0xd6, 0x3b, 0x5c, 0x53, 0x37, 0xc6, 0x83, 0xc7, 0xa5, 0xb6, 0x28,
	0x63, 0xf7, 0xac, 0xaa, 0x37, 0x57, 0xab, 0x63, 0xfe, 0x05, 0x06, 0x2c, 0x0a, 0xde, 0xf3, 0x35,
	0x7f, 0x8c, 0x49, 0xe2, 0x4c, 0x43, 0x56, 0x6d, 0x0a, 0x36, 0x52, 0xd7, 0xdb, 0x38, 0x81, 0x21,
	0x8b, 0x82, 0x4b, 0xe2, 0x52, 0xdd, 0x11, 0x57, 0x38, 0x34, 0x74, 0x31, 0xc4, 0x9a, 0xf8, 0xae,
	0xab, 0xff, 0x30, 0x67, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0x39, 0x31, 0xca, 0x55, 0x71, 0x04,
	0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// RPCNodeClient is the client API for RPCNode service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RPCNodeClient interface {
	RPCStore(ctx context.Context, in *KeyValue, opts ...grpc.CallOption) (*Empty, error)
	RPCStorePairs(ctx context.Context, in *KeyValues, opts ...grpc.CallOption) (*Empty, error)
	RPCReceive(ctx context.Context, in *KeyReq, opts ...grpc.CallOption) (*KeyValue, error)
	RPCRemove(ctx context.Context, in *KeyReq, opts ...grpc.CallOption) (*Empty, error)
	RPCExplore(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ExploreRes, error)
	RPCMeta(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*NodeMeta, error)
	RPCMove(ctx context.Context, in *MoveReq, opts ...grpc.CallOption) (*Empty, error)
}

type rPCNodeClient struct {
	cc grpc.ClientConnInterface
}

func NewRPCNodeClient(cc grpc.ClientConnInterface) RPCNodeClient {
	return &rPCNodeClient{cc}
}

func (c *rPCNodeClient) RPCStore(ctx context.Context, in *KeyValue, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCNode/RPCStore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCNodeClient) RPCStorePairs(ctx context.Context, in *KeyValues, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCNode/RPCStorePairs", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCNodeClient) RPCReceive(ctx context.Context, in *KeyReq, opts ...grpc.CallOption) (*KeyValue, error) {
	out := new(KeyValue)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCNode/RPCReceive", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCNodeClient) RPCRemove(ctx context.Context, in *KeyReq, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCNode/RPCRemove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCNodeClient) RPCExplore(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*ExploreRes, error) {
	out := new(ExploreRes)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCNode/RPCExplore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCNodeClient) RPCMeta(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*NodeMeta, error) {
	out := new(NodeMeta)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCNode/RPCMeta", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCNodeClient) RPCMove(ctx context.Context, in *MoveReq, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCNode/RPCMove", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RPCNodeServer is the server API for RPCNode service.
type RPCNodeServer interface {
	RPCStore(context.Context, *KeyValue) (*Empty, error)
	RPCStorePairs(context.Context, *KeyValues) (*Empty, error)
	RPCReceive(context.Context, *KeyReq) (*KeyValue, error)
	RPCRemove(context.Context, *KeyReq) (*Empty, error)
	RPCExplore(context.Context, *Empty) (*ExploreRes, error)
	RPCMeta(context.Context, *Empty) (*NodeMeta, error)
	RPCMove(context.Context, *MoveReq) (*Empty, error)
}

// UnimplementedRPCNodeServer can be embedded to have forward compatible implementations.
type UnimplementedRPCNodeServer struct {
}

func (*UnimplementedRPCNodeServer) RPCStore(ctx context.Context, req *KeyValue) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCStore not implemented")
}
func (*UnimplementedRPCNodeServer) RPCStorePairs(ctx context.Context, req *KeyValues) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCStorePairs not implemented")
}
func (*UnimplementedRPCNodeServer) RPCReceive(ctx context.Context, req *KeyReq) (*KeyValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCReceive not implemented")
}
func (*UnimplementedRPCNodeServer) RPCRemove(ctx context.Context, req *KeyReq) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCRemove not implemented")
}
func (*UnimplementedRPCNodeServer) RPCExplore(ctx context.Context, req *Empty) (*ExploreRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCExplore not implemented")
}
func (*UnimplementedRPCNodeServer) RPCMeta(ctx context.Context, req *Empty) (*NodeMeta, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCMeta not implemented")
}
func (*UnimplementedRPCNodeServer) RPCMove(ctx context.Context, req *MoveReq) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCMove not implemented")
}

func RegisterRPCNodeServer(s *grpc.Server, srv RPCNodeServer) {
	s.RegisterService(&_RPCNode_serviceDesc, srv)
}

func _RPCNode_RPCStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCNodeServer).RPCStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCNode/RPCStore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCNodeServer).RPCStore(ctx, req.(*KeyValue))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCNode_RPCStorePairs_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyValues)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCNodeServer).RPCStorePairs(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCNode/RPCStorePairs",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCNodeServer).RPCStorePairs(ctx, req.(*KeyValues))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCNode_RPCReceive_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCNodeServer).RPCReceive(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCNode/RPCReceive",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCNodeServer).RPCReceive(ctx, req.(*KeyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCNode_RPCRemove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(KeyReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCNodeServer).RPCRemove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCNode/RPCRemove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCNodeServer).RPCRemove(ctx, req.(*KeyReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCNode_RPCExplore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCNodeServer).RPCExplore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCNode/RPCExplore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCNodeServer).RPCExplore(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCNode_RPCMeta_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCNodeServer).RPCMeta(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCNode/RPCMeta",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCNodeServer).RPCMeta(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCNode_RPCMove_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MoveReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCNodeServer).RPCMove(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCNode/RPCMove",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCNodeServer).RPCMove(ctx, req.(*MoveReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _RPCNode_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpcapi.RPCNode",
	HandlerType: (*RPCNodeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RPCStore",
			Handler:    _RPCNode_RPCStore_Handler,
		},
		{
			MethodName: "RPCStorePairs",
			Handler:    _RPCNode_RPCStorePairs_Handler,
		},
		{
			MethodName: "RPCReceive",
			Handler:    _RPCNode_RPCReceive_Handler,
		},
		{
			MethodName: "RPCRemove",
			Handler:    _RPCNode_RPCRemove_Handler,
		},
		{
			MethodName: "RPCExplore",
			Handler:    _RPCNode_RPCExplore_Handler,
		},
		{
			MethodName: "RPCMeta",
			Handler:    _RPCNode_RPCMeta_Handler,
		},
		{
			MethodName: "RPCMove",
			Handler:    _RPCNode_RPCMove_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpcapi.proto",
}

// RPCBalancerClient is the client API for RPCBalancer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type RPCBalancerClient interface {
	RPCRegister(ctx context.Context, in *NodeMeta, opts ...grpc.CallOption) (*Empty, error)
	RPCHeartbeat(ctx context.Context, in *Ping, opts ...grpc.CallOption) (*Empty, error)
}

type rPCBalancerClient struct {
	cc grpc.ClientConnInterface
}

func NewRPCBalancerClient(cc grpc.ClientConnInterface) RPCBalancerClient {
	return &rPCBalancerClient{cc}
}

func (c *rPCBalancerClient) RPCRegister(ctx context.Context, in *NodeMeta, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCBalancer/RPCRegister", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rPCBalancerClient) RPCHeartbeat(ctx context.Context, in *Ping, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/rpcapi.RPCBalancer/RPCHeartbeat", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RPCBalancerServer is the server API for RPCBalancer service.
type RPCBalancerServer interface {
	RPCRegister(context.Context, *NodeMeta) (*Empty, error)
	RPCHeartbeat(context.Context, *Ping) (*Empty, error)
}

// UnimplementedRPCBalancerServer can be embedded to have forward compatible implementations.
type UnimplementedRPCBalancerServer struct {
}

func (*UnimplementedRPCBalancerServer) RPCRegister(ctx context.Context, req *NodeMeta) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCRegister not implemented")
}
func (*UnimplementedRPCBalancerServer) RPCHeartbeat(ctx context.Context, req *Ping) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RPCHeartbeat not implemented")
}

func RegisterRPCBalancerServer(s *grpc.Server, srv RPCBalancerServer) {
	s.RegisterService(&_RPCBalancer_serviceDesc, srv)
}

func _RPCBalancer_RPCRegister_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(NodeMeta)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCBalancerServer).RPCRegister(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCBalancer/RPCRegister",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCBalancerServer).RPCRegister(ctx, req.(*NodeMeta))
	}
	return interceptor(ctx, in, info, handler)
}

func _RPCBalancer_RPCHeartbeat_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Ping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RPCBalancerServer).RPCHeartbeat(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/rpcapi.RPCBalancer/RPCHeartbeat",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RPCBalancerServer).RPCHeartbeat(ctx, req.(*Ping))
	}
	return interceptor(ctx, in, info, handler)
}

var _RPCBalancer_serviceDesc = grpc.ServiceDesc{
	ServiceName: "rpcapi.RPCBalancer",
	HandlerType: (*RPCBalancerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RPCRegister",
			Handler:    _RPCBalancer_RPCRegister_Handler,
		},
		{
			MethodName: "RPCHeartbeat",
			Handler:    _RPCBalancer_RPCHeartbeat_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "rpcapi.proto",
}