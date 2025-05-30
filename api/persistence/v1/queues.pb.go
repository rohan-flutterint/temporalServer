// Code generated by protoc-gen-go. DO NOT EDIT.
// plugins:
// 	protoc-gen-go
// 	protoc
// source: temporal/server/api/persistence/v1/queues.proto

package persistence

import (
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"

	v1 "go.temporal.io/api/common/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type QueueState struct {
	state                        protoimpl.MessageState      `protogen:"open.v1"`
	ReaderStates                 map[int64]*QueueReaderState `protobuf:"bytes,1,rep,name=reader_states,json=readerStates,proto3" json:"reader_states,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	ExclusiveReaderHighWatermark *TaskKey                    `protobuf:"bytes,2,opt,name=exclusive_reader_high_watermark,json=exclusiveReaderHighWatermark,proto3" json:"exclusive_reader_high_watermark,omitempty"`
	unknownFields                protoimpl.UnknownFields
	sizeCache                    protoimpl.SizeCache
}

func (x *QueueState) Reset() {
	*x = QueueState{}
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueueState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueState) ProtoMessage() {}

func (x *QueueState) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueState.ProtoReflect.Descriptor instead.
func (*QueueState) Descriptor() ([]byte, []int) {
	return file_temporal_server_api_persistence_v1_queues_proto_rawDescGZIP(), []int{0}
}

func (x *QueueState) GetReaderStates() map[int64]*QueueReaderState {
	if x != nil {
		return x.ReaderStates
	}
	return nil
}

func (x *QueueState) GetExclusiveReaderHighWatermark() *TaskKey {
	if x != nil {
		return x.ExclusiveReaderHighWatermark
	}
	return nil
}

type QueueReaderState struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Scopes        []*QueueSliceScope     `protobuf:"bytes,1,rep,name=scopes,proto3" json:"scopes,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QueueReaderState) Reset() {
	*x = QueueReaderState{}
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueueReaderState) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueReaderState) ProtoMessage() {}

func (x *QueueReaderState) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueReaderState.ProtoReflect.Descriptor instead.
func (*QueueReaderState) Descriptor() ([]byte, []int) {
	return file_temporal_server_api_persistence_v1_queues_proto_rawDescGZIP(), []int{1}
}

func (x *QueueReaderState) GetScopes() []*QueueSliceScope {
	if x != nil {
		return x.Scopes
	}
	return nil
}

type QueueSliceScope struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Range         *QueueSliceRange       `protobuf:"bytes,1,opt,name=range,proto3" json:"range,omitempty"`
	Predicate     *Predicate             `protobuf:"bytes,2,opt,name=predicate,proto3" json:"predicate,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QueueSliceScope) Reset() {
	*x = QueueSliceScope{}
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueueSliceScope) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueSliceScope) ProtoMessage() {}

func (x *QueueSliceScope) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueSliceScope.ProtoReflect.Descriptor instead.
func (*QueueSliceScope) Descriptor() ([]byte, []int) {
	return file_temporal_server_api_persistence_v1_queues_proto_rawDescGZIP(), []int{2}
}

func (x *QueueSliceScope) GetRange() *QueueSliceRange {
	if x != nil {
		return x.Range
	}
	return nil
}

func (x *QueueSliceScope) GetPredicate() *Predicate {
	if x != nil {
		return x.Predicate
	}
	return nil
}

type QueueSliceRange struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	InclusiveMin  *TaskKey               `protobuf:"bytes,1,opt,name=inclusive_min,json=inclusiveMin,proto3" json:"inclusive_min,omitempty"`
	ExclusiveMax  *TaskKey               `protobuf:"bytes,2,opt,name=exclusive_max,json=exclusiveMax,proto3" json:"exclusive_max,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QueueSliceRange) Reset() {
	*x = QueueSliceRange{}
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueueSliceRange) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueueSliceRange) ProtoMessage() {}

func (x *QueueSliceRange) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueueSliceRange.ProtoReflect.Descriptor instead.
func (*QueueSliceRange) Descriptor() ([]byte, []int) {
	return file_temporal_server_api_persistence_v1_queues_proto_rawDescGZIP(), []int{3}
}

func (x *QueueSliceRange) GetInclusiveMin() *TaskKey {
	if x != nil {
		return x.InclusiveMin
	}
	return nil
}

func (x *QueueSliceRange) GetExclusiveMax() *TaskKey {
	if x != nil {
		return x.ExclusiveMax
	}
	return nil
}

type ReadQueueMessagesNextPageToken struct {
	state             protoimpl.MessageState `protogen:"open.v1"`
	LastReadMessageId int64                  `protobuf:"varint,1,opt,name=last_read_message_id,json=lastReadMessageId,proto3" json:"last_read_message_id,omitempty"`
	unknownFields     protoimpl.UnknownFields
	sizeCache         protoimpl.SizeCache
}

func (x *ReadQueueMessagesNextPageToken) Reset() {
	*x = ReadQueueMessagesNextPageToken{}
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ReadQueueMessagesNextPageToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ReadQueueMessagesNextPageToken) ProtoMessage() {}

func (x *ReadQueueMessagesNextPageToken) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ReadQueueMessagesNextPageToken.ProtoReflect.Descriptor instead.
func (*ReadQueueMessagesNextPageToken) Descriptor() ([]byte, []int) {
	return file_temporal_server_api_persistence_v1_queues_proto_rawDescGZIP(), []int{4}
}

func (x *ReadQueueMessagesNextPageToken) GetLastReadMessageId() int64 {
	if x != nil {
		return x.LastReadMessageId
	}
	return 0
}

type ListQueuesNextPageToken struct {
	state               protoimpl.MessageState `protogen:"open.v1"`
	LastReadQueueNumber int64                  `protobuf:"varint,1,opt,name=last_read_queue_number,json=lastReadQueueNumber,proto3" json:"last_read_queue_number,omitempty"`
	unknownFields       protoimpl.UnknownFields
	sizeCache           protoimpl.SizeCache
}

func (x *ListQueuesNextPageToken) Reset() {
	*x = ListQueuesNextPageToken{}
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListQueuesNextPageToken) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListQueuesNextPageToken) ProtoMessage() {}

func (x *ListQueuesNextPageToken) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListQueuesNextPageToken.ProtoReflect.Descriptor instead.
func (*ListQueuesNextPageToken) Descriptor() ([]byte, []int) {
	return file_temporal_server_api_persistence_v1_queues_proto_rawDescGZIP(), []int{5}
}

func (x *ListQueuesNextPageToken) GetLastReadQueueNumber() int64 {
	if x != nil {
		return x.LastReadQueueNumber
	}
	return 0
}

// HistoryTask represents an internal history service task for a particular shard. We use a blob because there is no
// common proto for all task proto types.
type HistoryTask struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// shard_id that this task belonged to when it was created. Technically, you can derive this from the task data
	// blob, but it's useful to have it here for quick access and to avoid deserializing the blob. Note that this may be
	// different from the shard id of this task in the current cluster because it could have come from a cluster with a
	// different shard id. This will always be the shard id of the task in its original cluster.
	ShardId int32 `protobuf:"varint,1,opt,name=shard_id,json=shardId,proto3" json:"shard_id,omitempty"`
	// blob that contains the history task proto. There is a GoLang-specific generic deserializer for this blob, but
	// there is no common proto for all task proto types, so deserializing in other languages will require a custom
	// switch on the task category, which should be available from the metadata for the queue that this task came from.
	Blob          *v1.DataBlob `protobuf:"bytes,2,opt,name=blob,proto3" json:"blob,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *HistoryTask) Reset() {
	*x = HistoryTask{}
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *HistoryTask) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HistoryTask) ProtoMessage() {}

func (x *HistoryTask) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HistoryTask.ProtoReflect.Descriptor instead.
func (*HistoryTask) Descriptor() ([]byte, []int) {
	return file_temporal_server_api_persistence_v1_queues_proto_rawDescGZIP(), []int{6}
}

func (x *HistoryTask) GetShardId() int32 {
	if x != nil {
		return x.ShardId
	}
	return 0
}

func (x *HistoryTask) GetBlob() *v1.DataBlob {
	if x != nil {
		return x.Blob
	}
	return nil
}

type QueuePartition struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// min_message_id is less than or equal to the id of every message in the queue. The min_message_id is mainly used to
	// skip over tombstones in Cassandra: let's say we deleted the first 1K messages from a queue with 1.1K messages. If
	//
	//	an operator asked for the first 100 messages, without the min_message_id, we would have to scan over the 1K
	//
	// tombstone rows before we could return the 100 messages. With the min_message_id, we can skip over all of the
	// tombstones by specifying message_id >= queue.min_message_id. Note: it is possible for this to be less than the id
	// of the lowest message in the queue temporarily because we delete messages before we update the queue metadata.
	// However, such errors surface to clients with an "Unavailable" code, so clients retry, and the id should be updated
	// soon. Additionally, we only use min_message_id to skip over tombstones, so it will only affect read performance,
	// not correctness.
	MinMessageId  int64 `protobuf:"varint,1,opt,name=min_message_id,json=minMessageId,proto3" json:"min_message_id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *QueuePartition) Reset() {
	*x = QueuePartition{}
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[7]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *QueuePartition) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QueuePartition) ProtoMessage() {}

func (x *QueuePartition) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[7]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QueuePartition.ProtoReflect.Descriptor instead.
func (*QueuePartition) Descriptor() ([]byte, []int) {
	return file_temporal_server_api_persistence_v1_queues_proto_rawDescGZIP(), []int{7}
}

func (x *QueuePartition) GetMinMessageId() int64 {
	if x != nil {
		return x.MinMessageId
	}
	return 0
}

type Queue struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// A map from partition index (0-based) to the partition metadata.
	Partitions    map[int32]*QueuePartition `protobuf:"bytes,1,rep,name=partitions,proto3" json:"partitions,omitempty" protobuf_key:"varint,1,opt,name=key" protobuf_val:"bytes,2,opt,name=value"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Queue) Reset() {
	*x = Queue{}
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[8]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Queue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Queue) ProtoMessage() {}

func (x *Queue) ProtoReflect() protoreflect.Message {
	mi := &file_temporal_server_api_persistence_v1_queues_proto_msgTypes[8]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Queue.ProtoReflect.Descriptor instead.
func (*Queue) Descriptor() ([]byte, []int) {
	return file_temporal_server_api_persistence_v1_queues_proto_rawDescGZIP(), []int{8}
}

func (x *Queue) GetPartitions() map[int32]*QueuePartition {
	if x != nil {
		return x.Partitions
	}
	return nil
}

var File_temporal_server_api_persistence_v1_queues_proto protoreflect.FileDescriptor

const file_temporal_server_api_persistence_v1_queues_proto_rawDesc = "" +
	"\n" +
	"/temporal/server/api/persistence/v1/queues.proto\x12\"temporal.server.api.persistence.v1\x1a$temporal/api/common/v1/message.proto\x1a3temporal/server/api/persistence/v1/predicates.proto\x1a.temporal/server/api/persistence/v1/tasks.proto\"\xde\x02\n" +
	"\n" +
	"QueueState\x12e\n" +
	"\rreader_states\x18\x01 \x03(\v2@.temporal.server.api.persistence.v1.QueueState.ReaderStatesEntryR\freaderStates\x12r\n" +
	"\x1fexclusive_reader_high_watermark\x18\x02 \x01(\v2+.temporal.server.api.persistence.v1.TaskKeyR\x1cexclusiveReaderHighWatermark\x1au\n" +
	"\x11ReaderStatesEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\x03R\x03key\x12J\n" +
	"\x05value\x18\x02 \x01(\v24.temporal.server.api.persistence.v1.QueueReaderStateR\x05value:\x028\x01\"_\n" +
	"\x10QueueReaderState\x12K\n" +
	"\x06scopes\x18\x01 \x03(\v23.temporal.server.api.persistence.v1.QueueSliceScopeR\x06scopes\"\xa9\x01\n" +
	"\x0fQueueSliceScope\x12I\n" +
	"\x05range\x18\x01 \x01(\v23.temporal.server.api.persistence.v1.QueueSliceRangeR\x05range\x12K\n" +
	"\tpredicate\x18\x02 \x01(\v2-.temporal.server.api.persistence.v1.PredicateR\tpredicate\"\xb5\x01\n" +
	"\x0fQueueSliceRange\x12P\n" +
	"\rinclusive_min\x18\x01 \x01(\v2+.temporal.server.api.persistence.v1.TaskKeyR\finclusiveMin\x12P\n" +
	"\rexclusive_max\x18\x02 \x01(\v2+.temporal.server.api.persistence.v1.TaskKeyR\fexclusiveMax\"Q\n" +
	"\x1eReadQueueMessagesNextPageToken\x12/\n" +
	"\x14last_read_message_id\x18\x01 \x01(\x03R\x11lastReadMessageId\"N\n" +
	"\x17ListQueuesNextPageToken\x123\n" +
	"\x16last_read_queue_number\x18\x01 \x01(\x03R\x13lastReadQueueNumber\"^\n" +
	"\vHistoryTask\x12\x19\n" +
	"\bshard_id\x18\x01 \x01(\x05R\ashardId\x124\n" +
	"\x04blob\x18\x02 \x01(\v2 .temporal.api.common.v1.DataBlobR\x04blob\"6\n" +
	"\x0eQueuePartition\x12$\n" +
	"\x0emin_message_id\x18\x01 \x01(\x03R\fminMessageId\"\xd5\x01\n" +
	"\x05Queue\x12Y\n" +
	"\n" +
	"partitions\x18\x01 \x03(\v29.temporal.server.api.persistence.v1.Queue.PartitionsEntryR\n" +
	"partitions\x1aq\n" +
	"\x0fPartitionsEntry\x12\x10\n" +
	"\x03key\x18\x01 \x01(\x05R\x03key\x12H\n" +
	"\x05value\x18\x02 \x01(\v22.temporal.server.api.persistence.v1.QueuePartitionR\x05value:\x028\x01B6Z4go.temporal.io/server/api/persistence/v1;persistenceb\x06proto3"

var (
	file_temporal_server_api_persistence_v1_queues_proto_rawDescOnce sync.Once
	file_temporal_server_api_persistence_v1_queues_proto_rawDescData []byte
)

func file_temporal_server_api_persistence_v1_queues_proto_rawDescGZIP() []byte {
	file_temporal_server_api_persistence_v1_queues_proto_rawDescOnce.Do(func() {
		file_temporal_server_api_persistence_v1_queues_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_temporal_server_api_persistence_v1_queues_proto_rawDesc), len(file_temporal_server_api_persistence_v1_queues_proto_rawDesc)))
	})
	return file_temporal_server_api_persistence_v1_queues_proto_rawDescData
}

var file_temporal_server_api_persistence_v1_queues_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_temporal_server_api_persistence_v1_queues_proto_goTypes = []any{
	(*QueueState)(nil),                     // 0: temporal.server.api.persistence.v1.QueueState
	(*QueueReaderState)(nil),               // 1: temporal.server.api.persistence.v1.QueueReaderState
	(*QueueSliceScope)(nil),                // 2: temporal.server.api.persistence.v1.QueueSliceScope
	(*QueueSliceRange)(nil),                // 3: temporal.server.api.persistence.v1.QueueSliceRange
	(*ReadQueueMessagesNextPageToken)(nil), // 4: temporal.server.api.persistence.v1.ReadQueueMessagesNextPageToken
	(*ListQueuesNextPageToken)(nil),        // 5: temporal.server.api.persistence.v1.ListQueuesNextPageToken
	(*HistoryTask)(nil),                    // 6: temporal.server.api.persistence.v1.HistoryTask
	(*QueuePartition)(nil),                 // 7: temporal.server.api.persistence.v1.QueuePartition
	(*Queue)(nil),                          // 8: temporal.server.api.persistence.v1.Queue
	nil,                                    // 9: temporal.server.api.persistence.v1.QueueState.ReaderStatesEntry
	nil,                                    // 10: temporal.server.api.persistence.v1.Queue.PartitionsEntry
	(*TaskKey)(nil),                        // 11: temporal.server.api.persistence.v1.TaskKey
	(*Predicate)(nil),                      // 12: temporal.server.api.persistence.v1.Predicate
	(*v1.DataBlob)(nil),                    // 13: temporal.api.common.v1.DataBlob
}
var file_temporal_server_api_persistence_v1_queues_proto_depIdxs = []int32{
	9,  // 0: temporal.server.api.persistence.v1.QueueState.reader_states:type_name -> temporal.server.api.persistence.v1.QueueState.ReaderStatesEntry
	11, // 1: temporal.server.api.persistence.v1.QueueState.exclusive_reader_high_watermark:type_name -> temporal.server.api.persistence.v1.TaskKey
	2,  // 2: temporal.server.api.persistence.v1.QueueReaderState.scopes:type_name -> temporal.server.api.persistence.v1.QueueSliceScope
	3,  // 3: temporal.server.api.persistence.v1.QueueSliceScope.range:type_name -> temporal.server.api.persistence.v1.QueueSliceRange
	12, // 4: temporal.server.api.persistence.v1.QueueSliceScope.predicate:type_name -> temporal.server.api.persistence.v1.Predicate
	11, // 5: temporal.server.api.persistence.v1.QueueSliceRange.inclusive_min:type_name -> temporal.server.api.persistence.v1.TaskKey
	11, // 6: temporal.server.api.persistence.v1.QueueSliceRange.exclusive_max:type_name -> temporal.server.api.persistence.v1.TaskKey
	13, // 7: temporal.server.api.persistence.v1.HistoryTask.blob:type_name -> temporal.api.common.v1.DataBlob
	10, // 8: temporal.server.api.persistence.v1.Queue.partitions:type_name -> temporal.server.api.persistence.v1.Queue.PartitionsEntry
	1,  // 9: temporal.server.api.persistence.v1.QueueState.ReaderStatesEntry.value:type_name -> temporal.server.api.persistence.v1.QueueReaderState
	7,  // 10: temporal.server.api.persistence.v1.Queue.PartitionsEntry.value:type_name -> temporal.server.api.persistence.v1.QueuePartition
	11, // [11:11] is the sub-list for method output_type
	11, // [11:11] is the sub-list for method input_type
	11, // [11:11] is the sub-list for extension type_name
	11, // [11:11] is the sub-list for extension extendee
	0,  // [0:11] is the sub-list for field type_name
}

func init() { file_temporal_server_api_persistence_v1_queues_proto_init() }
func file_temporal_server_api_persistence_v1_queues_proto_init() {
	if File_temporal_server_api_persistence_v1_queues_proto != nil {
		return
	}
	file_temporal_server_api_persistence_v1_predicates_proto_init()
	file_temporal_server_api_persistence_v1_tasks_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_temporal_server_api_persistence_v1_queues_proto_rawDesc), len(file_temporal_server_api_persistence_v1_queues_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_temporal_server_api_persistence_v1_queues_proto_goTypes,
		DependencyIndexes: file_temporal_server_api_persistence_v1_queues_proto_depIdxs,
		MessageInfos:      file_temporal_server_api_persistence_v1_queues_proto_msgTypes,
	}.Build()
	File_temporal_server_api_persistence_v1_queues_proto = out.File
	file_temporal_server_api_persistence_v1_queues_proto_goTypes = nil
	file_temporal_server_api_persistence_v1_queues_proto_depIdxs = nil
}
