// automatically generated, do not modify

package fbMessage

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type NotifyAction1 struct {
	_tab flatbuffers.Table
}

func (rcv *NotifyAction1) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *NotifyAction1) UserID() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func NotifyAction1Start(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func NotifyAction1AddUserID(builder *flatbuffers.Builder, userID int64) {
	builder.PrependInt64Slot(0, userID, 0)
}
func NotifyAction1End(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
