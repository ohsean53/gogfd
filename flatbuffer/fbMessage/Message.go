// automatically generated, do not modify

package fbMessage

import (
	flatbuffers "github.com/google/flatbuffers/go"
)
type Message struct {
	_tab flatbuffers.Table
}

func GetRootAsMessage(buf []byte, offset flatbuffers.UOffsetT) *Message {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &Message{}
	x.Init(buf, n + offset)
	return x
}

func (rcv *Message) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *Message) BodyType() byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.GetByte(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *Message) Body(obj *flatbuffers.Table) bool {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		rcv._tab.Union(obj, o)
		return true
	}
	return false
}

func MessageStart(builder *flatbuffers.Builder) {
	builder.StartObject(2)
}
func MessageAddBodyType(builder *flatbuffers.Builder, bodyType byte) {
	builder.PrependByteSlot(0, bodyType, 0)
}
func MessageAddBody(builder *flatbuffers.Builder, body flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(body), 0)
}
func MessageEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
