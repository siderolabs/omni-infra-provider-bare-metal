// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.2
// 	protoc        v4.24.4
// source: specs/specs.proto

package specs

import (
	reflect "reflect"
	sync "sync"

	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

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
	return file_specs_specs_proto_enumTypes[0].Descriptor()
}

func (PowerState) Type() protoreflect.EnumType {
	return &file_specs_specs_proto_enumTypes[0]
}

func (x PowerState) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use PowerState.Descriptor instead.
func (PowerState) EnumDescriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{0}
}

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
	return file_specs_specs_proto_enumTypes[1].Descriptor()
}

func (BootMode) Type() protoreflect.EnumType {
	return &file_specs_specs_proto_enumTypes[1]
}

func (x BootMode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use BootMode.Descriptor instead.
func (BootMode) EnumDescriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{1}
}

type PowerOperationSpec struct {
	state                protoimpl.MessageState `protogen:"open.v1"`
	LastPowerOperation   PowerState             `protobuf:"varint,1,opt,name=last_power_operation,json=lastPowerOperation,proto3,enum=baremetalproviderspecs.PowerState" json:"last_power_operation,omitempty"`
	LastPowerOnTimestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=last_power_on_timestamp,json=lastPowerOnTimestamp,proto3" json:"last_power_on_timestamp,omitempty"`
	unknownFields        protoimpl.UnknownFields
	sizeCache            protoimpl.SizeCache
}

func (x *PowerOperationSpec) Reset() {
	*x = PowerOperationSpec{}
	mi := &file_specs_specs_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *PowerOperationSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PowerOperationSpec) ProtoMessage() {}

func (x *PowerOperationSpec) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use PowerOperationSpec.ProtoReflect.Descriptor instead.
func (*PowerOperationSpec) Descriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{0}
}

func (x *PowerOperationSpec) GetLastPowerOperation() PowerState {
	if x != nil {
		return x.LastPowerOperation
	}
	return PowerState_POWER_STATE_UNKNOWN
}

func (x *PowerOperationSpec) GetLastPowerOnTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.LastPowerOnTimestamp
	}
	return nil
}

type BMCConfigurationSpec struct {
	state              protoimpl.MessageState     `protogen:"open.v1"`
	Ipmi               *BMCConfigurationSpec_IPMI `protobuf:"bytes,1,opt,name=ipmi,proto3" json:"ipmi,omitempty"`
	Api                *BMCConfigurationSpec_API  `protobuf:"bytes,2,opt,name=api,proto3" json:"api,omitempty"`
	ManuallyConfigured bool                       `protobuf:"varint,3,opt,name=manually_configured,json=manuallyConfigured,proto3" json:"manually_configured,omitempty"`
	unknownFields      protoimpl.UnknownFields
	sizeCache          protoimpl.SizeCache
}

func (x *BMCConfigurationSpec) Reset() {
	*x = BMCConfigurationSpec{}
	mi := &file_specs_specs_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BMCConfigurationSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BMCConfigurationSpec) ProtoMessage() {}

func (x *BMCConfigurationSpec) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use BMCConfigurationSpec.ProtoReflect.Descriptor instead.
func (*BMCConfigurationSpec) Descriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{1}
}

func (x *BMCConfigurationSpec) GetIpmi() *BMCConfigurationSpec_IPMI {
	if x != nil {
		return x.Ipmi
	}
	return nil
}

func (x *BMCConfigurationSpec) GetApi() *BMCConfigurationSpec_API {
	if x != nil {
		return x.Api
	}
	return nil
}

func (x *BMCConfigurationSpec) GetManuallyConfigured() bool {
	if x != nil {
		return x.ManuallyConfigured
	}
	return false
}

type MachineStatusSpec struct {
	state           protoimpl.MessageState `protogen:"open.v1"`
	AgentAccessible bool                   `protobuf:"varint,1,opt,name=agent_accessible,json=agentAccessible,proto3" json:"agent_accessible,omitempty"`
	PowerState      PowerState             `protobuf:"varint,2,opt,name=power_state,json=powerState,proto3,enum=baremetalproviderspecs.PowerState" json:"power_state,omitempty"`
	LastPxeBootMode BootMode               `protobuf:"varint,3,opt,name=last_pxe_boot_mode,json=lastPxeBootMode,proto3,enum=baremetalproviderspecs.BootMode" json:"last_pxe_boot_mode,omitempty"`
	unknownFields   protoimpl.UnknownFields
	sizeCache       protoimpl.SizeCache
}

func (x *MachineStatusSpec) Reset() {
	*x = MachineStatusSpec{}
	mi := &file_specs_specs_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MachineStatusSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MachineStatusSpec) ProtoMessage() {}

func (x *MachineStatusSpec) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use MachineStatusSpec.ProtoReflect.Descriptor instead.
func (*MachineStatusSpec) Descriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{2}
}

func (x *MachineStatusSpec) GetAgentAccessible() bool {
	if x != nil {
		return x.AgentAccessible
	}
	return false
}

func (x *MachineStatusSpec) GetPowerState() PowerState {
	if x != nil {
		return x.PowerState
	}
	return PowerState_POWER_STATE_UNKNOWN
}

func (x *MachineStatusSpec) GetLastPxeBootMode() BootMode {
	if x != nil {
		return x.LastPxeBootMode
	}
	return BootMode_BOOT_MODE_UNKNOWN
}

type WipeStatusSpec struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// LastWipeId is the ID of the last wipe operation that was performed on the machine.
	//
	// It is used to track if the machine needs to be wiped for an allocation.
	LastWipeId string `protobuf:"bytes,1,opt,name=last_wipe_id,json=lastWipeId,proto3" json:"last_wipe_id,omitempty"`
	// LastWipeInstallEventId is set to the same value of InfraMachine.InstallEventId field each time machine gets wiped.
	//
	// Using this, the provider is able to track the installation state of Talos on the machine. It does it by comparing this stored value
	// with the value of InfraMachine.InstallEventId field.
	//
	// If the value of InfraMachine.InstallEventId field is greater than the value of this field,
	// it means that Omni observed, after the wipe, at least one event indicating Talos is installed on that machine.
	LastWipeInstallEventId uint64 `protobuf:"varint,2,opt,name=last_wipe_install_event_id,json=lastWipeInstallEventId,proto3" json:"last_wipe_install_event_id,omitempty"`
	InitialWipeDone        bool   `protobuf:"varint,3,opt,name=initial_wipe_done,json=initialWipeDone,proto3" json:"initial_wipe_done,omitempty"`
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *WipeStatusSpec) Reset() {
	*x = WipeStatusSpec{}
	mi := &file_specs_specs_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *WipeStatusSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*WipeStatusSpec) ProtoMessage() {}

func (x *WipeStatusSpec) ProtoReflect() protoreflect.Message {
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

// Deprecated: Use WipeStatusSpec.ProtoReflect.Descriptor instead.
func (*WipeStatusSpec) Descriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{3}
}

func (x *WipeStatusSpec) GetLastWipeId() string {
	if x != nil {
		return x.LastWipeId
	}
	return ""
}

func (x *WipeStatusSpec) GetLastWipeInstallEventId() uint64 {
	if x != nil {
		return x.LastWipeInstallEventId
	}
	return 0
}

func (x *WipeStatusSpec) GetInitialWipeDone() bool {
	if x != nil {
		return x.InitialWipeDone
	}
	return false
}

type RebootStatusSpec struct {
	state        protoimpl.MessageState `protogen:"open.v1"`
	LastRebootId string                 `protobuf:"bytes,1,opt,name=last_reboot_id,json=lastRebootId,proto3" json:"last_reboot_id,omitempty"`
	// LastRebootTimestamp is the timestamp of the last reboot (or power on) of the machine.
	//
	// It is used to track the last reboot time of the machine, and to enforce the MinRebootInterval.
	LastRebootTimestamp *timestamppb.Timestamp `protobuf:"bytes,2,opt,name=last_reboot_timestamp,json=lastRebootTimestamp,proto3" json:"last_reboot_timestamp,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *RebootStatusSpec) Reset() {
	*x = RebootStatusSpec{}
	mi := &file_specs_specs_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *RebootStatusSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RebootStatusSpec) ProtoMessage() {}

func (x *RebootStatusSpec) ProtoReflect() protoreflect.Message {
	mi := &file_specs_specs_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RebootStatusSpec.ProtoReflect.Descriptor instead.
func (*RebootStatusSpec) Descriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{4}
}

func (x *RebootStatusSpec) GetLastRebootId() string {
	if x != nil {
		return x.LastRebootId
	}
	return ""
}

func (x *RebootStatusSpec) GetLastRebootTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.LastRebootTimestamp
	}
	return nil
}

type TLSConfigSpec struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	CaCert        string                 `protobuf:"bytes,1,opt,name=ca_cert,json=caCert,proto3" json:"ca_cert,omitempty"`
	CaKey         string                 `protobuf:"bytes,2,opt,name=ca_key,json=caKey,proto3" json:"ca_key,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *TLSConfigSpec) Reset() {
	*x = TLSConfigSpec{}
	mi := &file_specs_specs_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TLSConfigSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TLSConfigSpec) ProtoMessage() {}

func (x *TLSConfigSpec) ProtoReflect() protoreflect.Message {
	mi := &file_specs_specs_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TLSConfigSpec.ProtoReflect.Descriptor instead.
func (*TLSConfigSpec) Descriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{5}
}

func (x *TLSConfigSpec) GetCaCert() string {
	if x != nil {
		return x.CaCert
	}
	return ""
}

func (x *TLSConfigSpec) GetCaKey() string {
	if x != nil {
		return x.CaKey
	}
	return ""
}

type BMCConfigurationSpec_IPMI struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Address       string                 `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	Port          uint32                 `protobuf:"varint,2,opt,name=port,proto3" json:"port,omitempty"`
	Username      string                 `protobuf:"bytes,3,opt,name=username,proto3" json:"username,omitempty"`
	Password      string                 `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BMCConfigurationSpec_IPMI) Reset() {
	*x = BMCConfigurationSpec_IPMI{}
	mi := &file_specs_specs_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BMCConfigurationSpec_IPMI) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BMCConfigurationSpec_IPMI) ProtoMessage() {}

func (x *BMCConfigurationSpec_IPMI) ProtoReflect() protoreflect.Message {
	mi := &file_specs_specs_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BMCConfigurationSpec_IPMI.ProtoReflect.Descriptor instead.
func (*BMCConfigurationSpec_IPMI) Descriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{1, 0}
}

func (x *BMCConfigurationSpec_IPMI) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

func (x *BMCConfigurationSpec_IPMI) GetPort() uint32 {
	if x != nil {
		return x.Port
	}
	return 0
}

func (x *BMCConfigurationSpec_IPMI) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *BMCConfigurationSpec_IPMI) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

type BMCConfigurationSpec_API struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Address       string                 `protobuf:"bytes,1,opt,name=address,proto3" json:"address,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *BMCConfigurationSpec_API) Reset() {
	*x = BMCConfigurationSpec_API{}
	mi := &file_specs_specs_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *BMCConfigurationSpec_API) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BMCConfigurationSpec_API) ProtoMessage() {}

func (x *BMCConfigurationSpec_API) ProtoReflect() protoreflect.Message {
	mi := &file_specs_specs_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BMCConfigurationSpec_API.ProtoReflect.Descriptor instead.
func (*BMCConfigurationSpec_API) Descriptor() ([]byte, []int) {
	return file_specs_specs_proto_rawDescGZIP(), []int{1, 1}
}

func (x *BMCConfigurationSpec_API) GetAddress() string {
	if x != nil {
		return x.Address
	}
	return ""
}

var File_specs_specs_proto protoreflect.FileDescriptor

var file_specs_specs_proto_rawDesc = []byte{
	0x0a, 0x11, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2f, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x16, 0x62, 0x61, 0x72, 0x65, 0x6d, 0x65, 0x74, 0x61, 0x6c, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x73, 0x1a, 0x1f, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xbd, 0x01, 0x0a,
	0x12, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53,
	0x70, 0x65, 0x63, 0x12, 0x54, 0x0a, 0x14, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x70, 0x6f, 0x77, 0x65,
	0x72, 0x5f, 0x6f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x22, 0x2e, 0x62, 0x61, 0x72, 0x65, 0x6d, 0x65, 0x74, 0x61, 0x6c, 0x70, 0x72, 0x6f,
	0x76, 0x69, 0x64, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2e, 0x50, 0x6f, 0x77, 0x65, 0x72,
	0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x12, 0x6c, 0x61, 0x73, 0x74, 0x50, 0x6f, 0x77, 0x65, 0x72,
	0x4f, 0x70, 0x65, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x51, 0x0a, 0x17, 0x6c, 0x61, 0x73,
	0x74, 0x5f, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x5f, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x14, 0x6c, 0x61, 0x73, 0x74, 0x50, 0x6f, 0x77, 0x65,
	0x72, 0x4f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0xe1, 0x02, 0x0a,
	0x14, 0x42, 0x4d, 0x43, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f,
	0x6e, 0x53, 0x70, 0x65, 0x63, 0x12, 0x45, 0x0a, 0x04, 0x69, 0x70, 0x6d, 0x69, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x31, 0x2e, 0x62, 0x61, 0x72, 0x65, 0x6d, 0x65, 0x74, 0x61, 0x6c, 0x70,
	0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2e, 0x42, 0x4d, 0x43,
	0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x53, 0x70, 0x65,
	0x63, 0x2e, 0x49, 0x50, 0x4d, 0x49, 0x52, 0x04, 0x69, 0x70, 0x6d, 0x69, 0x12, 0x42, 0x0a, 0x03,
	0x61, 0x70, 0x69, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x30, 0x2e, 0x62, 0x61, 0x72, 0x65,
	0x6d, 0x65, 0x74, 0x61, 0x6c, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x73, 0x70, 0x65,
	0x63, 0x73, 0x2e, 0x42, 0x4d, 0x43, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x53, 0x70, 0x65, 0x63, 0x2e, 0x41, 0x50, 0x49, 0x52, 0x03, 0x61, 0x70, 0x69,
	0x12, 0x2f, 0x0a, 0x13, 0x6d, 0x61, 0x6e, 0x75, 0x61, 0x6c, 0x6c, 0x79, 0x5f, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x75, 0x72, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x08, 0x52, 0x12, 0x6d,
	0x61, 0x6e, 0x75, 0x61, 0x6c, 0x6c, 0x79, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x65,
	0x64, 0x1a, 0x6c, 0x0a, 0x04, 0x49, 0x50, 0x4d, 0x49, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64,
	0x72, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72,
	0x65, 0x73, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x04, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x1a,
	0x1f, 0x0a, 0x03, 0x41, 0x50, 0x49, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x61, 0x64, 0x64, 0x72, 0x65, 0x73, 0x73,
	0x22, 0xd2, 0x01, 0x0a, 0x11, 0x4d, 0x61, 0x63, 0x68, 0x69, 0x6e, 0x65, 0x53, 0x74, 0x61, 0x74,
	0x75, 0x73, 0x53, 0x70, 0x65, 0x63, 0x12, 0x29, 0x0a, 0x10, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x5f,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x69, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08,
	0x52, 0x0f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x41, 0x63, 0x63, 0x65, 0x73, 0x73, 0x69, 0x62, 0x6c,
	0x65, 0x12, 0x43, 0x0a, 0x0b, 0x70, 0x6f, 0x77, 0x65, 0x72, 0x5f, 0x73, 0x74, 0x61, 0x74, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x22, 0x2e, 0x62, 0x61, 0x72, 0x65, 0x6d, 0x65, 0x74,
	0x61, 0x6c, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2e,
	0x50, 0x6f, 0x77, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x52, 0x0a, 0x70, 0x6f, 0x77, 0x65,
	0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x4d, 0x0a, 0x12, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x70,
	0x78, 0x65, 0x5f, 0x62, 0x6f, 0x6f, 0x74, 0x5f, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x20, 0x2e, 0x62, 0x61, 0x72, 0x65, 0x6d, 0x65, 0x74, 0x61, 0x6c, 0x70, 0x72,
	0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x73, 0x70, 0x65, 0x63, 0x73, 0x2e, 0x42, 0x6f, 0x6f, 0x74,
	0x4d, 0x6f, 0x64, 0x65, 0x52, 0x0f, 0x6c, 0x61, 0x73, 0x74, 0x50, 0x78, 0x65, 0x42, 0x6f, 0x6f,
	0x74, 0x4d, 0x6f, 0x64, 0x65, 0x22, 0x9a, 0x01, 0x0a, 0x0e, 0x57, 0x69, 0x70, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x53, 0x70, 0x65, 0x63, 0x12, 0x20, 0x0a, 0x0c, 0x6c, 0x61, 0x73, 0x74,
	0x5f, 0x77, 0x69, 0x70, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a,
	0x6c, 0x61, 0x73, 0x74, 0x57, 0x69, 0x70, 0x65, 0x49, 0x64, 0x12, 0x3a, 0x0a, 0x1a, 0x6c, 0x61,
	0x73, 0x74, 0x5f, 0x77, 0x69, 0x70, 0x65, 0x5f, 0x69, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x5f,
	0x65, 0x76, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x16,
	0x6c, 0x61, 0x73, 0x74, 0x57, 0x69, 0x70, 0x65, 0x49, 0x6e, 0x73, 0x74, 0x61, 0x6c, 0x6c, 0x45,
	0x76, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x2a, 0x0a, 0x11, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61,
	0x6c, 0x5f, 0x77, 0x69, 0x70, 0x65, 0x5f, 0x64, 0x6f, 0x6e, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x0f, 0x69, 0x6e, 0x69, 0x74, 0x69, 0x61, 0x6c, 0x57, 0x69, 0x70, 0x65, 0x44, 0x6f,
	0x6e, 0x65, 0x22, 0x88, 0x01, 0x0a, 0x10, 0x52, 0x65, 0x62, 0x6f, 0x6f, 0x74, 0x53, 0x74, 0x61,
	0x74, 0x75, 0x73, 0x53, 0x70, 0x65, 0x63, 0x12, 0x24, 0x0a, 0x0e, 0x6c, 0x61, 0x73, 0x74, 0x5f,
	0x72, 0x65, 0x62, 0x6f, 0x6f, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0c, 0x6c, 0x61, 0x73, 0x74, 0x52, 0x65, 0x62, 0x6f, 0x6f, 0x74, 0x49, 0x64, 0x12, 0x4e, 0x0a,
	0x15, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x62, 0x6f, 0x6f, 0x74, 0x5f, 0x74, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x13, 0x6c, 0x61, 0x73, 0x74, 0x52, 0x65,
	0x62, 0x6f, 0x6f, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x3f, 0x0a,
	0x0d, 0x54, 0x4c, 0x53, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x53, 0x70, 0x65, 0x63, 0x12, 0x17,
	0x0a, 0x07, 0x63, 0x61, 0x5f, 0x63, 0x65, 0x72, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x63, 0x61, 0x43, 0x65, 0x72, 0x74, 0x12, 0x15, 0x0a, 0x06, 0x63, 0x61, 0x5f, 0x6b, 0x65,
	0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x63, 0x61, 0x4b, 0x65, 0x79, 0x2a, 0x4e,
	0x0a, 0x0a, 0x50, 0x6f, 0x77, 0x65, 0x72, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12, 0x17, 0x0a, 0x13,
	0x50, 0x4f, 0x57, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x55, 0x4e, 0x4b, 0x4e,
	0x4f, 0x57, 0x4e, 0x10, 0x00, 0x12, 0x13, 0x0a, 0x0f, 0x50, 0x4f, 0x57, 0x45, 0x52, 0x5f, 0x53,
	0x54, 0x41, 0x54, 0x45, 0x5f, 0x4f, 0x46, 0x46, 0x10, 0x01, 0x12, 0x12, 0x0a, 0x0e, 0x50, 0x4f,
	0x57, 0x45, 0x52, 0x5f, 0x53, 0x54, 0x41, 0x54, 0x45, 0x5f, 0x4f, 0x4e, 0x10, 0x02, 0x2a, 0x6d,
	0x0a, 0x08, 0x42, 0x6f, 0x6f, 0x74, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x15, 0x0a, 0x11, 0x42, 0x4f,
	0x4f, 0x54, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x10,
	0x00, 0x12, 0x17, 0x0a, 0x13, 0x42, 0x4f, 0x4f, 0x54, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x41,
	0x47, 0x45, 0x4e, 0x54, 0x5f, 0x50, 0x58, 0x45, 0x10, 0x01, 0x12, 0x17, 0x0a, 0x13, 0x42, 0x4f,
	0x4f, 0x54, 0x5f, 0x4d, 0x4f, 0x44, 0x45, 0x5f, 0x54, 0x41, 0x4c, 0x4f, 0x53, 0x5f, 0x50, 0x58,
	0x45, 0x10, 0x02, 0x12, 0x18, 0x0a, 0x14, 0x42, 0x4f, 0x4f, 0x54, 0x5f, 0x4d, 0x4f, 0x44, 0x45,
	0x5f, 0x54, 0x41, 0x4c, 0x4f, 0x53, 0x5f, 0x44, 0x49, 0x53, 0x4b, 0x10, 0x03, 0x42, 0x40, 0x5a,
	0x3e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x73, 0x69, 0x64, 0x65,
	0x72, 0x6f, 0x6c, 0x61, 0x62, 0x73, 0x2f, 0x6f, 0x6d, 0x6e, 0x69, 0x2d, 0x69, 0x6e, 0x66, 0x72,
	0x61, 0x2d, 0x70, 0x72, 0x6f, 0x76, 0x69, 0x64, 0x65, 0x72, 0x2d, 0x62, 0x61, 0x72, 0x65, 0x2d,
	0x6d, 0x65, 0x74, 0x61, 0x6c, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x73, 0x70, 0x65, 0x63, 0x73, 0x62,
	0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
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
var file_specs_specs_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_specs_specs_proto_goTypes = []any{
	(PowerState)(0),                   // 0: baremetalproviderspecs.PowerState
	(BootMode)(0),                     // 1: baremetalproviderspecs.BootMode
	(*PowerOperationSpec)(nil),        // 2: baremetalproviderspecs.PowerOperationSpec
	(*BMCConfigurationSpec)(nil),      // 3: baremetalproviderspecs.BMCConfigurationSpec
	(*MachineStatusSpec)(nil),         // 4: baremetalproviderspecs.MachineStatusSpec
	(*WipeStatusSpec)(nil),            // 5: baremetalproviderspecs.WipeStatusSpec
	(*RebootStatusSpec)(nil),          // 6: baremetalproviderspecs.RebootStatusSpec
	(*TLSConfigSpec)(nil),             // 7: baremetalproviderspecs.TLSConfigSpec
	(*BMCConfigurationSpec_IPMI)(nil), // 8: baremetalproviderspecs.BMCConfigurationSpec.IPMI
	(*BMCConfigurationSpec_API)(nil),  // 9: baremetalproviderspecs.BMCConfigurationSpec.API
	(*timestamppb.Timestamp)(nil),     // 10: google.protobuf.Timestamp
}
var file_specs_specs_proto_depIdxs = []int32{
	0,  // 0: baremetalproviderspecs.PowerOperationSpec.last_power_operation:type_name -> baremetalproviderspecs.PowerState
	10, // 1: baremetalproviderspecs.PowerOperationSpec.last_power_on_timestamp:type_name -> google.protobuf.Timestamp
	8,  // 2: baremetalproviderspecs.BMCConfigurationSpec.ipmi:type_name -> baremetalproviderspecs.BMCConfigurationSpec.IPMI
	9,  // 3: baremetalproviderspecs.BMCConfigurationSpec.api:type_name -> baremetalproviderspecs.BMCConfigurationSpec.API
	0,  // 4: baremetalproviderspecs.MachineStatusSpec.power_state:type_name -> baremetalproviderspecs.PowerState
	1,  // 5: baremetalproviderspecs.MachineStatusSpec.last_pxe_boot_mode:type_name -> baremetalproviderspecs.BootMode
	10, // 6: baremetalproviderspecs.RebootStatusSpec.last_reboot_timestamp:type_name -> google.protobuf.Timestamp
	7,  // [7:7] is the sub-list for method output_type
	7,  // [7:7] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
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
			NumMessages:   8,
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
