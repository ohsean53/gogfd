// automatically generated, do not modify

package fbMessage

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type ResCreate struct {
	_tab flatbuffers.Table
}

func (rcv *ResCreate) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ResCreate) UserID() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ResCreate) RoomID() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ResCreate) ResultCode() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 100
}

func ResCreateStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func ResCreateAddUserID(builder *flatbuffers.Builder, userID int64) {
	builder.PrependInt64Slot(0, userID, 0)
}
func ResCreateAddRoomID(builder *flatbuffers.Builder, roomID int64) {
	builder.PrependInt64Slot(1, roomID, 0)
}
func ResCreateAddResultCode(builder *flatbuffers.Builder, resultCode int32) {
	builder.PrependInt32Slot(2, resultCode, 100)
}
func ResCreateEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
