// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.11.2
// source: lbbft.proto

package proto

import (
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type PrePrepareArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	View     uint32     `protobuf:"varint,1,opt,name=View,proto3" json:"View,omitempty"`
	Seq      uint32     `protobuf:"varint,2,opt,name=Seq,proto3" json:"Seq,omitempty"`
	Commands []*Command `protobuf:"bytes,3,rep,name=Commands,proto3" json:"Commands,omitempty"`
}

func (x *PrePrepareArgs) Reset() {
	*x = PrePrepareArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbbft_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrePrepareArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrePrepareArgs) ProtoMessage() {}

func (x *PrePrepareArgs) ProtoReflect() protoreflect.Message {
	mi := &file_lbbft_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrePrepareArgs.ProtoReflect.Descriptor instead.
func (*PrePrepareArgs) Descriptor() ([]byte, []int) {
	return file_lbbft_proto_rawDescGZIP(), []int{0}
}

func (x *PrePrepareArgs) GetView() uint32 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *PrePrepareArgs) GetSeq() uint32 {
	if x != nil {
		return x.Seq
	}
	return 0
}

func (x *PrePrepareArgs) GetCommands() []*Command {
	if x != nil {
		return x.Commands
	}
	return nil
}

type PrePrepareReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sig *PartialSig `protobuf:"bytes,1,opt,name=Sig,proto3" json:"Sig,omitempty"`
}

func (x *PrePrepareReply) Reset() {
	*x = PrePrepareReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbbft_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrePrepareReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrePrepareReply) ProtoMessage() {}

func (x *PrePrepareReply) ProtoReflect() protoreflect.Message {
	mi := &file_lbbft_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrePrepareReply.ProtoReflect.Descriptor instead.
func (*PrePrepareReply) Descriptor() ([]byte, []int) {
	return file_lbbft_proto_rawDescGZIP(), []int{1}
}

func (x *PrePrepareReply) GetSig() *PartialSig {
	if x != nil {
		return x.Sig
	}
	return nil
}

type PrepareArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	View uint32      `protobuf:"varint,1,opt,name=View,proto3" json:"View,omitempty"`
	Seq  uint32      `protobuf:"varint,2,opt,name=Seq,proto3" json:"Seq,omitempty"`
	QC   *QuorumCert `protobuf:"bytes,3,opt,name=QC,proto3" json:"QC,omitempty"`
}

func (x *PrepareArgs) Reset() {
	*x = PrepareArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbbft_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrepareArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareArgs) ProtoMessage() {}

func (x *PrepareArgs) ProtoReflect() protoreflect.Message {
	mi := &file_lbbft_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareArgs.ProtoReflect.Descriptor instead.
func (*PrepareArgs) Descriptor() ([]byte, []int) {
	return file_lbbft_proto_rawDescGZIP(), []int{2}
}

func (x *PrepareArgs) GetView() uint32 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *PrepareArgs) GetSeq() uint32 {
	if x != nil {
		return x.Seq
	}
	return 0
}

func (x *PrepareArgs) GetQC() *QuorumCert {
	if x != nil {
		return x.QC
	}
	return nil
}

type PrepareReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sig *PartialSig `protobuf:"bytes,1,opt,name=Sig,proto3" json:"Sig,omitempty"`
}

func (x *PrepareReply) Reset() {
	*x = PrepareReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbbft_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PrepareReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PrepareReply) ProtoMessage() {}

func (x *PrepareReply) ProtoReflect() protoreflect.Message {
	mi := &file_lbbft_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PrepareReply.ProtoReflect.Descriptor instead.
func (*PrepareReply) Descriptor() ([]byte, []int) {
	return file_lbbft_proto_rawDescGZIP(), []int{3}
}

func (x *PrepareReply) GetSig() *PartialSig {
	if x != nil {
		return x.Sig
	}
	return nil
}

type CommitArgs struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	View uint32      `protobuf:"varint,1,opt,name=View,proto3" json:"View,omitempty"`
	Seq  uint32      `protobuf:"varint,2,opt,name=Seq,proto3" json:"Seq,omitempty"`
	QC   *QuorumCert `protobuf:"bytes,3,opt,name=QC,proto3" json:"QC,omitempty"`
}

func (x *CommitArgs) Reset() {
	*x = CommitArgs{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbbft_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommitArgs) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitArgs) ProtoMessage() {}

func (x *CommitArgs) ProtoReflect() protoreflect.Message {
	mi := &file_lbbft_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitArgs.ProtoReflect.Descriptor instead.
func (*CommitArgs) Descriptor() ([]byte, []int) {
	return file_lbbft_proto_rawDescGZIP(), []int{4}
}

func (x *CommitArgs) GetView() uint32 {
	if x != nil {
		return x.View
	}
	return 0
}

func (x *CommitArgs) GetSeq() uint32 {
	if x != nil {
		return x.Seq
	}
	return 0
}

func (x *CommitArgs) GetQC() *QuorumCert {
	if x != nil {
		return x.QC
	}
	return nil
}

type Command struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Data []byte `protobuf:"bytes,1,opt,name=Data,proto3" json:"Data,omitempty"`
}

func (x *Command) Reset() {
	*x = Command{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbbft_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Command) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Command) ProtoMessage() {}

func (x *Command) ProtoReflect() protoreflect.Message {
	mi := &file_lbbft_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Command.ProtoReflect.Descriptor instead.
func (*Command) Descriptor() ([]byte, []int) {
	return file_lbbft_proto_rawDescGZIP(), []int{5}
}

func (x *Command) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

// ecdsa的签名
type PartialSig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ReplicaID int32  `protobuf:"varint,1,opt,name=ReplicaID,proto3" json:"ReplicaID,omitempty"`
	R         []byte `protobuf:"bytes,2,opt,name=R,proto3" json:"R,omitempty"`
	S         []byte `protobuf:"bytes,3,opt,name=S,proto3" json:"S,omitempty"`
}

func (x *PartialSig) Reset() {
	*x = PartialSig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbbft_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PartialSig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PartialSig) ProtoMessage() {}

func (x *PartialSig) ProtoReflect() protoreflect.Message {
	mi := &file_lbbft_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PartialSig.ProtoReflect.Descriptor instead.
func (*PartialSig) Descriptor() ([]byte, []int) {
	return file_lbbft_proto_rawDescGZIP(), []int{6}
}

func (x *PartialSig) GetReplicaID() int32 {
	if x != nil {
		return x.ReplicaID
	}
	return 0
}

func (x *PartialSig) GetR() []byte {
	if x != nil {
		return x.R
	}
	return nil
}

func (x *PartialSig) GetS() []byte {
	if x != nil {
		return x.S
	}
	return nil
}

type QuorumCert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Sigs       []*PartialSig `protobuf:"bytes,1,rep,name=Sigs,proto3" json:"Sigs,omitempty"`
	SigContent []byte        `protobuf:"bytes,2,opt,name=SigContent,proto3" json:"SigContent,omitempty"`
}

func (x *QuorumCert) Reset() {
	*x = QuorumCert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lbbft_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuorumCert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuorumCert) ProtoMessage() {}

func (x *QuorumCert) ProtoReflect() protoreflect.Message {
	mi := &file_lbbft_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuorumCert.ProtoReflect.Descriptor instead.
func (*QuorumCert) Descriptor() ([]byte, []int) {
	return file_lbbft_proto_rawDescGZIP(), []int{7}
}

func (x *QuorumCert) GetSigs() []*PartialSig {
	if x != nil {
		return x.Sigs
	}
	return nil
}

func (x *QuorumCert) GetSigContent() []byte {
	if x != nil {
		return x.SigContent
	}
	return nil
}

var File_lbbft_proto protoreflect.FileDescriptor

var file_lbbft_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6c, 0x62, 0x62, 0x66, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x62, 0x0a, 0x0e, 0x50, 0x72, 0x65, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x41,
	0x72, 0x67, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x56, 0x69, 0x65, 0x77, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x04, 0x56, 0x69, 0x65, 0x77, 0x12, 0x10, 0x0a, 0x03, 0x53, 0x65, 0x71, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x53, 0x65, 0x71, 0x12, 0x2a, 0x0a, 0x08, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x08, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x73, 0x22, 0x36, 0x0a, 0x0f, 0x50, 0x72, 0x65, 0x50, 0x72, 0x65, 0x70,
	0x61, 0x72, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x23, 0x0a, 0x03, 0x53, 0x69, 0x67, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61,
	0x72, 0x74, 0x69, 0x61, 0x6c, 0x53, 0x69, 0x67, 0x52, 0x03, 0x53, 0x69, 0x67, 0x22, 0x56, 0x0a,
	0x0b, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x41, 0x72, 0x67, 0x73, 0x12, 0x12, 0x0a, 0x04,
	0x56, 0x69, 0x65, 0x77, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x56, 0x69, 0x65, 0x77,
	0x12, 0x10, 0x0a, 0x03, 0x53, 0x65, 0x71, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x53,
	0x65, 0x71, 0x12, 0x21, 0x0a, 0x02, 0x51, 0x43, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72,
	0x74, 0x52, 0x02, 0x51, 0x43, 0x22, 0x33, 0x0a, 0x0c, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65,
	0x52, 0x65, 0x70, 0x6c, 0x79, 0x12, 0x23, 0x0a, 0x03, 0x53, 0x69, 0x67, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x72, 0x74, 0x69,
	0x61, 0x6c, 0x53, 0x69, 0x67, 0x52, 0x03, 0x53, 0x69, 0x67, 0x22, 0x55, 0x0a, 0x0a, 0x43, 0x6f,
	0x6d, 0x6d, 0x69, 0x74, 0x41, 0x72, 0x67, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x56, 0x69, 0x65, 0x77,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x56, 0x69, 0x65, 0x77, 0x12, 0x10, 0x0a, 0x03,
	0x53, 0x65, 0x71, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x03, 0x53, 0x65, 0x71, 0x12, 0x21,
	0x0a, 0x02, 0x51, 0x43, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x2e, 0x51, 0x75, 0x6f, 0x72, 0x75, 0x6d, 0x43, 0x65, 0x72, 0x74, 0x52, 0x02, 0x51,
	0x43, 0x22, 0x1d, 0x0a, 0x07, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x44, 0x61, 0x74, 0x61, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x44, 0x61, 0x74, 0x61,
	0x22, 0x46, 0x0a, 0x0a, 0x50, 0x61, 0x72, 0x74, 0x69, 0x61, 0x6c, 0x53, 0x69, 0x67, 0x12, 0x1c,
	0x0a, 0x09, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x09, 0x52, 0x65, 0x70, 0x6c, 0x69, 0x63, 0x61, 0x49, 0x44, 0x12, 0x0c, 0x0a, 0x01,
	0x52, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x01, 0x52, 0x12, 0x0c, 0x0a, 0x01, 0x53, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x01, 0x53, 0x22, 0x53, 0x0a, 0x0a, 0x51, 0x75, 0x6f, 0x72,
	0x75, 0x6d, 0x43, 0x65, 0x72, 0x74, 0x12, 0x25, 0x0a, 0x04, 0x53, 0x69, 0x67, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x61, 0x72,
	0x74, 0x69, 0x61, 0x6c, 0x53, 0x69, 0x67, 0x52, 0x04, 0x53, 0x69, 0x67, 0x73, 0x12, 0x1e, 0x0a,
	0x0a, 0x53, 0x69, 0x67, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x0a, 0x53, 0x69, 0x67, 0x43, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x32, 0xb3, 0x01,
	0x0a, 0x05, 0x4c, 0x42, 0x42, 0x46, 0x54, 0x12, 0x3d, 0x0a, 0x0a, 0x50, 0x72, 0x65, 0x50, 0x72,
	0x65, 0x70, 0x61, 0x72, 0x65, 0x12, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72,
	0x65, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x41, 0x72, 0x67, 0x73, 0x1a, 0x16, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x65, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72, 0x65, 0x52,
	0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x34, 0x0a, 0x07, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72,
	0x65, 0x12, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72, 0x65, 0x70, 0x61, 0x72,
	0x65, 0x41, 0x72, 0x67, 0x73, 0x1a, 0x13, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x50, 0x72,
	0x65, 0x70, 0x61, 0x72, 0x65, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x00, 0x12, 0x35, 0x0a, 0x06,
	0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43,
	0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x41, 0x72, 0x67, 0x73, 0x1a, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0x00, 0x42, 0x29, 0x5a, 0x27, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x6a, 0x6f, 0x65, 0x2d, 0x7a, 0x78, 0x68, 0x2f, 0x6c, 0x62, 0x62, 0x66, 0x74, 0x2f,
	0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_lbbft_proto_rawDescOnce sync.Once
	file_lbbft_proto_rawDescData = file_lbbft_proto_rawDesc
)

func file_lbbft_proto_rawDescGZIP() []byte {
	file_lbbft_proto_rawDescOnce.Do(func() {
		file_lbbft_proto_rawDescData = protoimpl.X.CompressGZIP(file_lbbft_proto_rawDescData)
	})
	return file_lbbft_proto_rawDescData
}

var file_lbbft_proto_msgTypes = make([]protoimpl.MessageInfo, 8)
var file_lbbft_proto_goTypes = []interface{}{
	(*PrePrepareArgs)(nil),  // 0: proto.PrePrepareArgs
	(*PrePrepareReply)(nil), // 1: proto.PrePrepareReply
	(*PrepareArgs)(nil),     // 2: proto.PrepareArgs
	(*PrepareReply)(nil),    // 3: proto.PrepareReply
	(*CommitArgs)(nil),      // 4: proto.CommitArgs
	(*Command)(nil),         // 5: proto.Command
	(*PartialSig)(nil),      // 6: proto.PartialSig
	(*QuorumCert)(nil),      // 7: proto.QuorumCert
	(*empty.Empty)(nil),     // 8: google.protobuf.Empty
}
var file_lbbft_proto_depIdxs = []int32{
	5, // 0: proto.PrePrepareArgs.Commands:type_name -> proto.Command
	6, // 1: proto.PrePrepareReply.Sig:type_name -> proto.PartialSig
	7, // 2: proto.PrepareArgs.QC:type_name -> proto.QuorumCert
	6, // 3: proto.PrepareReply.Sig:type_name -> proto.PartialSig
	7, // 4: proto.CommitArgs.QC:type_name -> proto.QuorumCert
	6, // 5: proto.QuorumCert.Sigs:type_name -> proto.PartialSig
	0, // 6: proto.LBBFT.PrePrepare:input_type -> proto.PrePrepareArgs
	2, // 7: proto.LBBFT.Prepare:input_type -> proto.PrepareArgs
	4, // 8: proto.LBBFT.Commit:input_type -> proto.CommitArgs
	1, // 9: proto.LBBFT.PrePrepare:output_type -> proto.PrePrepareReply
	3, // 10: proto.LBBFT.Prepare:output_type -> proto.PrepareReply
	8, // 11: proto.LBBFT.Commit:output_type -> google.protobuf.Empty
	9, // [9:12] is the sub-list for method output_type
	6, // [6:9] is the sub-list for method input_type
	6, // [6:6] is the sub-list for extension type_name
	6, // [6:6] is the sub-list for extension extendee
	0, // [0:6] is the sub-list for field type_name
}

func init() { file_lbbft_proto_init() }
func file_lbbft_proto_init() {
	if File_lbbft_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_lbbft_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrePrepareArgs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbbft_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrePrepareReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbbft_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrepareArgs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbbft_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PrepareReply); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbbft_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommitArgs); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbbft_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Command); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbbft_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PartialSig); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_lbbft_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuorumCert); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_lbbft_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   8,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_lbbft_proto_goTypes,
		DependencyIndexes: file_lbbft_proto_depIdxs,
		MessageInfos:      file_lbbft_proto_msgTypes,
	}.Build()
	File_lbbft_proto = out.File
	file_lbbft_proto_rawDesc = nil
	file_lbbft_proto_goTypes = nil
	file_lbbft_proto_depIdxs = nil
}
