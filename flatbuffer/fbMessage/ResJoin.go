// automatically generated, do not modify

package fbMessage

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type ResJoin struct {
	_tab flatbuffers.Table
}

func (rcv *ResJoin) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *ResJoin) UserID() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *ResJoin) Members(j int) int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		a := rcv._tab.Vector(o)
		return rcv._tab.GetInt64(a + flatbuffers.UOffsetT(j * 8))
	}
	return 0
}

func (rcv *ResJoin) MembersLength() int {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.VectorLen(o)
	}
	return 0
}

func (rcv *ResJoin) ResultCode() int32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt32(o + rcv._tab.Pos)
	}
	return 100
}

func ResJoinStart(builder *flatbuffers.Builder) {
	builder.StartObject(3)
}
func ResJoinAddUserID(builder *flatbuffers.Builder, userID int64) {
	builder.PrependInt64Slot(0, userID, 0)
}
func ResJoinAddMembers(builder *flatbuffers.Builder, members flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(members), 0)
}
func ResJoinStartMembersVector(builder *flatbuffers.Builder, numElems int) flatbuffers.UOffsetT {
	return builder.StartVector(8, numElems, 8)
}
func ResJoinAddResultCode(builder *flatbuffers.Builder, resultCode int32) {
	builder.PrependInt32Slot(2, resultCode, 100)
}
func ResJoinEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
