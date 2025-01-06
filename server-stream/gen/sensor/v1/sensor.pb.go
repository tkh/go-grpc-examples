// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.1
// 	protoc        (unknown)
// source: sensor/v1/sensor.proto

package sensorv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Request to stream a specific sensor's data.
type StreamSensorDataRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	SensorId      string                 `protobuf:"bytes,1,opt,name=sensor_id,json=sensorId,proto3" json:"sensor_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamSensorDataRequest) Reset() {
	*x = StreamSensorDataRequest{}
	mi := &file_sensor_v1_sensor_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamSensorDataRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamSensorDataRequest) ProtoMessage() {}

func (x *StreamSensorDataRequest) ProtoReflect() protoreflect.Message {
	mi := &file_sensor_v1_sensor_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamSensorDataRequest.ProtoReflect.Descriptor instead.
func (*StreamSensorDataRequest) Descriptor() ([]byte, []int) {
	return file_sensor_v1_sensor_proto_rawDescGZIP(), []int{0}
}

func (x *StreamSensorDataRequest) GetSensorId() string {
	if x != nil {
		return x.SensorId
	}
	return ""
}

// Timestamped sensor data.
type StreamSensorDataResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Temperature   float32                `protobuf:"fixed32,1,opt,name=temperature,proto3" json:"temperature,omitempty"`
	Humidity      float32                `protobuf:"fixed32,2,opt,name=humidity,proto3" json:"humidity,omitempty"`
	Timestamp     *timestamppb.Timestamp `protobuf:"bytes,3,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *StreamSensorDataResponse) Reset() {
	*x = StreamSensorDataResponse{}
	mi := &file_sensor_v1_sensor_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *StreamSensorDataResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*StreamSensorDataResponse) ProtoMessage() {}

func (x *StreamSensorDataResponse) ProtoReflect() protoreflect.Message {
	mi := &file_sensor_v1_sensor_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use StreamSensorDataResponse.ProtoReflect.Descriptor instead.
func (*StreamSensorDataResponse) Descriptor() ([]byte, []int) {
	return file_sensor_v1_sensor_proto_rawDescGZIP(), []int{1}
}

func (x *StreamSensorDataResponse) GetTemperature() float32 {
	if x != nil {
		return x.Temperature
	}
	return 0
}

func (x *StreamSensorDataResponse) GetHumidity() float32 {
	if x != nil {
		return x.Humidity
	}
	return 0
}

func (x *StreamSensorDataResponse) GetTimestamp() *timestamppb.Timestamp {
	if x != nil {
		return x.Timestamp
	}
	return nil
}

var File_sensor_v1_sensor_proto protoreflect.FileDescriptor

var file_sensor_v1_sensor_proto_rawDesc = []byte{
	0x0a, 0x16, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x2f, 0x76, 0x31, 0x2f, 0x73, 0x65, 0x6e, 0x73,
	0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72,
	0x2e, 0x76, 0x31, 0x1a, 0x1f, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x2f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x22, 0x36, 0x0a, 0x17, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x65,
	0x6e, 0x73, 0x6f, 0x72, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12,
	0x1b, 0x0a, 0x09, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x08, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x49, 0x64, 0x22, 0x92, 0x01, 0x0a,
	0x18, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x44, 0x61, 0x74,
	0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x0b, 0x74, 0x65, 0x6d,
	0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x02, 0x52, 0x0b,
	0x74, 0x65, 0x6d, 0x70, 0x65, 0x72, 0x61, 0x74, 0x75, 0x72, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x68,
	0x75, 0x6d, 0x69, 0x64, 0x69, 0x74, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x02, 0x52, 0x08, 0x68,
	0x75, 0x6d, 0x69, 0x64, 0x69, 0x74, 0x79, 0x12, 0x38, 0x0a, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73,
	0x74, 0x61, 0x6d, 0x70, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x1a, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x54, 0x69, 0x6d,
	0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d,
	0x70, 0x32, 0x70, 0x0a, 0x0d, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x5f, 0x0a, 0x10, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x65, 0x6e, 0x73,
	0x6f, 0x72, 0x44, 0x61, 0x74, 0x61, 0x12, 0x22, 0x2e, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x2e,
	0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x44,
	0x61, 0x74, 0x61, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e, 0x73, 0x65, 0x6e,
	0x73, 0x6f, 0x72, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x53, 0x65, 0x6e,
	0x73, 0x6f, 0x72, 0x44, 0x61, 0x74, 0x61, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22,
	0x00, 0x30, 0x01, 0x42, 0x46, 0x5a, 0x44, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x74, 0x6b, 0x68, 0x2f, 0x67, 0x6f, 0x2d, 0x67, 0x72, 0x70, 0x63, 0x2d, 0x65, 0x78,
	0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x65, 0x72, 0x2d, 0x73, 0x74,
	0x72, 0x65, 0x61, 0x6d, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x2f,
	0x76, 0x31, 0x3b, 0x73, 0x65, 0x6e, 0x73, 0x6f, 0x72, 0x76, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_sensor_v1_sensor_proto_rawDescOnce sync.Once
	file_sensor_v1_sensor_proto_rawDescData = file_sensor_v1_sensor_proto_rawDesc
)

func file_sensor_v1_sensor_proto_rawDescGZIP() []byte {
	file_sensor_v1_sensor_proto_rawDescOnce.Do(func() {
		file_sensor_v1_sensor_proto_rawDescData = protoimpl.X.CompressGZIP(file_sensor_v1_sensor_proto_rawDescData)
	})
	return file_sensor_v1_sensor_proto_rawDescData
}

var file_sensor_v1_sensor_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_sensor_v1_sensor_proto_goTypes = []any{
	(*StreamSensorDataRequest)(nil),  // 0: sensor.v1.StreamSensorDataRequest
	(*StreamSensorDataResponse)(nil), // 1: sensor.v1.StreamSensorDataResponse
	(*timestamppb.Timestamp)(nil),    // 2: google.protobuf.Timestamp
}
var file_sensor_v1_sensor_proto_depIdxs = []int32{
	2, // 0: sensor.v1.StreamSensorDataResponse.timestamp:type_name -> google.protobuf.Timestamp
	0, // 1: sensor.v1.SensorService.StreamSensorData:input_type -> sensor.v1.StreamSensorDataRequest
	1, // 2: sensor.v1.SensorService.StreamSensorData:output_type -> sensor.v1.StreamSensorDataResponse
	2, // [2:3] is the sub-list for method output_type
	1, // [1:2] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_sensor_v1_sensor_proto_init() }
func file_sensor_v1_sensor_proto_init() {
	if File_sensor_v1_sensor_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_sensor_v1_sensor_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_sensor_v1_sensor_proto_goTypes,
		DependencyIndexes: file_sensor_v1_sensor_proto_depIdxs,
		MessageInfos:      file_sensor_v1_sensor_proto_msgTypes,
	}.Build()
	File_sensor_v1_sensor_proto = out.File
	file_sensor_v1_sensor_proto_rawDesc = nil
	file_sensor_v1_sensor_proto_goTypes = nil
	file_sensor_v1_sensor_proto_depIdxs = nil
}
