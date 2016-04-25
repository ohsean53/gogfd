// automatically generated, do not modify

package fbMessage

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type ResRoomList struct {
	_tab flatbuffers.Table
}

func (rcv *ResRoomList) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ResRoomList) UserID() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ResRoomList) RoomIDs(j int) int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetInt64(a + flatbuffers.UOffsetT(j * 8))
	}
	return 0
}

func (rcv *ResRoomList) RoomIDsLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ResRoomList) ResultCode() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 100
}

func ResRoomListStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func ResRoomListAddUserID(builder *flatbuffers.Builder, userID int64) {
	builder.PrependInt64Slot(0, userID, 0)
}
func ResRoomListAddRoomIDs(builder *flatbuffers.Builder, roomIDs flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(roomIDs), 0)
}
func ResRoomListStartRoomIDsVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(8, numElems, 8)
}
func ResRoomListAddResultCode(builder *flatbuffers.Builder, resultCode int32) {
	builder.PrependInt32Slot(2, resultCode, 100)
}
func ResRoomListEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
