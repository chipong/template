// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: table.proto

package oz

import (
	proto "github.com/golang/protobuf/proto"
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

type Template struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string         `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	GroupId  string         `protobuf:"bytes,2,opt,name=group_id,json=groupId,proto3" json:"group_id,omitempty"`
	RewardId string         `protobuf:"bytes,3,opt,name=reward_id,json=rewardId,proto3" json:"reward_id,omitempty"`
	Count    int64          `protobuf:"varint,4,opt,name=count,proto3" json:"count,omitempty"`
	Enum     TemplateEnum_T `protobuf:"varint,5,opt,name=enum,proto3,enum=oz.TemplateEnum_T" json:"enum,omitempty"`
}

func (x *Template) Reset() {
	*x = Template{}
	if protoimpl.UnsafeEnabled {
		mi := &file_table_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Template) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Template) ProtoMessage() {}

func (x *Template) ProtoReflect() protoreflect.Message {
	mi := &file_table_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Template.ProtoReflect.Descriptor instead.
func (*Template) Descriptor() ([]byte, []int) {
	return file_table_proto_rawDescGZIP(), []int{0}
}

func (x *Template) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Template) GetGroupId() string {
	if x != nil {
		return x.GroupId
	}
	return ""
}

func (x *Template) GetRewardId() string {
	if x != nil {
		return x.RewardId
	}
	return ""
}

func (x *Template) GetCount() int64 {
	if x != nil {
		return x.Count
	}
	return 0
}

func (x *Template) GetEnum() TemplateEnum_T {
	if x != nil {
		return x.Enum
	}
	return TemplateEnum_NONE
}

var File_table_proto protoreflect.FileDescriptor

var file_table_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x74, 0x61, 0x62, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02, 0x6f,
	0x7a, 0x1a, 0x0c, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22,
	0x90, 0x01, 0x0a, 0x08, 0x54, 0x65, 0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x12, 0x0e, 0x0a, 0x02,
	0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x19, 0x0a, 0x08,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x5f, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x67, 0x72, 0x6f, 0x75, 0x70, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x72, 0x65, 0x77, 0x61, 0x72,
	0x64, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x77, 0x61,
	0x72, 0x64, 0x49, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x05, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x26, 0x0a, 0x04, 0x65, 0x6e,
	0x75, 0x6d, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x6f, 0x7a, 0x2e, 0x54, 0x65,
	0x6d, 0x70, 0x6c, 0x61, 0x74, 0x65, 0x45, 0x6e, 0x75, 0x6d, 0x2e, 0x54, 0x52, 0x04, 0x65, 0x6e,
	0x75, 0x6d, 0x42, 0x08, 0x50, 0x01, 0x5a, 0x04, 0x2e, 0x2f, 0x6f, 0x7a, 0x62, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_table_proto_rawDescOnce sync.Once
	file_table_proto_rawDescData = file_table_proto_rawDesc
)

func file_table_proto_rawDescGZIP() []byte {
	file_table_proto_rawDescOnce.Do(func() {
		file_table_proto_rawDescData = protoimpl.X.CompressGZIP(file_table_proto_rawDescData)
	})
	return file_table_proto_rawDescData
}

var file_table_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_table_proto_goTypes = []interface{}{
	(*Template)(nil),    // 0: oz.Template
	(TemplateEnum_T)(0), // 1: oz.TemplateEnum.T
}
var file_table_proto_depIdxs = []int32{
	1, // 0: oz.Template.enum:type_name -> oz.TemplateEnum.T
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_table_proto_init() }
func file_table_proto_init() {
	if File_table_proto != nil {
		return
	}
	file_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_table_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Template); i {
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
			RawDescriptor: file_table_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_table_proto_goTypes,
		DependencyIndexes: file_table_proto_depIdxs,
		MessageInfos:      file_table_proto_msgTypes,
	}.Build()
	File_table_proto = out.File
	file_table_proto_rawDesc = nil
	file_table_proto_goTypes = nil
	file_table_proto_depIdxs = nil
}
