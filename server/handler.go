package main

import (
	"gogfd/lib"
	"gogfd/flatbuffer/fbMessage"
	flatbuffers "github.com/google/flatbuffers/go"
	"strconv"
)

type MsgHandlerFunc func(user *User, message *fbMessage.Message) bool

var msgHandler = map[byte]MsgHandlerFunc{
	fbMessage.MessageBodyReqLogin:          LoginHandler,
	fbMessage.MessageBodyReqCreate:         CreateHandler,
	fbMessage.MessageBodyReqJoin:           JoinHandler,
	fbMessage.MessageBodyReqAction1:     Action1Handler,
	fbMessage.MessageBodyReqQuit:           QuitHandler,
	fbMessage.MessageBodyReqRoomList:       RoomListHandler,
}

func LoginHandler(user *User, message *fbMessage.Message) bool {
	unionTable := new(flatbuffers.Table)
	if message.Body(unionTable) == false {
		lib.Log("message.Body fail")
		return false
	}

	req := new(fbMessage.ReqLogin)
	req.Init(unionTable.Bytes, unionTable.Pos)
	user.userID = req.UserID()
	lib.Log("server recv user id : ", user.userID)

	builder := flatbuffers.NewBuilder(0)
	fbMessage.ResLoginStart(builder)
	fbMessage.ResLoginAddUserID(builder, user.userID)
	fbMessage.ResLoginAddResultCode(builder, fbMessage.ResultCodeSuccess)
	resBody := fbMessage.ResLoginEnd(builder)

	fbMessage.MessageStart(builder)
	fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyResLogin)
	fbMessage.MessageAddBody(builder, resBody)
	msg := fbMessage.MessageEnd(builder)

	builder.Finish(msg)

	user.Push(NewMessage(user.userID, fbMessage.MessageBodyResLogin, builder.FinishedBytes()))
	return true
}

func CreateHandler(user *User, message *fbMessage.Message) bool {

	unionTable := new(flatbuffers.Table)
	if message.Body(unionTable) == false {
		lib.Log("message.Body fail")
		return false
	}

	req := new(fbMessage.ReqCreate)
	req.Init(unionTable.Bytes, unionTable.Pos)

	if user.userID != req.UserID() {
		if DEBUG {
			lib.Log("Fail room create, user id missmatch")
		}
		return false
	}


	// room create
	roomID := GetRandomRoomID()
	r := NewRoom(roomID)
	r.users.Set(user.userID, user) // insert user
	user.room = r                  // set room
	rooms.Set(roomID, r)           // set room into global shared map
	if DEBUG {
		lib.Log("Get rand room id : ", lib.Itoa64(roomID))
	}

	builder := flatbuffers.NewBuilder(0)
	fbMessage.ResCreateStart(builder)
	fbMessage.ResCreateAddUserID(builder, user.userID)
	fbMessage.ResCreateAddRoomID(builder, roomID)
	fbMessage.ResCreateAddResultCode(builder, fbMessage.ResultCodeSuccess)
	resBody := fbMessage.ResCreateEnd(builder)

	fbMessage.MessageStart(builder)
	fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyResCreate)
	fbMessage.MessageAddBody(builder, resBody)
	msg := fbMessage.MessageEnd(builder)

	builder.Finish(msg)

	if DEBUG {
		lib.Log("Room create, room id : ", lib.Itoa64(roomID))
	}

	user.Push(NewMessage(user.userID, fbMessage.MessageBodyResCreate, builder.FinishedBytes()))
	return true
}

func JoinHandler(user *User, message *fbMessage.Message) bool {

	// request message
	unionTable := new(flatbuffers.Table)
	if message.Body(unionTable) == false {
		lib.Log("message.Body fail")
		return false
	}

	req := new(fbMessage.ReqJoin)
	req.Init(unionTable.Bytes, unionTable.Pos)
	roomID := req.RoomID()
	value, ok := rooms.Get(roomID)

	if !ok {
		if DEBUG {
			lib.Log("Fail room join, room does not exist, room id : ", lib.Itoa64(roomID))
		}
		return false
	}

	r := value.(*Room)
	r.users.Set(user.userID, user)
	user.room = r

	// broadcast message
	builder := flatbuffers.NewBuilder(0)
	fbMessage.NotifyJoinStart(builder)
	fbMessage.NotifyJoinAddUserID(builder, user.userID)
	fbMessage.NotifyJoinAddRoomID(builder, roomID)
	notifyBody := fbMessage.NotifyJoinEnd(builder)

	// message body
	fbMessage.MessageStart(builder)
	fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyNotifyJoin)
	fbMessage.MessageAddBody(builder, notifyBody)
	msg := fbMessage.MessageEnd(builder)
	builder.Finish(msg)
	user.SendToAll(NewMessage(user.userID, fbMessage.MessageBodyNotifyJoin, builder.FinishedBytes()))

	// response message
	builder = flatbuffers.NewBuilder(0)

	////////////////////////////////////////////////////////
	// set member vector (must do init first)
	////////////////////////////////////////////////////////
	userCount := r.users.Count()
	lib.Log("room id : " + lib.Itoa64(roomID) + " member count : " + strconv.Itoa(userCount))
	fbMessage.ResJoinStartMembersVector(builder, userCount)
	for userID, _ := range r.users.Map() {
		lib.Log("PrependInt64 room id : " + lib.Itoa64(roomID) + " user id : " + lib.Itoa64(userID))
		builder.PrependInt64(userID)
	}
	member := builder.EndVector(userCount)

	// response message
	fbMessage.ResJoinStart(builder)
	fbMessage.ResJoinAddUserID(builder, user.userID)
	fbMessage.ResJoinAddMembers(builder, member)
	fbMessage.ResJoinAddResultCode(builder, fbMessage.ResultCodeSuccess)
	resBody := fbMessage.ResJoinEnd(builder)

	// message body
	fbMessage.MessageStart(builder)
	fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyResJoin)
	fbMessage.MessageAddBody(builder, resBody)
	msg = fbMessage.MessageEnd(builder)

	builder.Finish(msg)
	user.Push(NewMessage(user.userID, fbMessage.MessageBodyResJoin, builder.FinishedBytes()))

	return true
}

func Action1Handler(user *User, message *fbMessage.Message) bool {

	// request message
	unionTable := new(flatbuffers.Table)
	if message.Body(unionTable) == false {
		lib.Log("message.Body fail")
		return false
	}

	req := new(fbMessage.ReqAction1)
	req.Init(unionTable.Bytes, unionTable.Pos)

	if nil == user.room {
		if DEBUG {
			lib.Log("Fail action 1, not exist room info")
		}
		return false
	}
	if user.userID != req.UserID() {
		if DEBUG {
			lib.Log("Fail action 1, user id missmatch")
		}
		return false
	}

	// broadcast message
	builder := flatbuffers.NewBuilder(0)
	fbMessage.NotifyAction1Start(builder)
	fbMessage.NotifyAction1AddUserID(builder, user.userID)
	notifyBody := fbMessage.NotifyAction1End(builder)

	// message body
	fbMessage.MessageStart(builder)
	fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyNotifyAction1)
	fbMessage.MessageAddBody(builder, notifyBody)
	msg := fbMessage.MessageEnd(builder)
	builder.Finish(msg)
	user.SendToAll(NewMessage(user.userID, fbMessage.MessageBodyNotifyAction1, builder.FinishedBytes()))

	// response message
	builder = flatbuffers.NewBuilder(0)

	// response message
	fbMessage.ResAction1Start(builder)
	fbMessage.ResAction1AddUserID(builder, user.userID)
	fbMessage.ResAction1AddResultCode(builder, fbMessage.ResultCodeSuccess)
	resBody := fbMessage.ResAction1End(builder)

	// message body
	fbMessage.MessageStart(builder)
	fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyResAction1)
	fbMessage.MessageAddBody(builder, resBody)
	msg = fbMessage.MessageEnd(builder)

	builder.Finish(msg)
	user.Push(NewMessage(user.userID, fbMessage.MessageBodyResAction1, builder.FinishedBytes()))
	return true
}

func QuitHandler(user *User, message *fbMessage.Message) bool {

	unionTable := new(flatbuffers.Table)
	if message.Body(unionTable) == false {
		lib.Log("message.Body fail")
		return false
	}

	req := new(fbMessage.ReqQuit)
	req.Init(unionTable.Bytes, unionTable.Pos)
	user.userID = req.UserID()
	lib.Log("QuitHandler server recv user id : ", user.userID)

	builder := flatbuffers.NewBuilder(0)
	fbMessage.ResQuitStart(builder)
	fbMessage.ResQuitAddUserID(builder, user.userID)
	fbMessage.ResQuitAddResultCode(builder, fbMessage.ResultCodeSuccess)
	resBody := fbMessage.ResQuitEnd(builder)

	fbMessage.MessageStart(builder)
	fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyResQuit)
	fbMessage.MessageAddBody(builder, resBody)
	msg := fbMessage.MessageEnd(builder)

	builder.Finish(msg)

	user.Push(NewMessage(user.userID, fbMessage.MessageBodyResQuit, builder.FinishedBytes()))

	// same act user.Leave()
	user.exit <- true
	return true
}

func RoomListHandler(user *User, message *fbMessage.Message) bool {
	// request message
	unionTable := new(flatbuffers.Table)
	if message.Body(unionTable) == false {
		lib.Log("message.Body fail")
		return false
	}

	req := new(fbMessage.ReqRoomList)
	req.Init(unionTable.Bytes, unionTable.Pos)

	// response message
	builder := flatbuffers.NewBuilder(0)

	////////////////////////////////////////////////////////
	// set room ids vector (must do init first)
	////////////////////////////////////////////////////////
	roomCount := rooms.Count()
	fbMessage.ResRoomListStartRoomIDsVector(builder, roomCount)
	for roomID, _ := range rooms.Map() {
		builder.PrependInt64(roomID)
	}
	rooms := builder.EndVector(roomCount)

	// response message
	fbMessage.ResRoomListStart(builder)
	fbMessage.ResRoomListAddUserID(builder, user.userID)
	fbMessage.ResRoomListAddRoomIDs(builder, rooms)
	fbMessage.ResRoomListAddResultCode(builder, fbMessage.ResultCodeSuccess)
	resBody := fbMessage.ResRoomListEnd(builder)

	// message body
	fbMessage.MessageStart(builder)
	fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyResRoomList)
	fbMessage.MessageAddBody(builder, resBody)
	msg := fbMessage.MessageEnd(builder)

	builder.Finish(msg)
	user.Push(NewMessage(user.userID, fbMessage.MessageBodyResRoomList, builder.FinishedBytes()))
	return true
}
