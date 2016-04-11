package main

import (
	"gogfd/lib"
	"gogfd/flatbuffer/fbMessage"
	"github.com/google/flatbuffers/go"
)

type User struct {
	userID int64
	room   *Room
	recv   chan *UserMessage
	exit   chan bool // signal
}

func NewUser(uid int64, room *Room) *User {
	return &User{
		userID: uid,
		recv:   make(chan *UserMessage),
		exit:   make(chan bool, 1),
		room:   room,
	}
}

func (u *User) Leave() {

	if DEBUG {
		lib.Log("Leave user id : ", lib.Itoa64(u.userID))
	}

	if u.room != nil {
		if DEBUG {
			lib.Log("Leave room id : ", lib.Itoa64(u.room.roomID))
		}

		builder := flatbuffers.NewBuilder(0)
		fbMessage.NotifyQuitStart(builder)
		fbMessage.NotifyQuitAddUserID(builder, u.userID)
		fbMessage.NotifyQuitAddRoomID(builder, u.room.roomID)
		notifyBody := fbMessage.NotifyQuitEnd(builder)

		// message body
		fbMessage.MessageStart(builder)
		fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyNotifyQuit)
		fbMessage.MessageAddBody(builder, notifyBody)
		msg := fbMessage.MessageEnd(builder)
		builder.Finish(msg)

		// race condition : broadcast goroutine vs ClientSender goroutine
		u.room.Leave(u.userID)

		// notify all members in the room
		u.SendToAll(NewMessage(u.userID, fbMessage.MessageBodyNotifyQuit, builder.FinishedBytes()))
		if DEBUG {
			lib.Log("NotifyQuit message send")
		}
	}

	if DEBUG {
		lib.Log("Leave func end")
	}
}

func (u *User) Push(m *UserMessage) {
	u.recv <- m // send message to user
}

func (u *User) SendToAll(m *UserMessage) {
	if u.room.IsEmptyRoom() == false {
		u.room.messages <- m
	}
}
