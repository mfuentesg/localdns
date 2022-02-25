// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: pb/dns.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type RecordType int32

const (
	RecordType_A     RecordType = 0
	RecordType_AAAA  RecordType = 1
	RecordType_CNAME RecordType = 2
	RecordType_MX    RecordType = 3
	RecordType_TXT   RecordType = 4
	RecordType_NS    RecordType = 5
	RecordType_SOA   RecordType = 6
	RecordType_SRV   RecordType = 7
	RecordType_PTR   RecordType = 8
)

// Enum value maps for RecordType.
var (
	RecordType_name = map[int32]string{
		0: "A",
		1: "AAAA",
		2: "CNAME",
		3: "MX",
		4: "TXT",
		5: "NS",
		6: "SOA",
		7: "SRV",
		8: "PTR",
	}
	RecordType_value = map[string]int32{
		"A":     0,
		"AAAA":  1,
		"CNAME": 2,
		"MX":    3,
		"TXT":   4,
		"NS":    5,
		"SOA":   6,
		"SRV":   7,
		"PTR":   8,
	}
)

func (x RecordType) Enum() *RecordType {
	p := new(RecordType)
	*p = x
	return p
}

func (x RecordType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (RecordType) Descriptor() protoreflect.EnumDescriptor {
	return file_pb_dns_proto_enumTypes[0].Descriptor()
}

func (RecordType) Type() protoreflect.EnumType {
	return &file_pb_dns_proto_enumTypes[0]
}

func (x RecordType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use RecordType.Descriptor instead.
func (RecordType) EnumDescriptor() ([]byte, []int) {
	return file_pb_dns_proto_rawDescGZIP(), []int{0}
}

type Record struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type   RecordType `protobuf:"varint,1,opt,name=type,proto3,enum=localdns.dns.RecordType" json:"type,omitempty"`
	Domain string     `protobuf:"bytes,3,opt,name=domain,proto3" json:"domain,omitempty"`
	Ip     string     `protobuf:"bytes,4,opt,name=ip,proto3" json:"ip,omitempty"`
	Ttl    int32      `protobuf:"varint,5,opt,name=ttl,proto3" json:"ttl,omitempty"`
}

func (x *Record) Reset() {
	*x = Record{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_dns_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Record) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Record) ProtoMessage() {}

func (x *Record) ProtoReflect() protoreflect.Message {
	mi := &file_pb_dns_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Record.ProtoReflect.Descriptor instead.
func (*Record) Descriptor() ([]byte, []int) {
	return file_pb_dns_proto_rawDescGZIP(), []int{0}
}

func (x *Record) GetType() RecordType {
	if x != nil {
		return x.Type
	}
	return RecordType_A
}

func (x *Record) GetDomain() string {
	if x != nil {
		return x.Domain
	}
	return ""
}

func (x *Record) GetIp() string {
	if x != nil {
		return x.Ip
	}
	return ""
}

func (x *Record) GetTtl() int32 {
	if x != nil {
		return x.Ttl
	}
	return 0
}

type RecordList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Records []*Record `protobuf:"bytes,1,rep,name=records,proto3" json:"records,omitempty"`
	Page    int32     `protobuf:"varint,2,opt,name=page,proto3" json:"page,omitempty"`
	Total   int32     `protobuf:"varint,3,opt,name=total,proto3" json:"total,omitempty"`
}

func (x *RecordList) Reset() {
	*x = RecordList{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_dns_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecordList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordList) ProtoMessage() {}

func (x *RecordList) ProtoReflect() protoreflect.Message {
	mi := &file_pb_dns_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordList.ProtoReflect.Descriptor instead.
func (*RecordList) Descriptor() ([]byte, []int) {
	return file_pb_dns_proto_rawDescGZIP(), []int{1}
}

func (x *RecordList) GetRecords() []*Record {
	if x != nil {
		return x.Records
	}
	return nil
}

func (x *RecordList) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *RecordList) GetTotal() int32 {
	if x != nil {
		return x.Total
	}
	return 0
}

type RecordsFilter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Page   int32  `protobuf:"varint,1,opt,name=page,proto3" json:"page,omitempty"`
	Offset int32  `protobuf:"varint,2,opt,name=offset,proto3" json:"offset,omitempty"`
	Query  string `protobuf:"bytes,3,opt,name=query,proto3" json:"query,omitempty"`
}

func (x *RecordsFilter) Reset() {
	*x = RecordsFilter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_dns_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RecordsFilter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RecordsFilter) ProtoMessage() {}

func (x *RecordsFilter) ProtoReflect() protoreflect.Message {
	mi := &file_pb_dns_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RecordsFilter.ProtoReflect.Descriptor instead.
func (*RecordsFilter) Descriptor() ([]byte, []int) {
	return file_pb_dns_proto_rawDescGZIP(), []int{2}
}

func (x *RecordsFilter) GetPage() int32 {
	if x != nil {
		return x.Page
	}
	return 0
}

func (x *RecordsFilter) GetOffset() int32 {
	if x != nil {
		return x.Offset
	}
	return 0
}

func (x *RecordsFilter) GetQuery() string {
	if x != nil {
		return x.Query
	}
	return ""
}

var File_pb_dns_proto protoreflect.FileDescriptor

var file_pb_dns_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x70, 0x62, 0x2f, 0x64, 0x6e, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x0c,
	0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x64, 0x6e, 0x73, 0x2e, 0x64, 0x6e, 0x73, 0x1a, 0x1b, 0x67, 0x6f,
	0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d,
	0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x70, 0x0a, 0x06, 0x52, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x12, 0x2c, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0e, 0x32, 0x18, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x64, 0x6e, 0x73, 0x2e, 0x64, 0x6e, 0x73,
	0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70,
	0x65, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x70, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x02, 0x69, 0x70, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x74, 0x6c,
	0x18, 0x05, 0x20, 0x01, 0x28, 0x05, 0x52, 0x03, 0x74, 0x74, 0x6c, 0x22, 0x66, 0x0a, 0x0a, 0x52,
	0x65, 0x63, 0x6f, 0x72, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x2e, 0x0a, 0x07, 0x72, 0x65, 0x63,
	0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x6c, 0x6f, 0x63,
	0x61, 0x6c, 0x64, 0x6e, 0x73, 0x2e, 0x64, 0x6e, 0x73, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x52, 0x07, 0x72, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x14, 0x0a,
	0x05, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x74, 0x6f,
	0x74, 0x61, 0x6c, 0x22, 0x51, 0x0a, 0x0d, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x46, 0x69,
	0x6c, 0x74, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x05, 0x52, 0x04, 0x70, 0x61, 0x67, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x66, 0x66, 0x73,
	0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x06, 0x6f, 0x66, 0x66, 0x73, 0x65, 0x74,
	0x12, 0x14, 0x0a, 0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x05, 0x71, 0x75, 0x65, 0x72, 0x79, 0x2a, 0x5c, 0x0a, 0x0a, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x54, 0x79, 0x70, 0x65, 0x12, 0x05, 0x0a, 0x01, 0x41, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x41,
	0x41, 0x41, 0x41, 0x10, 0x01, 0x12, 0x09, 0x0a, 0x05, 0x43, 0x4e, 0x41, 0x4d, 0x45, 0x10, 0x02,
	0x12, 0x06, 0x0a, 0x02, 0x4d, 0x58, 0x10, 0x03, 0x12, 0x07, 0x0a, 0x03, 0x54, 0x58, 0x54, 0x10,
	0x04, 0x12, 0x06, 0x0a, 0x02, 0x4e, 0x53, 0x10, 0x05, 0x12, 0x07, 0x0a, 0x03, 0x53, 0x4f, 0x41,
	0x10, 0x06, 0x12, 0x07, 0x0a, 0x03, 0x53, 0x52, 0x56, 0x10, 0x07, 0x12, 0x07, 0x0a, 0x03, 0x50,
	0x54, 0x52, 0x10, 0x08, 0x32, 0x85, 0x02, 0x0a, 0x0a, 0x44, 0x6e, 0x73, 0x53, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x12, 0x37, 0x0a, 0x09, 0x41, 0x64, 0x64, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64,
	0x12, 0x14, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x64, 0x6e, 0x73, 0x2e, 0x64, 0x6e, 0x73, 0x2e,
	0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x1a, 0x14, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x64, 0x6e,
	0x73, 0x2e, 0x64, 0x6e, 0x73, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x3a, 0x0a, 0x0c,
	0x55, 0x70, 0x64, 0x61, 0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x14, 0x2e, 0x6c,
	0x6f, 0x63, 0x61, 0x6c, 0x64, 0x6e, 0x73, 0x2e, 0x64, 0x6e, 0x73, 0x2e, 0x52, 0x65, 0x63, 0x6f,
	0x72, 0x64, 0x1a, 0x14, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x64, 0x6e, 0x73, 0x2e, 0x64, 0x6e,
	0x73, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x3c, 0x0a, 0x0c, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x12, 0x14, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c,
	0x64, 0x6e, 0x73, 0x2e, 0x64, 0x6e, 0x73, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x1a, 0x16,
	0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x44, 0x0a, 0x0b, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x65,
	0x63, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x1b, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x64, 0x6e, 0x73,
	0x2e, 0x64, 0x6e, 0x73, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x73, 0x46, 0x69, 0x6c, 0x74,
	0x65, 0x72, 0x1a, 0x18, 0x2e, 0x6c, 0x6f, 0x63, 0x61, 0x6c, 0x64, 0x6e, 0x73, 0x2e, 0x64, 0x6e,
	0x73, 0x2e, 0x52, 0x65, 0x63, 0x6f, 0x72, 0x64, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x07, 0x5a, 0x05,
	0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_dns_proto_rawDescOnce sync.Once
	file_pb_dns_proto_rawDescData = file_pb_dns_proto_rawDesc
)

func file_pb_dns_proto_rawDescGZIP() []byte {
	file_pb_dns_proto_rawDescOnce.Do(func() {
		file_pb_dns_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_dns_proto_rawDescData)
	})
	return file_pb_dns_proto_rawDescData
}

var file_pb_dns_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pb_dns_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_pb_dns_proto_goTypes = []interface{}{
	(RecordType)(0),       // 0: localdns.dns.RecordType
	(*Record)(nil),        // 1: localdns.dns.Record
	(*RecordList)(nil),    // 2: localdns.dns.RecordList
	(*RecordsFilter)(nil), // 3: localdns.dns.RecordsFilter
	(*emptypb.Empty)(nil), // 4: google.protobuf.Empty
}
var file_pb_dns_proto_depIdxs = []int32{
	0, // 0: localdns.dns.Record.type:type_name -> localdns.dns.RecordType
	1, // 1: localdns.dns.RecordList.records:type_name -> localdns.dns.Record
	1, // 2: localdns.dns.DnsService.AddRecord:input_type -> localdns.dns.Record
	1, // 3: localdns.dns.DnsService.UpdateRecord:input_type -> localdns.dns.Record
	1, // 4: localdns.dns.DnsService.DeleteRecord:input_type -> localdns.dns.Record
	3, // 5: localdns.dns.DnsService.ListRecords:input_type -> localdns.dns.RecordsFilter
	1, // 6: localdns.dns.DnsService.AddRecord:output_type -> localdns.dns.Record
	1, // 7: localdns.dns.DnsService.UpdateRecord:output_type -> localdns.dns.Record
	4, // 8: localdns.dns.DnsService.DeleteRecord:output_type -> google.protobuf.Empty
	2, // 9: localdns.dns.DnsService.ListRecords:output_type -> localdns.dns.RecordList
	6, // [6:10] is the sub-list for method output_type
	2, // [2:6] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_pb_dns_proto_init() }
func file_pb_dns_proto_init() {
	if File_pb_dns_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_dns_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Record); i {
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
		file_pb_dns_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecordList); i {
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
		file_pb_dns_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RecordsFilter); i {
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
			RawDescriptor: file_pb_dns_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_dns_proto_goTypes,
		DependencyIndexes: file_pb_dns_proto_depIdxs,
		EnumInfos:         file_pb_dns_proto_enumTypes,
		MessageInfos:      file_pb_dns_proto_msgTypes,
	}.Build()
	File_pb_dns_proto = out.File
	file_pb_dns_proto_rawDesc = nil
	file_pb_dns_proto_goTypes = nil
	file_pb_dns_proto_depIdxs = nil
}
