// Code generated by protoc-gen-go. DO NOT EDIT.
// source: csgeo.proto

package proto

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type ConsumerGeo struct {
	ConsumerGeo          []*ConsumerGeo_Item `protobuf:"bytes,1,rep,name=ConsumerGeo" json:"ConsumerGeo,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *ConsumerGeo) Reset()         { *m = ConsumerGeo{} }
func (m *ConsumerGeo) String() string { return proto.CompactTextString(m) }
func (*ConsumerGeo) ProtoMessage()    {}
func (*ConsumerGeo) Descriptor() ([]byte, []int) {
	return fileDescriptor_csgeo_9584f557c60be3d0, []int{0}
}
func (m *ConsumerGeo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsumerGeo.Unmarshal(m, b)
}
func (m *ConsumerGeo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsumerGeo.Marshal(b, m, deterministic)
}
func (dst *ConsumerGeo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsumerGeo.Merge(dst, src)
}
func (m *ConsumerGeo) XXX_Size() int {
	return xxx_messageInfo_ConsumerGeo.Size(m)
}
func (m *ConsumerGeo) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsumerGeo.DiscardUnknown(m)
}

var xxx_messageInfo_ConsumerGeo proto.InternalMessageInfo

func (m *ConsumerGeo) GetConsumerGeo() []*ConsumerGeo_Item {
	if m != nil {
		return m.ConsumerGeo
	}
	return nil
}

type ConsumerGeo_Item struct {
	ConsumerId           int64                `protobuf:"varint,1,opt,name=ConsumerId" json:"ConsumerId,omitempty"`
	Timestamp            *timestamp.Timestamp `protobuf:"bytes,2,opt,name=Timestamp" json:"Timestamp,omitempty"`
	Lat                  float64              `protobuf:"fixed64,3,opt,name=Lat" json:"Lat,omitempty"`
	Lng                  float64              `protobuf:"fixed64,4,opt,name=Lng" json:"Lng,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *ConsumerGeo_Item) Reset()         { *m = ConsumerGeo_Item{} }
func (m *ConsumerGeo_Item) String() string { return proto.CompactTextString(m) }
func (*ConsumerGeo_Item) ProtoMessage()    {}
func (*ConsumerGeo_Item) Descriptor() ([]byte, []int) {
	return fileDescriptor_csgeo_9584f557c60be3d0, []int{0, 0}
}
func (m *ConsumerGeo_Item) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsumerGeo_Item.Unmarshal(m, b)
}
func (m *ConsumerGeo_Item) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsumerGeo_Item.Marshal(b, m, deterministic)
}
func (dst *ConsumerGeo_Item) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsumerGeo_Item.Merge(dst, src)
}
func (m *ConsumerGeo_Item) XXX_Size() int {
	return xxx_messageInfo_ConsumerGeo_Item.Size(m)
}
func (m *ConsumerGeo_Item) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsumerGeo_Item.DiscardUnknown(m)
}

var xxx_messageInfo_ConsumerGeo_Item proto.InternalMessageInfo

func (m *ConsumerGeo_Item) GetConsumerId() int64 {
	if m != nil {
		return m.ConsumerId
	}
	return 0
}

func (m *ConsumerGeo_Item) GetTimestamp() *timestamp.Timestamp {
	if m != nil {
		return m.Timestamp
	}
	return nil
}

func (m *ConsumerGeo_Item) GetLat() float64 {
	if m != nil {
		return m.Lat
	}
	return 0
}

func (m *ConsumerGeo_Item) GetLng() float64 {
	if m != nil {
		return m.Lng
	}
	return 0
}

func init() {
	proto.RegisterType((*ConsumerGeo)(nil), "proto.ConsumerGeo")
	proto.RegisterType((*ConsumerGeo_Item)(nil), "proto.ConsumerGeo.Item")
}

func init() { proto.RegisterFile("csgeo.proto", fileDescriptor_csgeo_9584f557c60be3d0) }

var fileDescriptor_csgeo_9584f557c60be3d0 = []byte{
	// 190 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x4e, 0x2e, 0x4e, 0x4f,
	0xcd, 0xd7, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x05, 0x53, 0x52, 0xf2, 0xe9, 0xf9, 0xf9,
	0xe9, 0x39, 0xa9, 0xfa, 0x60, 0x5e, 0x52, 0x69, 0x9a, 0x7e, 0x49, 0x66, 0x6e, 0x6a, 0x71, 0x49,
	0x62, 0x6e, 0x01, 0x44, 0x9d, 0xd2, 0x79, 0x46, 0x2e, 0x6e, 0xe7, 0xfc, 0xbc, 0xe2, 0xd2, 0xdc,
	0xd4, 0x22, 0xf7, 0xd4, 0x7c, 0x21, 0x4b, 0x14, 0xae, 0x04, 0xa3, 0x02, 0xb3, 0x06, 0xb7, 0x91,
	0x38, 0x44, 0xb1, 0x1e, 0x92, 0x8c, 0x9e, 0x67, 0x49, 0x6a, 0x6e, 0x10, 0xb2, 0x5a, 0xa9, 0x16,
	0x46, 0x2e, 0x16, 0x90, 0xa8, 0x90, 0x1c, 0x17, 0x17, 0x4c, 0xdc, 0x33, 0x45, 0x82, 0x51, 0x81,
	0x51, 0x83, 0x39, 0x08, 0x49, 0x44, 0xc8, 0x82, 0x8b, 0x33, 0x04, 0xe6, 0x0c, 0x09, 0x26, 0x05,
	0x46, 0x0d, 0x6e, 0x23, 0x29, 0x3d, 0x88, 0x43, 0xf5, 0x60, 0x0e, 0xd5, 0x83, 0xab, 0x08, 0x42,
	0x28, 0x16, 0x12, 0xe0, 0x62, 0xf6, 0x49, 0x2c, 0x91, 0x60, 0x56, 0x60, 0xd4, 0x60, 0x0c, 0x02,
	0x31, 0xc1, 0x22, 0x79, 0xe9, 0x12, 0x2c, 0x50, 0x91, 0xbc, 0xf4, 0x24, 0x36, 0xb0, 0x11, 0xc6,
	0x80, 0x00, 0x00, 0x00, 0xff, 0xff, 0x2a, 0x28, 0xe9, 0x2e, 0x0f, 0x01, 0x00, 0x00,
}
