package main

import (
//flatbuffers "github.com/google/flatbuffers/go"
	"gogfd/flatbuffer/fbMessage"
	"math/rand"
	"net"
	"runtime"
	"time"
	"gogfd/lib"
	"gogfd/logger"
	"gogfd/config"
)



// global variable
var (
	rooms lib.SharedMap
)

type UserMessage struct {
	userID    int64  // sender
	msgType   byte
	timestamp int    // send time
	contents  []byte // serialized google protocol-buffer message
}

func NewMessage(userID int64, eventType byte, msg []byte) *UserMessage {
	return &UserMessage{
		userID,
		eventType,
		int(time.Now().Unix()),
		msg,
	}
}

func InitRooms() {
	rooms = lib.NewSMap(lib.RW_MUTEX)
	rand.Seed(time.Now().UTC().UnixNano())
}

func ClientSender(user *User, c net.Conn) {

	defer user.Leave()

	for {
		select {
		case <-user.exit:
		// when receive signal then finish the program
			logger.Log(logger.DEBUG, "Leave user id :" + lib.Itoa64(user.userID))
			return
		case m := <-user.recv:
		// on receive message
			logger.Log(logger.DEBUG, "Client recv, user id : " + lib.Itoa64(user.userID))
			_, err := c.Write(m.contents) // send data to client
			if err != nil {
				logger.Log(logger.ERROR, err)
				return
			}
		}
	}
}

func ClientReader(user *User, c net.Conn) {

	data := make([]byte, 4096) // 4096 byte slice (dynamic resize)

	for {
		n, err := c.Read(data)
		if err != nil {
			logger.Log(logger.DEBUG, "Fail Stream read, err : ", err)
			break
		}

		rawData := data[:n]
		message := fbMessage.GetRootAsMessage(rawData, 0)
		messageType := message.BodyType()
		handler, ok := msgHandler[messageType]

		logger.Log(logger.DEBUG, "req : type : ", messageType)

		if ok {
			ret := handler(user, message) // calling proper handler function
			if !ret {
				logger.Log(logger.ERROR, "Fail handler : ", messageType)
			}
		} else {
			logger.Log(logger.ERROR, "Fail no function defined for type : ", messageType)
			break
		}
	}

	// fail read
	user.exit <- true
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	ln, err := net.Listen("tcp", ":" + config.SERVER_PORT) // using TCP protocol over 8000 port
	if err != nil {
		logger.Log(logger.CRITICAL, err)
		return
	}

	InitRooms()

	defer ln.Close() // reserve listen wait close
	for {
		conn, err := ln.Accept() // server accept client connection -> return connection
		if err != nil {
			logger.Log(logger.CRITICAL, "Fail Accept err : ", err)
			conn.Close()
			continue
		}
		logger.Log(logger.INFO, "New Connection: ", conn.RemoteAddr())
		logger.WriteScribe("access", "test")
		user := NewUser(0, nil) // empty user data
		go ClientReader(user, conn)
		go ClientSender(user, conn)
	}
}
