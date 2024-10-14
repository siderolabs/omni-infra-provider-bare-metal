// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.24.4
// source: specs/specs.proto

package specs

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type BootMode int32

const (
	BootMode_BOOT_MODE_UNKNOWN    BootMode = 0
	BootMode_BOOT_MODE_AGENT_PXE  BootMode = 1
	BootMode_BOOT_MODE_TALOS_PXE  BootMode = 2
	BootMode_BOOT_MODE_TALOS_DISK BootMode = 3
)

// Enum value maps for BootMode.
var (
	BootMode_name = map[int32]string{
		0: "BOOT_MODE_UNKNOWN",
		1: "BOOT_MODE_AGENT_PXE",
		2: "BOOT_MODE_TALOS_PXE",
		3: "BOOT_MODE_TALOS_DISK",
	}
	BootMode_value = map[string]int32{
		"BOOT_MODE_UNKNOWN":    0,
		"BOOT_MODE_AGENT_PXE":  1,
		"BOOT_MODE_TALOS_PXE":  2,
		"BOOT_MODE_TALOS_DISK": 3,
	}
)

func (x BootMode) Enum() *BootMode {
	p := new(BootMode)
	*p = x
	return p
}

func (x BootMode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (BootMode) Descriptor() protoreflect.EnumDescriptor {
	return file_specs_specs_proto_enumTypes[0].Descriptor()
}

func (BootMode) Type() protoreflect.EnumType {
	return &file_specs_specs_proto_enumTypes[0]
}

func (x BootMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BootMode.Descriptor instead.
func (BootMode) EnumDescriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{0}
}

type PowerState int32

const (
	PowerState_POWER_STATE_UNKNOWN PowerState = 0
	PowerState_POWER_STATE_OFF     PowerState = 1
	PowerState_POWER_STATE_ON      PowerState = 2
)

// Enum value maps for PowerState.
var (
	PowerState_name = map[int32]string{
		0: "POWER_STATE_UNKNOWN",
		1: "POWER_STATE_OFF",
		2: "POWER_STATE_ON",
	}
	PowerState_value = map[string]int32{
		"POWER_STATE_UNKNOWN": 0,
		"POWER_STATE_OFF":     1,
		"POWER_STATE_ON":      2,
	}
)

func (x PowerState) Enum() *PowerState {
	p := new(PowerState)
	*p = x
	return p
}

func (x PowerState) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (PowerState) Descriptor() protoreflect.EnumDescriptor {
	return file_specs_specs_proto_enumTypes[1].Descriptor()
}

func (PowerState) Type() protoreflect.EnumType {
	return &file_specs_specs_proto_enumTypes[1]
}

func (x PowerState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PowerState.Descriptor instead.
func (PowerState) EnumDescriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{1}
}

type PowerManagement struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Ipmi *PowerManagement_IPMI `protobuf:"bytes,1,opt,name=ipmi,proto3" json:"ipmi,omitempty"`
	Api  *PowerManagement_API  `protobuf:"bytes,2,opt,name=api,proto3" json:"api,omitempty"`
}

func (x *PowerManagement) Reset() {
	*x = PowerManagement{}
	mi := &file_specs_specs_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PowerManagement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PowerManagement) ProtoMessage() {}

func (x *PowerManagement) ProtoReflect() protoreflect.Message {
	mi := &file_specs_specs_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PowerManagement.ProtoReflect.Descriptor instead.
func (*PowerManagement) Descriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{0}
}

func (x *PowerManagement) GetIpmi() *PowerManagement_IPMI {
	if x != nil {
		return x.Ipmi
	}
	return nil
}

func (x *PowerManagement) GetApi() *PowerManagement_API {
	if x != nil {
		return x.Api
	}
	return nil
}

type MachineStatusSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PowerManagement *PowerManagement `protobuf:"bytes,1,opt,name=power_management,json=powerManagement,proto3" json:"power_management,omitempty"`
	PowerState      PowerState       `protobuf:"varint,2,opt,name=power_state,json=powerState,proto3,enum=baremetalproviderspecs.PowerState" json:"power_state,omitempty"`
	// LastBootMode is the last observed boot mode of the machine. It is updated by the PXE server each time it boots a server,
	// and is also updated by the status of the agent connectivity.
	BootMode BootMode `protobuf:"varint,3,opt,name=boot_mode,json=bootMode,proto3,enum=baremetalproviderspecs.BootMode" json:"boot_mode,omitempty"`
	// LastWipeId is the ID of the last wipe operation that was performed on the machine.
	//
	// It is used to track if the machine needs to be wiped for an allocation.
	LastWipeId string `protobuf:"bytes,4,opt,name=last_wipe_id,json=lastWipeId,proto3" json:"last_wipe_id,omitempty"`
}

func (x *MachineStatusSpec) Reset() {
	*x = MachineStatusSpec{}
	mi := &file_specs_specs_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MachineStatusSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MachineStatusSpec) ProtoMessage() {}

func (x *MachineStatusSpec) ProtoReflect() protoreflect.Message {
	mi := &file_specs_specs_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MachineStatusSpec.ProtoReflect.Descriptor instead.
func (*MachineStatusSpec) Descriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{1}
}

func (x *MachineStatusSpec) GetPowerManagement() *PowerManagement {
	if x != nil {
		return x.PowerManagement
	}
	return nil
}

func (x *MachineStatusSpec) GetPowerState() PowerState {
	if x != nil {
		return x.PowerState
	}
	return PowerState_POWER_STATE_UNKNOWN
}

func (x *MachineStatusSpec) GetBootMode() BootMode {
	if x != nil {
		return x.BootMode
	}
	return BootMode_BOOT_MODE_UNKNOWN
}

func (x *MachineStatusSpec) GetLastWipeId() string {
	if x != nil {
		return x.LastWipeId
	}
	return ""
}

type PowerManagement_IPMI struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address  string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Port     uint32 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	Username string `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *PowerManagement_IPMI) Reset() {
	*x = PowerManagement_IPMI{}
	mi := &file_specs_specs_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PowerManagement_IPMI) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PowerManagement_IPMI) ProtoMessage() {}

func (x *PowerManagement_IPMI) ProtoReflect() protoreflect.Message {
	mi := &file_specs_specs_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PowerManagement_IPMI.ProtoReflect.Descriptor instead.
func (*PowerManagement_IPMI) Descriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{0, 0}
}

func (x *PowerManagement_IPMI) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *PowerManagement_IPMI) GetPort() uint32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *PowerManagement_IPMI) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *PowerManagement_IPMI) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type PowerManagement_API struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Address string `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
}

func (x *PowerManagement_API) Reset() {
	*x = PowerManagement_API{}
	mi := &file_specs_specs_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PowerManagement_API) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PowerManagement_API) ProtoMessage() {}

func (x *PowerManagement_API) ProtoReflect() protoreflect.Message {
	mi := &file_specs_specs_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PowerManagement_API.ProtoReflect.Descriptor instead.
func (*PowerManagement_API) Descriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{0, 1}
}

func (x *PowerManagement_API) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

var File_specs_specs_proto protoreflect.FileDescriptor

var file_specs_specs_proto_rawDesc = []byte{
	0x0a, 0x11, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2f, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x16, 0x62, 0x61, 0x72, 0x65, 0x6d, 0x65, 0x74, 0x61, 0x6c, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x73, 0x22, 0xa1, 0x02, 0x0a, 0x0f,
	0x50, 0x6f, 0x77, 0x65, 0x72, 0x4d, 0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12,
	0x40, 0x0a, 0x04, 0x69, 0x70, 0x6d, 0x69, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2c, 0x2e,
	0x62, 0x61, 0x72, 0x65, 0x6d, 0x65, 0x74, 0x61, 0x6c, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65,
	0x72, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2e, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x4d, 0x61, 0x6e, 0x61,
	0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x49, 0x50, 0x4d, 0x49, 0x52, 0x04, 0x69, 0x70, 0x6d,
	0x69, 0x12, 0x3d, 0x0a, 0x03, 0x61, 0x70, 0x69, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2b,
	0x2e, 0x62, 0x61, 0x72, 0x65, 0x6d, 0x65, 0x74, 0x61, 0x6c, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2e, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x4d, 0x61, 0x6e,
	0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x41, 0x50, 0x49, 0x52, 0x03, 0x61, 0x70, 0x69,
	0x1a, 0x6c, 0x0a, 0x04, 0x49, 0x50, 0x4d, 0x49, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65,
	0x73, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d,
	0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x1a, 0x1f,
	0x0a, 0x03, 0x41, 0x50, 0x49, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73, 0x22,
	0x8d, 0x02, 0x0a, 0x11, 0x4d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x53, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x53, 0x70, 0x65, 0x63, 0x12, 0x52, 0x0a, 0x10, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x5f, 0x6d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x27, 0x2e, 0x62, 0x61, 0x72, 0x65, 0x6d, 0x65, 0x74, 0x61, 0x6c, 0x70, 0x72, 0x6f, 0x76, 0x69,
	0x64, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2e, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x4d, 0x61,
	0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x52, 0x0f, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x4d,
	0x61, 0x6e, 0x61, 0x67, 0x65, 0x6d, 0x65, 0x6e, 0x74, 0x12, 0x43, 0x0a, 0x0b, 0x70, 0x6f, 0x77,
	0x65, 0x72, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22,
	0x2e, 0x62, 0x61, 0x72, 0x65, 0x6d, 0x65, 0x74, 0x61, 0x6c, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64,
	0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2e, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x52, 0x0a, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x3d,
	0x0a, 0x09, 0x62, 0x6f, 0x6f, 0x74, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x20, 0x2e, 0x62, 0x61, 0x72, 0x65, 0x6d, 0x65, 0x74, 0x61, 0x6c, 0x70, 0x72, 0x6f,
	0x76, 0x69, 0x64, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2e, 0x42, 0x6f, 0x6f, 0x74, 0x4d,
	0x6f, 0x64, 0x65, 0x52, 0x08, 0x62, 0x6f, 0x6f, 0x74, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x20, 0x0a,
	0x0c, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x77, 0x69, 0x70, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x0a, 0x6c, 0x61, 0x73, 0x74, 0x57, 0x69, 0x70, 0x65, 0x49, 0x64, 0x2a,
	0x6d, 0x0a, 0x08, 0x42, 0x6f, 0x6f, 0x74, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x15, 0x0a, 0x11, 0x42,
	0x4f, 0x4f, 0x54, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e,
	0x10, 0x00, 0x12, 0x17, 0x0a, 0x13, 0x42, 0x4f, 0x4f, 0x54, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f,
	0x41, 0x47, 0x45, 0x4e, 0x54, 0x5f, 0x50, 0x58, 0x45, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x42,
	0x4f, 0x4f, 0x54, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x54, 0x41, 0x4c, 0x4f, 0x53, 0x5f, 0x50,
	0x58, 0x45, 0x10, 0x02, 0x12, 0x18, 0x0a, 0x14, 0x42, 0x4f, 0x4f, 0x54, 0x5f, 0x4d, 0x4f, 0x44,
	0x45, 0x5f, 0x54, 0x41, 0x4c, 0x4f, 0x53, 0x5f, 0x44, 0x49, 0x53, 0x4b, 0x10, 0x03, 0x2a, 0x4e,
	0x0a, 0x0a, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x17, 0x0a, 0x13,
	0x50, 0x4f, 0x57, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x55, 0x4e, 0x4b, 0x4e,
	0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x50, 0x4f, 0x57, 0x45, 0x52, 0x5f, 0x53,
	0x54, 0x41, 0x54, 0x45, 0x5f, 0x4f, 0x46, 0x46, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x50, 0x4f,
	0x57, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x4f, 0x4e, 0x10, 0x02, 0x42, 0x40,
	0x5a, 0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x69, 0x64,
	0x65, 0x72, 0x6f, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x6f, 0x6d, 0x6e, 0x69, 0x2d, 0x69, 0x6e, 0x66,
	0x72, 0x61, 0x2d, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2d, 0x62, 0x61, 0x72, 0x65,
	0x2d, 0x6d, 0x65, 0x74, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x70, 0x65, 0x63, 0x73,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_specs_specs_proto_rawDescOnce sync.Once
	file_specs_specs_proto_rawDescData = file_specs_specs_proto_rawDesc
)

func file_specs_specs_proto_rawDescGZIP() []byte {
	file_specs_specs_proto_rawDescOnce.Do(func() {
		file_specs_specs_proto_rawDescData = protoimpl.X.CompressGZIP(file_specs_specs_proto_rawDescData)
	})
	return file_specs_specs_proto_rawDescData
}

var file_specs_specs_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_specs_specs_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_specs_specs_proto_goTypes = []any{
	(BootMode)(0),                // 0: baremetalproviderspecs.BootMode
	(PowerState)(0),              // 1: baremetalproviderspecs.PowerState
	(*PowerManagement)(nil),      // 2: baremetalproviderspecs.PowerManagement
	(*MachineStatusSpec)(nil),    // 3: baremetalproviderspecs.MachineStatusSpec
	(*PowerManagement_IPMI)(nil), // 4: baremetalproviderspecs.PowerManagement.IPMI
	(*PowerManagement_API)(nil),  // 5: baremetalproviderspecs.PowerManagement.API
}
var file_specs_specs_proto_depIdxs = []int32{
	4, // 0: baremetalproviderspecs.PowerManagement.ipmi:type_name -> baremetalproviderspecs.PowerManagement.IPMI
	5, // 1: baremetalproviderspecs.PowerManagement.api:type_name -> baremetalproviderspecs.PowerManagement.API
	2, // 2: baremetalproviderspecs.MachineStatusSpec.power_management:type_name -> baremetalproviderspecs.PowerManagement
	1, // 3: baremetalproviderspecs.MachineStatusSpec.power_state:type_name -> baremetalproviderspecs.PowerState
	0, // 4: baremetalproviderspecs.MachineStatusSpec.boot_mode:type_name -> baremetalproviderspecs.BootMode
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_specs_specs_proto_init() }
func file_specs_specs_proto_init() {
	if File_specs_specs_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_specs_specs_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_specs_specs_proto_goTypes,
		DependencyIndexes: file_specs_specs_proto_depIdxs,
		EnumInfos:         file_specs_specs_proto_enumTypes,
		MessageInfos:      file_specs_specs_proto_msgTypes,
	}.Build()
	File_specs_specs_proto = out.File
	file_specs_specs_proto_rawDesc = nil
	file_specs_specs_proto_goTypes = nil
	file_specs_specs_proto_depIdxs = nil
}