// automatically generated, do not modify

package fbMessage

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type NotifyJoin struct {
	_tab flatbuffers.Table
}

func (rcv *NotifyJoin) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *NotifyJoin) UserID() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *NotifyJoin) RoomID() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func NotifyJoinStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func NotifyJoinAddUserID(builder *flatbuffers.Builder, userID int64) {
	builder.PrependInt64Slot(0, userID, 0)
}
func NotifyJoinAddRoomID(builder *flatbuffers.Builder, roomID int64) {
	builder.PrependInt64Slot(1, roomID, 0)
}
func NotifyJoinEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
