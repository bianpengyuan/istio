// Code generated by protoc-gen-go. DO NOT EDIT.
// source: mixer_frontend.proto

package mixer

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import types "github.com/gogo/protobuf/types"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type MixerFrontendConfig struct {
	TemplateDescriptors  map[string]string `protobuf:"bytes,1,rep,name=template_descriptors,json=templateDescriptors" json:"template_descriptors,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	AdapterDescriptors   map[string]string `protobuf:"bytes,2,rep,name=adapter_descriptors,json=adapterDescriptors" json:"adapter_descriptors,omitempty" protobuf_key:"bytes,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	InstanceConfig       []*InstanceConfig `protobuf:"bytes,3,rep,name=instance_config,json=instanceConfig" json:"instance_config,omitempty"`
	HandlerConfig        []*HandlerConfig  `protobuf:"bytes,4,rep,name=handler_config,json=handlerConfig" json:"handler_config,omitempty"`
	DispatchSpec         []*DispatchSpec   `protobuf:"bytes,5,rep,name=dispatch_spec,json=dispatchSpec" json:"dispatch_spec,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *MixerFrontendConfig) Reset()         { *m = MixerFrontendConfig{} }
func (m *MixerFrontendConfig) String() string { return proto.CompactTextString(m) }
func (*MixerFrontendConfig) ProtoMessage()    {}
func (*MixerFrontendConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_mixer_frontend_b953796d700de19f, []int{0}
}
func (m *MixerFrontendConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MixerFrontendConfig.Unmarshal(m, b)
}
func (m *MixerFrontendConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MixerFrontendConfig.Marshal(b, m, deterministic)
}
func (dst *MixerFrontendConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MixerFrontendConfig.Merge(dst, src)
}
func (m *MixerFrontendConfig) XXX_Size() int {
	return xxx_messageInfo_MixerFrontendConfig.Size(m)
}
func (m *MixerFrontendConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_MixerFrontendConfig.DiscardUnknown(m)
}

var xxx_messageInfo_MixerFrontendConfig proto.InternalMessageInfo

func (m *MixerFrontendConfig) GetTemplateDescriptors() map[string]string {
	if m != nil {
		return m.TemplateDescriptors
	}
	return nil
}

func (m *MixerFrontendConfig) GetAdapterDescriptors() map[string]string {
	if m != nil {
		return m.AdapterDescriptors
	}
	return nil
}

func (m *MixerFrontendConfig) GetInstanceConfig() []*InstanceConfig {
	if m != nil {
		return m.InstanceConfig
	}
	return nil
}

func (m *MixerFrontendConfig) GetHandlerConfig() []*HandlerConfig {
	if m != nil {
		return m.HandlerConfig
	}
	return nil
}

func (m *MixerFrontendConfig) GetDispatchSpec() []*DispatchSpec {
	if m != nil {
		return m.DispatchSpec
	}
	return nil
}

type InstanceConfig struct {
	Name     string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Template string `protobuf:"bytes,2,opt,name=template" json:"template,omitempty"`
	// This params is a struct of key to expressions.
	Params               *types.Struct `protobuf:"bytes,3,opt,name=params" json:"params,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *InstanceConfig) Reset()         { *m = InstanceConfig{} }
func (m *InstanceConfig) String() string { return proto.CompactTextString(m) }
func (*InstanceConfig) ProtoMessage()    {}
func (*InstanceConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_mixer_frontend_b953796d700de19f, []int{1}
}
func (m *InstanceConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InstanceConfig.Unmarshal(m, b)
}
func (m *InstanceConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InstanceConfig.Marshal(b, m, deterministic)
}
func (dst *InstanceConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InstanceConfig.Merge(dst, src)
}
func (m *InstanceConfig) XXX_Size() int {
	return xxx_messageInfo_InstanceConfig.Size(m)
}
func (m *InstanceConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_InstanceConfig.DiscardUnknown(m)
}

var xxx_messageInfo_InstanceConfig proto.InternalMessageInfo

func (m *InstanceConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *InstanceConfig) GetTemplate() string {
	if m != nil {
		return m.Template
	}
	return ""
}

func (m *InstanceConfig) GetParams() *types.Struct {
	if m != nil {
		return m.Params
	}
	return nil
}

type HandlerConfig struct {
	Name                 string        `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Adapter              string        `protobuf:"bytes,2,opt,name=adapter" json:"adapter,omitempty"`
	Address              string        `protobuf:"bytes,3,opt,name=address" json:"address,omitempty"`
	Params               *types.Struct `protobuf:"bytes,4,opt,name=params" json:"params,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *HandlerConfig) Reset()         { *m = HandlerConfig{} }
func (m *HandlerConfig) String() string { return proto.CompactTextString(m) }
func (*HandlerConfig) ProtoMessage()    {}
func (*HandlerConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_mixer_frontend_b953796d700de19f, []int{2}
}
func (m *HandlerConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_HandlerConfig.Unmarshal(m, b)
}
func (m *HandlerConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_HandlerConfig.Marshal(b, m, deterministic)
}
func (dst *HandlerConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_HandlerConfig.Merge(dst, src)
}
func (m *HandlerConfig) XXX_Size() int {
	return xxx_messageInfo_HandlerConfig.Size(m)
}
func (m *HandlerConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_HandlerConfig.DiscardUnknown(m)
}

var xxx_messageInfo_HandlerConfig proto.InternalMessageInfo

func (m *HandlerConfig) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *HandlerConfig) GetAdapter() string {
	if m != nil {
		return m.Adapter
	}
	return ""
}

func (m *HandlerConfig) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *HandlerConfig) GetParams() *types.Struct {
	if m != nil {
		return m.Params
	}
	return nil
}

type DispatchSpec struct {
	MatchExpression      string   `protobuf:"bytes,1,opt,name=match_expression,json=matchExpression" json:"match_expression,omitempty"`
	Instances            []string `protobuf:"bytes,2,rep,name=instances" json:"instances,omitempty"`
	Handlers             []string `protobuf:"bytes,3,rep,name=handlers" json:"handlers,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DispatchSpec) Reset()         { *m = DispatchSpec{} }
func (m *DispatchSpec) String() string { return proto.CompactTextString(m) }
func (*DispatchSpec) ProtoMessage()    {}
func (*DispatchSpec) Descriptor() ([]byte, []int) {
	return fileDescriptor_mixer_frontend_b953796d700de19f, []int{3}
}
func (m *DispatchSpec) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DispatchSpec.Unmarshal(m, b)
}
func (m *DispatchSpec) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DispatchSpec.Marshal(b, m, deterministic)
}
func (dst *DispatchSpec) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DispatchSpec.Merge(dst, src)
}
func (m *DispatchSpec) XXX_Size() int {
	return xxx_messageInfo_DispatchSpec.Size(m)
}
func (m *DispatchSpec) XXX_DiscardUnknown() {
	xxx_messageInfo_DispatchSpec.DiscardUnknown(m)
}

var xxx_messageInfo_DispatchSpec proto.InternalMessageInfo

func (m *DispatchSpec) GetMatchExpression() string {
	if m != nil {
		return m.MatchExpression
	}
	return ""
}

func (m *DispatchSpec) GetInstances() []string {
	if m != nil {
		return m.Instances
	}
	return nil
}

func (m *DispatchSpec) GetHandlers() []string {
	if m != nil {
		return m.Handlers
	}
	return nil
}

func init() {
	proto.RegisterType((*MixerFrontendConfig)(nil), "mixer.MixerFrontendConfig")
	proto.RegisterMapType((map[string]string)(nil), "mixer.MixerFrontendConfig.AdapterDescriptorsEntry")
	proto.RegisterMapType((map[string]string)(nil), "mixer.MixerFrontendConfig.TemplateDescriptorsEntry")
	proto.RegisterType((*InstanceConfig)(nil), "mixer.InstanceConfig")
	proto.RegisterType((*HandlerConfig)(nil), "mixer.HandlerConfig")
	proto.RegisterType((*DispatchSpec)(nil), "mixer.DispatchSpec")
}

func init() {
	proto.RegisterFile("mixer_frontend.proto", fileDescriptor_mixer_frontend_b953796d700de19f)
}

var fileDescriptor_mixer_frontend_b953796d700de19f = []byte{
	// 445 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x52, 0x4d, 0x6f, 0x1a, 0x31,
	0x10, 0xd5, 0xf2, 0x91, 0x96, 0x49, 0x20, 0x91, 0xa1, 0x8a, 0x85, 0x72, 0x40, 0x9c, 0xe8, 0x65,
	0x91, 0xc8, 0x25, 0x6a, 0xa5, 0x4a, 0x55, 0x43, 0xd4, 0x1e, 0x7a, 0xd9, 0xf4, 0x8e, 0x1c, 0xef,
	0x00, 0xab, 0xb2, 0xb6, 0x6b, 0x9b, 0x2a, 0xf9, 0x07, 0xfd, 0x29, 0xfd, 0x99, 0x15, 0x63, 0x2f,
	0x01, 0x11, 0xaa, 0xf6, 0xe6, 0x37, 0xcf, 0x6f, 0xde, 0x68, 0xe6, 0x41, 0xaf, 0x2c, 0x1e, 0xd1,
	0xce, 0xe6, 0x56, 0x2b, 0x8f, 0x2a, 0x4f, 0x8d, 0xd5, 0x5e, 0xb3, 0x26, 0x55, 0xfb, 0x57, 0x0b,
	0xad, 0x17, 0x2b, 0x1c, 0x53, 0xf1, 0x61, 0x3d, 0x1f, 0x3b, 0x6f, 0xd7, 0xd2, 0x87, 0x4f, 0xc3,
	0xdf, 0x0d, 0xe8, 0x7e, 0xdd, 0xfc, 0xbb, 0x8b, 0xe2, 0x4f, 0x5a, 0xcd, 0x8b, 0x05, 0x9b, 0x43,
	0xcf, 0x63, 0x69, 0x56, 0xc2, 0xe3, 0x2c, 0x47, 0x27, 0x6d, 0x61, 0xbc, 0xb6, 0x8e, 0x27, 0x83,
	0xfa, 0xe8, 0x74, 0x72, 0x9d, 0x52, 0xef, 0xf4, 0x05, 0x65, 0xfa, 0x2d, 0xca, 0x6e, 0x9f, 0x55,
	0x53, 0xe5, 0xed, 0x53, 0xd6, 0xf5, 0x87, 0x0c, 0x93, 0xd0, 0x15, 0xb9, 0x30, 0x1e, 0xed, 0x9e,
	0x4d, 0x8d, 0x6c, 0x26, 0x7f, 0xb1, 0xf9, 0x18, 0x54, 0x07, 0x2e, 0x4c, 0x1c, 0x10, 0xec, 0x03,
	0x9c, 0x17, 0xca, 0x79, 0xa1, 0x24, 0xce, 0x24, 0xc9, 0x79, 0x9d, 0x0c, 0xde, 0x44, 0x83, 0x2f,
	0x91, 0x0d, 0xbd, 0xb3, 0x4e, 0xb1, 0x87, 0xd9, 0x7b, 0xe8, 0x2c, 0x85, 0xca, 0x57, 0x68, 0x2b,
	0x79, 0x83, 0xe4, 0xbd, 0x28, 0xff, 0x1c, 0xc8, 0xa8, 0x6e, 0x2f, 0x77, 0x21, 0xbb, 0x81, 0x76,
	0x5e, 0x38, 0x23, 0xbc, 0x5c, 0xce, 0x9c, 0x41, 0xc9, 0x9b, 0xa4, 0xed, 0x46, 0xed, 0x6d, 0xe4,
	0xee, 0x0d, 0xca, 0xec, 0x2c, 0xdf, 0x41, 0xfd, 0x3b, 0xe0, 0xc7, 0x96, 0xc9, 0x2e, 0xa0, 0xfe,
	0x1d, 0x9f, 0x78, 0x32, 0x48, 0x46, 0xad, 0x6c, 0xf3, 0x64, 0x3d, 0x68, 0xfe, 0x14, 0xab, 0x35,
	0xf2, 0x1a, 0xd5, 0x02, 0x78, 0x57, 0xbb, 0x49, 0xfa, 0x53, 0xb8, 0x3c, 0xb2, 0xad, 0xff, 0x69,
	0x33, 0xfc, 0x01, 0x9d, 0xfd, 0x3d, 0x31, 0x06, 0x0d, 0x25, 0x4a, 0x8c, 0x72, 0x7a, 0xb3, 0x3e,
	0xbc, 0xae, 0xee, 0x1c, 0x5b, 0x6c, 0x31, 0x1b, 0xc3, 0x89, 0x11, 0x56, 0x94, 0x8e, 0xd7, 0x07,
	0xc9, 0xe8, 0x74, 0x72, 0x99, 0x86, 0x6c, 0xa6, 0x55, 0x36, 0xd3, 0x7b, 0xca, 0x66, 0x16, 0xbf,
	0x0d, 0x7f, 0x25, 0xd0, 0xde, 0x5b, 0xee, 0x8b, 0x96, 0x1c, 0x5e, 0xc5, 0xa3, 0x47, 0xc7, 0x0a,
	0x06, 0x26, 0xb7, 0xe8, 0x82, 0x23, 0x31, 0x04, 0x77, 0x46, 0x69, 0xfc, 0xdb, 0x28, 0x0e, 0xce,
	0x76, 0x4f, 0xc5, 0xde, 0xc2, 0x45, 0x49, 0x37, 0xc5, 0x47, 0xb3, 0xe9, 0x58, 0x68, 0x15, 0x87,
	0x3a, 0xa7, 0xfa, 0x74, 0x5b, 0x66, 0x57, 0xd0, 0xaa, 0x02, 0x15, 0x92, 0xdd, 0xca, 0x9e, 0x0b,
	0x9b, 0x85, 0xc5, 0xc0, 0x38, 0x4a, 0x65, 0x2b, 0xdb, 0xe2, 0x87, 0x13, 0x9a, 0xe6, 0xfa, 0x4f,
	0x00, 0x00, 0x00, 0xff, 0xff, 0x6e, 0x42, 0x84, 0x37, 0xe1, 0x03, 0x00, 0x00,
}
