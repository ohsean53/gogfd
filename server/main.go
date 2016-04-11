package main

import (
//flatbuffers "github.com/google/flatbuffers/go"
	"gogfd/flatbuffer/fbMessage"
	"math"
	"math/rand"
	"net"
	"runtime"
	"time"
	"gogfd/lib"
)


// server config
const (
	maxRoom = math.MaxInt32
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
	rooms = lib.NewSMap(lib.RWMutex)
	rand.Seed(time.Now().UTC().UnixNano())
}

func ClientSender(user *User, c net.Conn) {

	defer user.Leave()

	for {
		select {
		case <-user.exit:
		// when receive signal then finish the program
			if DEBUG {
				lib.Log("Leave user id :" + lib.Itoa64(user.userID))
			}
			return
		case m := <-user.recv:
		// on receive message
			if DEBUG {
				lib.Log("Client recv, user id : " + lib.Itoa64(user.userID))
			}
			_, err := c.Write(m.contents) // send data to client
			if err != nil {
				if DEBUG {
					lib.Log(err)
				}
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
			if DEBUG {
				lib.Log("Fail Stream read, err : ", err)
			}
			break
		}

		rawData := data[:n]
		message := fbMessage.GetRootAsMessage(rawData, 0)
		messageType := message.BodyType()
		handler, ok := msgHandler[messageType]

		if DEBUG {
			lib.Log("req : type : ", messageType)
		}
		if ok {
			ret := handler(user, message) // calling proper handler function
			if !ret {
				lib.Log("Fail handler : ", messageType)
			}
		} else {
			if DEBUG {
				lib.Log("Fail no function defined for type : ", messageType)
			}
			break
		}
	}

	// fail read
	user.exit <- true
}

func main() {

	runtime.GOMAXPROCS(runtime.NumCPU())
	ln, err := net.Listen("tcp", ":8000") // using TCP protocol over 8000 port
	if err != nil {
		if DEBUG {
			lib.Log(err)
		}
		return
	}

	InitRooms()

	defer ln.Close() // reserve listen wait close
	for {
		conn, err := ln.Accept() // server accept client connection -> return connection
		if err != nil {
			lib.Log("Fail Accept err : ", err)
			conn.Close()
			continue
		}
		if DEBUG {
			lib.Log("New Connection: ", conn.RemoteAddr())
		}

		lib.WriteScribe("access", "test")
		user := NewUser(0, nil) // empty user data
		go ClientReader(user, conn)
		go ClientSender(user, conn)
	}
}
