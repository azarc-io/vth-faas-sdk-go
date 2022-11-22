// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1-devel
// 	protoc        (unknown)
// source: azarc/common/v1/struct.proto

package commonv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type NullValue int32

const (
	// Null value.
	NullValue_NULL_VALUE NullValue = 0
)

// Enum value maps for NullValue.
var (
	NullValue_name = map[int32]string{
		0: "NULL_VALUE",
	}
	NullValue_value = map[string]int32{
		"NULL_VALUE": 0,
	}
)

func (x NullValue) Enum() *NullValue {
	p := new(NullValue)
	*p = x
	return p
}

func (x NullValue) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (NullValue) Descriptor() protoreflect.EnumDescriptor {
	return file_azarc_common_v1_struct_proto_enumTypes[0].Descriptor()
}

func (NullValue) Type() protoreflect.EnumType {
	return &file_azarc_common_v1_struct_proto_enumTypes[0]
}

func (x NullValue) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use NullValue.Descriptor instead.
func (NullValue) EnumDescriptor() ([]byte, []int) {
	return file_azarc_common_v1_struct_proto_rawDescGZIP(), []int{0}
}

type Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// The kind of value.
	//
	// Types that are assignable to Kind:
	//
	//	*Value_NullValue
	//	*Value_DoubleValue
	//	*Value_FloatValue
	//	*Value_Int32Value
	//	*Value_Int64Value
	//	*Value_Uint32Value
	//	*Value_Uint64Value
	//	*Value_Sint32Value
	//	*Value_Sint64Value
	//	*Value_Fixed32Value
	//	*Value_Fixed64Value
	//	*Value_Sfixed32Value
	//	*Value_Sfixed64Value
	//	*Value_BoolValue
	//	*Value_StringValue
	//	*Value_BytesValue
	//	*Value_AnyValue
	//	*Value_StructValue
	//	*Value_ListValue
	Kind isValue_Kind `protobuf_oneof:"kind"`
}

func (x *Value) Reset() {
	*x = Value{}
	if protoimpl.UnsafeEnabled {
		mi := &file_azarc_common_v1_struct_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_azarc_common_v1_struct_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_azarc_common_v1_struct_proto_rawDescGZIP(), []int{0}
}

func (m *Value) GetKind() isValue_Kind {
	if m != nil {
		return m.Kind
	}
	return nil
}

func (x *Value) GetNullValue() NullValue {
	if x, ok := x.GetKind().(*Value_NullValue); ok {
		return x.NullValue
	}
	return NullValue_NULL_VALUE
}

func (x *Value) GetDoubleValue() float64 {
	if x, ok := x.GetKind().(*Value_DoubleValue); ok {
		return x.DoubleValue
	}
	return 0
}

func (x *Value) GetFloatValue() float32 {
	if x, ok := x.GetKind().(*Value_FloatValue); ok {
		return x.FloatValue
	}
	return 0
}

func (x *Value) GetInt32Value() int32 {
	if x, ok := x.GetKind().(*Value_Int32Value); ok {
		return x.Int32Value
	}
	return 0
}

func (x *Value) GetInt64Value() int64 {
	if x, ok := x.GetKind().(*Value_Int64Value); ok {
		return x.Int64Value
	}
	return 0
}

func (x *Value) GetUint32Value() uint32 {
	if x, ok := x.GetKind().(*Value_Uint32Value); ok {
		return x.Uint32Value
	}
	return 0
}

func (x *Value) GetUint64Value() uint64 {
	if x, ok := x.GetKind().(*Value_Uint64Value); ok {
		return x.Uint64Value
	}
	return 0
}

func (x *Value) GetSint32Value() int32 {
	if x, ok := x.GetKind().(*Value_Sint32Value); ok {
		return x.Sint32Value
	}
	return 0
}

func (x *Value) GetSint64Value() int64 {
	if x, ok := x.GetKind().(*Value_Sint64Value); ok {
		return x.Sint64Value
	}
	return 0
}

func (x *Value) GetFixed32Value() uint32 {
	if x, ok := x.GetKind().(*Value_Fixed32Value); ok {
		return x.Fixed32Value
	}
	return 0
}

func (x *Value) GetFixed64Value() uint64 {
	if x, ok := x.GetKind().(*Value_Fixed64Value); ok {
		return x.Fixed64Value
	}
	return 0
}

func (x *Value) GetSfixed32Value() int32 {
	if x, ok := x.GetKind().(*Value_Sfixed32Value); ok {
		return x.Sfixed32Value
	}
	return 0
}

func (x *Value) GetSfixed64Value() int64 {
	if x, ok := x.GetKind().(*Value_Sfixed64Value); ok {
		return x.Sfixed64Value
	}
	return 0
}

func (x *Value) GetBoolValue() bool {
	if x, ok := x.GetKind().(*Value_BoolValue); ok {
		return x.BoolValue
	}
	return false
}

func (x *Value) GetStringValue() string {
	if x, ok := x.GetKind().(*Value_StringValue); ok {
		return x.StringValue
	}
	return ""
}

func (x *Value) GetBytesValue() []byte {
	if x, ok := x.GetKind().(*Value_BytesValue); ok {
		return x.BytesValue
	}
	return nil
}

func (x *Value) GetAnyValue() *anypb.Any {
	if x, ok := x.GetKind().(*Value_AnyValue); ok {
		return x.AnyValue
	}
	return nil
}

func (x *Value) GetStructValue() *Struct {
	if x, ok := x.GetKind().(*Value_StructValue); ok {
		return x.StructValue
	}
	return nil
}

func (x *Value) GetListValue() *ListValue {
	if x, ok := x.GetKind().(*Value_ListValue); ok {
		return x.ListValue
	}
	return nil
}

type isValue_Kind interface {
	isValue_Kind()
}

type Value_NullValue struct {
	// Represents a null value.
	NullValue NullValue `protobuf:"varint,1,opt,name=null_value,json=nullValue,proto3,enum=common.v1.NullValue,oneof"`
}

type Value_DoubleValue struct {
	// Represents a double value.
	DoubleValue float64 `protobuf:"fixed64,2,opt,name=double_value,json=doubleValue,proto3,oneof"`
}

type Value_FloatValue struct {
	// Represents a float value.
	FloatValue float32 `protobuf:"fixed32,3,opt,name=float_value,json=floatValue,proto3,oneof"`
}

type Value_Int32Value struct {
	// Represents an int32 value.
	Int32Value int32 `protobuf:"varint,4,opt,name=int32_value,json=int32Value,proto3,oneof"`
}

type Value_Int64Value struct {
	// Represents an int64 value.
	Int64Value int64 `protobuf:"varint,5,opt,name=int64_value,json=int64Value,proto3,oneof"`
}

type Value_Uint32Value struct {
	// Represents an uint32 value.
	Uint32Value uint32 `protobuf:"varint,6,opt,name=uint32_value,json=uint32Value,proto3,oneof"`
}

type Value_Uint64Value struct {
	// Represents an uint64 value.
	Uint64Value uint64 `protobuf:"varint,7,opt,name=uint64_value,json=uint64Value,proto3,oneof"`
}

type Value_Sint32Value struct {
	// Represents a sint32 value.
	Sint32Value int32 `protobuf:"zigzag32,8,opt,name=sint32_value,json=sint32Value,proto3,oneof"`
}

type Value_Sint64Value struct {
	// Represents a sint64 value.
	Sint64Value int64 `protobuf:"zigzag64,9,opt,name=sint64_value,json=sint64Value,proto3,oneof"`
}

type Value_Fixed32Value struct {
	// Represents a fixed32 value.
	Fixed32Value uint32 `protobuf:"fixed32,10,opt,name=fixed32_value,json=fixed32Value,proto3,oneof"`
}

type Value_Fixed64Value struct {
	// Represents a fixed64 value.
	Fixed64Value uint64 `protobuf:"fixed64,11,opt,name=fixed64_value,json=fixed64Value,proto3,oneof"`
}

type Value_Sfixed32Value struct {
	// Represents a sfixed32 value.
	Sfixed32Value int32 `protobuf:"fixed32,12,opt,name=sfixed32_value,json=sfixed32Value,proto3,oneof"`
}

type Value_Sfixed64Value struct {
	// Represents a sfixed64 value.
	Sfixed64Value int64 `protobuf:"fixed64,13,opt,name=sfixed64_value,json=sfixed64Value,proto3,oneof"`
}

type Value_BoolValue struct {
	// Represents a boolean value.
	BoolValue bool `protobuf:"varint,14,opt,name=bool_value,json=boolValue,proto3,oneof"`
}

type Value_StringValue struct {
	// Represents a string value.
	StringValue string `protobuf:"bytes,15,opt,name=string_value,json=stringValue,proto3,oneof"`
}

type Value_BytesValue struct {
	// Represents a bytes value.
	BytesValue []byte `protobuf:"bytes,16,opt,name=bytes_value,json=bytesValue,proto3,oneof"`
}

type Value_AnyValue struct {
	// Represents an Any value.
	AnyValue *anypb.Any `protobuf:"bytes,17,opt,name=any_value,json=anyValue,proto3,oneof"`
}

type Value_StructValue struct {
	// Represents a structured value.
	StructValue *Struct `protobuf:"bytes,18,opt,name=struct_value,json=structValue,proto3,oneof"`
}

type Value_ListValue struct {
	// Represents a repeated `Value`.
	ListValue *ListValue `protobuf:"bytes,19,opt,name=list_value,json=listValue,proto3,oneof"`
}

func (*Value_NullValue) isValue_Kind() {}

func (*Value_DoubleValue) isValue_Kind() {}

func (*Value_FloatValue) isValue_Kind() {}

func (*Value_Int32Value) isValue_Kind() {}

func (*Value_Int64Value) isValue_Kind() {}

func (*Value_Uint32Value) isValue_Kind() {}

func (*Value_Uint64Value) isValue_Kind() {}

func (*Value_Sint32Value) isValue_Kind() {}

func (*Value_Sint64Value) isValue_Kind() {}

func (*Value_Fixed32Value) isValue_Kind() {}

func (*Value_Fixed64Value) isValue_Kind() {}

func (*Value_Sfixed32Value) isValue_Kind() {}

func (*Value_Sfixed64Value) isValue_Kind() {}

func (*Value_BoolValue) isValue_Kind() {}

func (*Value_StringValue) isValue_Kind() {}

func (*Value_BytesValue) isValue_Kind() {}

func (*Value_AnyValue) isValue_Kind() {}

func (*Value_StructValue) isValue_Kind() {}

func (*Value_ListValue) isValue_Kind() {}

type Struct struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Unordered map of dynamically typed values.
	Fields map[string]*Value `protobuf:"bytes,1,rep,name=fields,proto3" json:"fields,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Struct) Reset() {
	*x = Struct{}
	if protoimpl.UnsafeEnabled {
		mi := &file_azarc_common_v1_struct_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Struct) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Struct) ProtoMessage() {}

func (x *Struct) ProtoReflect() protoreflect.Message {
	mi := &file_azarc_common_v1_struct_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Struct.ProtoReflect.Descriptor instead.
func (*Struct) Descriptor() ([]byte, []int) {
	return file_azarc_common_v1_struct_proto_rawDescGZIP(), []int{1}
}

func (x *Struct) GetFields() map[string]*Value {
	if x != nil {
		return x.Fields
	}
	return nil
}

type ListValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Repeated field of dynamically typed values.
	Values []*Value `protobuf:"bytes,1,rep,name=values,proto3" json:"values,omitempty"`
}

func (x *ListValue) Reset() {
	*x = ListValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_azarc_common_v1_struct_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ListValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListValue) ProtoMessage() {}

func (x *ListValue) ProtoReflect() protoreflect.Message {
	mi := &file_azarc_common_v1_struct_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListValue.ProtoReflect.Descriptor instead.
func (*ListValue) Descriptor() ([]byte, []int) {
	return file_azarc_common_v1_struct_proto_rawDescGZIP(), []int{2}
}

func (x *ListValue) GetValues() []*Value {
	if x != nil {
		return x.Values
	}
	return nil
}

var File_azarc_common_v1_struct_proto protoreflect.FileDescriptor

var file_azarc_common_v1_struct_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x61, 0x7a, 0x61, 0x72, 0x63, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76,
	0x31, 0x2f, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09,
	0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x95, 0x06, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x35,
	0x0a, 0x0a, 0x6e, 0x75, 0x6c, 0x6c, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x0e, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4e,
	0x75, 0x6c, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48, 0x00, 0x52, 0x09, 0x6e, 0x75, 0x6c, 0x6c,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x64, 0x6f, 0x75, 0x62, 0x6c, 0x65, 0x5f,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01, 0x48, 0x00, 0x52, 0x0b, 0x64,
	0x6f, 0x75, 0x62, 0x6c, 0x65, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21, 0x0a, 0x0b, 0x66, 0x6c,
	0x6f, 0x61, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x02, 0x48,
	0x00, 0x52, 0x0a, 0x66, 0x6c, 0x6f, 0x61, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21, 0x0a,
	0x0b, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01,
	0x28, 0x05, 0x48, 0x00, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x21, 0x0a, 0x0b, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x03, 0x48, 0x00, 0x52, 0x0a, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x75, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x5f, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x00, 0x52, 0x0b, 0x75, 0x69, 0x6e,
	0x74, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x75, 0x69, 0x6e, 0x74,
	0x36, 0x34, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x48, 0x00,
	0x52, 0x0b, 0x75, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x23, 0x0a,
	0x0c, 0x73, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x11, 0x48, 0x00, 0x52, 0x0b, 0x73, 0x69, 0x6e, 0x74, 0x33, 0x32, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x73, 0x69, 0x6e, 0x74, 0x36, 0x34, 0x5f, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x12, 0x48, 0x00, 0x52, 0x0b, 0x73, 0x69, 0x6e, 0x74,
	0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x25, 0x0a, 0x0d, 0x66, 0x69, 0x78, 0x65, 0x64,
	0x33, 0x32, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x07, 0x48, 0x00,
	0x52, 0x0c, 0x66, 0x69, 0x78, 0x65, 0x64, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x25,
	0x0a, 0x0d, 0x66, 0x69, 0x78, 0x65, 0x64, 0x36, 0x34, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18,
	0x0b, 0x20, 0x01, 0x28, 0x06, 0x48, 0x00, 0x52, 0x0c, 0x66, 0x69, 0x78, 0x65, 0x64, 0x36, 0x34,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x27, 0x0a, 0x0e, 0x73, 0x66, 0x69, 0x78, 0x65, 0x64, 0x33,
	0x32, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x0c, 0x20, 0x01, 0x28, 0x0f, 0x48, 0x00, 0x52,
	0x0d, 0x73, 0x66, 0x69, 0x78, 0x65, 0x64, 0x33, 0x32, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x27,
	0x0a, 0x0e, 0x73, 0x66, 0x69, 0x78, 0x65, 0x64, 0x36, 0x34, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65,
	0x18, 0x0d, 0x20, 0x01, 0x28, 0x10, 0x48, 0x00, 0x52, 0x0d, 0x73, 0x66, 0x69, 0x78, 0x65, 0x64,
	0x36, 0x34, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x1f, 0x0a, 0x0a, 0x62, 0x6f, 0x6f, 0x6c, 0x5f,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x0e, 0x20, 0x01, 0x28, 0x08, 0x48, 0x00, 0x52, 0x09, 0x62,
	0x6f, 0x6f, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x23, 0x0a, 0x0c, 0x73, 0x74, 0x72, 0x69,
	0x6e, 0x67, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00,
	0x52, 0x0b, 0x73, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x21, 0x0a,
	0x0b, 0x62, 0x79, 0x74, 0x65, 0x73, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x10, 0x20, 0x01,
	0x28, 0x0c, 0x48, 0x00, 0x52, 0x0a, 0x62, 0x79, 0x74, 0x65, 0x73, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x12, 0x33, 0x0a, 0x09, 0x61, 0x6e, 0x79, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x11, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41, 0x6e, 0x79, 0x48, 0x00, 0x52, 0x08, 0x61, 0x6e, 0x79,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x36, 0x0a, 0x0c, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x5f,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x12, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x48, 0x00,
	0x52, 0x0b, 0x73, 0x74, 0x72, 0x75, 0x63, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x35, 0x0a,
	0x0a, 0x6c, 0x69, 0x73, 0x74, 0x5f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x13, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x69,
	0x73, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x48, 0x00, 0x52, 0x09, 0x6c, 0x69, 0x73, 0x74, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x42, 0x06, 0x0a, 0x04, 0x6b, 0x69, 0x6e, 0x64, 0x22, 0x8c, 0x01, 0x0a,
	0x06, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x12, 0x35, 0x0a, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64,
	0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x2e, 0x46, 0x69, 0x65, 0x6c, 0x64,
	0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x66, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x1a, 0x4b,
	0x0a, 0x0b, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12,
	0x26, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x35, 0x0a, 0x09, 0x4c,
	0x69, 0x73, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x28, 0x0a, 0x06, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f,
	0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x06, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x73, 0x2a, 0x1b, 0x0a, 0x09, 0x4e, 0x75, 0x6c, 0x6c, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12,
	0x0e, 0x0a, 0x0a, 0x4e, 0x55, 0x4c, 0x4c, 0x5f, 0x56, 0x41, 0x4c, 0x55, 0x45, 0x10, 0x00, 0x42,
	0xac, 0x01, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76,
	0x31, 0x42, 0x0b, 0x53, 0x74, 0x72, 0x75, 0x63, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01,
	0x5a, 0x49, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x61, 0x7a, 0x61,
	0x72, 0x63, 0x2d, 0x69, 0x6f, 0x2f, 0x76, 0x74, 0x68, 0x2d, 0x66, 0x61, 0x61, 0x73, 0x2d, 0x73,
	0x64, 0x6b, 0x2d, 0x67, 0x6f, 0x2f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x67,
	0x65, 0x6e, 0x2f, 0x61, 0x7a, 0x61, 0x72, 0x63, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f,
	0x76, 0x31, 0x3b, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x43, 0x58,
	0x58, 0xaa, 0x02, 0x09, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x09,
	0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x15, 0x43, 0x6f, 0x6d, 0x6d,
	0x6f, 0x6e, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61, 0x74,
	0x61, 0xea, 0x02, 0x0a, 0x43, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x3a, 0x3a, 0x56, 0x31, 0x62, 0x06,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_azarc_common_v1_struct_proto_rawDescOnce sync.Once
	file_azarc_common_v1_struct_proto_rawDescData = file_azarc_common_v1_struct_proto_rawDesc
)

func file_azarc_common_v1_struct_proto_rawDescGZIP() []byte {
	file_azarc_common_v1_struct_proto_rawDescOnce.Do(func() {
		file_azarc_common_v1_struct_proto_rawDescData = protoimpl.X.CompressGZIP(file_azarc_common_v1_struct_proto_rawDescData)
	})
	return file_azarc_common_v1_struct_proto_rawDescData
}

var file_azarc_common_v1_struct_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_azarc_common_v1_struct_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_azarc_common_v1_struct_proto_goTypes = []interface{}{
	(NullValue)(0),    // 0: common.v1.NullValue
	(*Value)(nil),     // 1: common.v1.Value
	(*Struct)(nil),    // 2: common.v1.Struct
	(*ListValue)(nil), // 3: common.v1.ListValue
	nil,               // 4: common.v1.Struct.FieldsEntry
	(*anypb.Any)(nil), // 5: google.protobuf.Any
}
var file_azarc_common_v1_struct_proto_depIdxs = []int32{
	0, // 0: common.v1.Value.null_value:type_name -> common.v1.NullValue
	5, // 1: common.v1.Value.any_value:type_name -> google.protobuf.Any
	2, // 2: common.v1.Value.struct_value:type_name -> common.v1.Struct
	3, // 3: common.v1.Value.list_value:type_name -> common.v1.ListValue
	4, // 4: common.v1.Struct.fields:type_name -> common.v1.Struct.FieldsEntry
	1, // 5: common.v1.ListValue.values:type_name -> common.v1.Value
	1, // 6: common.v1.Struct.FieldsEntry.value:type_name -> common.v1.Value
	7, // [7:7] is the sub-list for method output_type
	7, // [7:7] is the sub-list for method input_type
	7, // [7:7] is the sub-list for extension type_name
	7, // [7:7] is the sub-list for extension extendee
	0, // [0:7] is the sub-list for field type_name
}

func init() { file_azarc_common_v1_struct_proto_init() }
func file_azarc_common_v1_struct_proto_init() {
	if File_azarc_common_v1_struct_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_azarc_common_v1_struct_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Value); i {
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
		file_azarc_common_v1_struct_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Struct); i {
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
		file_azarc_common_v1_struct_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ListValue); i {
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
	file_azarc_common_v1_struct_proto_msgTypes[0].OneofWrappers = []interface{}{
		(*Value_NullValue)(nil),
		(*Value_DoubleValue)(nil),
		(*Value_FloatValue)(nil),
		(*Value_Int32Value)(nil),
		(*Value_Int64Value)(nil),
		(*Value_Uint32Value)(nil),
		(*Value_Uint64Value)(nil),
		(*Value_Sint32Value)(nil),
		(*Value_Sint64Value)(nil),
		(*Value_Fixed32Value)(nil),
		(*Value_Fixed64Value)(nil),
		(*Value_Sfixed32Value)(nil),
		(*Value_Sfixed64Value)(nil),
		(*Value_BoolValue)(nil),
		(*Value_StringValue)(nil),
		(*Value_BytesValue)(nil),
		(*Value_AnyValue)(nil),
		(*Value_StructValue)(nil),
		(*Value_ListValue)(nil),
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_azarc_common_v1_struct_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_azarc_common_v1_struct_proto_goTypes,
		DependencyIndexes: file_azarc_common_v1_struct_proto_depIdxs,
		EnumInfos:         file_azarc_common_v1_struct_proto_enumTypes,
		MessageInfos:      file_azarc_common_v1_struct_proto_msgTypes,
	}.Build()
	File_azarc_common_v1_struct_proto = out.File
	file_azarc_common_v1_struct_proto_rawDesc = nil
	file_azarc_common_v1_struct_proto_goTypes = nil
	file_azarc_common_v1_struct_proto_depIdxs = nil
}
