package main

import (
//"math"
	"gogfd/lib"
	"gogfd/logger"
	"sync"
)

type Room struct {
	users    lib.SharedMap
	messages chan *UserMessage // broadcast message channel
	roomID   int64
}

var (
	lastUseRoomID int64
	sc sync.RWMutex
)

func GetAutoIncRoomID() int64 {
	sc.Lock()
	defer sc.Unlock()
	lastUseRoomID += 1
	return lastUseRoomID
}

func GetRoom(roomID int64) (r *Room) {
	value, ok := rooms.Get(roomID)
	if !ok {
		logger.Log(logger.ERROR, "err: not exist room : ", roomID)
	}
	r = value.(*Room)
	return
}

func GenerateRoomID() (uuid int64) {
	for {
		// TODO change room-id generate strategy
		//uuid = int64(lib.RandInt32(1, 100))
		uuid = GetAutoIncRoomID()
		if _, ok := rooms.Get(uuid); ok {
			logger.Log(logger.ERROR, "err: exist same room id")
			continue
		}
		return
	}
}

// Room construct
func NewRoom(roomID int64) (r *Room) {
	r = new(Room)
	r.messages = make(chan *UserMessage)
	r.users = lib.NewSMap(lib.RW_MUTEX) // global shared map
	r.roomID = roomID
	go r.RoomMessageLoop() // for broadcast message
	return
}

func (r *Room) RoomMessageLoop() {
	// when messages channel is closed then "for-loop" will be break
	for m := range r.messages {
		for userID, _ := range r.users.Map() {
			if userID == m.userID {
				continue
			}
			value, ok := r.users.Get(userID)
			if !ok {
				continue
			}
			user := value.(*User)
			logger.Log(logger.DEBUG, "Push message for broadcast :" + lib.Itoa64(user.userID))
			user.Push(m)
		}
	}
}

func (r *Room) getRoomUsers() (userIDs []int64) {
	userIDs = r.users.GetKeys()
	return userIDs
}

func (r *Room) Leave(userID int64) {
	r.users.Remove(userID)

	if r.IsEmptyRoom() == false {
		return
	}

	close(r.messages)    // close broadcast channel
	rooms.Remove(r.roomID)
}

func (r *Room) IsEmptyRoom() bool {
	if r.users.Count() == 0 {
		return true
	}
	return false
}
