// automatically generated, do not modify

package fbMessage

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type ResAction1 struct {
	_tab flatbuffers.Table
}

func (rcv *ResAction1) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ResAction1) UserID() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ResAction1) ResultCode() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 100
}

func ResAction1Start(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func ResAction1AddUserID(builder *flatbuffers.Builder, userID int64) {
	builder.PrependInt64Slot(0, userID, 0)
}
func ResAction1AddResultCode(builder *flatbuffers.Builder, resultCode int32) {
	builder.PrependInt32Slot(1, resultCode, 100)
}
func ResAction1End(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
