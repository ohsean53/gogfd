// automatically generated, do not modify

package fbMessage

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type ReqQuit struct {
	_tab flatbuffers.Table
}

func (rcv *ReqQuit) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ReqQuit) UserID() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func ReqQuitStart(builder *flatbuffers.Builder) {
	builder.StartObject(1)
}
func ReqQuitAddUserID(builder *flatbuffers.Builder, userID int64) {
	builder.PrependInt64Slot(0, userID, 0)
}
func ReqQuitEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
