package main

import (
	"fmt"
	"net"
	"gogfd/flatbuffer/fbMessage"
	"github.com/google/flatbuffers/go"
	"gogfd/lib"
	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"log"
	"strconv"
	"os"
)

type MsgHandlerFunc func(message *fbMessage.Message) bool

var msgHandler = map[byte]MsgHandlerFunc{
	fbMessage.MessageBodyResLogin:          ResLogin,
	fbMessage.MessageBodyResCreate:         ResCreate,
	fbMessage.MessageBodyResJoin:           ResJoin,
	fbMessage.MessageBodyNotifyJoin:        NotifyJoinHandler,
	fbMessage.MessageBodyNotifyQuit:        NotifyQuitHandler,
	fbMessage.MessageBodyResAction1:        ResAction1,
	fbMessage.MessageBodyNotifyAction1:     NotifyAction1Handler,
	fbMessage.MessageBodyResQuit:           ResQuit,
	fbMessage.MessageBodyResRoomList:       ResRoomList,
}

var inputString string

func main() {
	client, err := net.Dial("tcp", "127.0.0.1:8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer client.Close()

	data := make([]byte, 4096)
	//exit := make(chan bool, 1)

	var userID int64
	var mw *walk.MainWindow

	if err := (MainWindow{
		AssignTo: &mw,
		Title:    "Walk LogView Example",
		MinSize:  Size{320, 240},
		Size:     Size{600, 400},
		Layout:   VBox{},
		Children: []Widget{
			HSplitter{
				Children: []Widget{
					PushButton{
						Text: "login",
						OnClicked: func() {
							if cmd, err := RunUserIdDialog(mw); err != nil {
								log.Print(err)
							} else if cmd == walk.DlgCmdOK {
								log.Println("dlg msg : " + inputString)
								num, err := strconv.Atoi(inputString)
								lib.CheckError(err)
								userID = int64(num)
								ReqLogin(client, userID, data)
							}
						},
					},
					PushButton{
						Text: "room create",
						OnClicked: func() {
							log.Println("req create user id : ", userID)
							ReqCreate(client, userID, data)
						},
					},
					PushButton{
						Text: "room list",
						OnClicked: func() {
							log.Println("room list user id : ", userID)
							ReqRoomList(client, userID, data)
						},
					},
					PushButton{
						Text: "join",
						OnClicked: func() {
							if cmd, err := RunRoomJoinDialog(mw); err != nil {
								log.Print(err)
							} else if cmd == walk.DlgCmdOK {
								log.Println("dlg msg : " + inputString)
								num, err := strconv.Atoi(inputString)
								lib.CheckError(err)
								roomID := int64(num)
								ReqJoin(client, userID, data, roomID)
							}
						},
					},
					PushButton{
						Text: "action1",
						OnClicked: func() {
							ReqAction1(client, userID, data)
						},
					},
					PushButton{
						Text: "quit",
						OnClicked: func() {
							log.Println("quit user id : ", userID)
							ReqQuit(client, userID, data)
							os.Exit(3)
						},
					},
				},
			},
		},
	}.Create()); err != nil {
		log.Fatal(err)
	}

	lv, err := NewLogView(mw)
	if err != nil {
		log.Fatal(err)
	}

	//logFile, err := os.OpenFile("log.txt", os.O_WRONLY, 0666)
	log.SetOutput(lv)

	go func() {
		data := make([]byte, 4096)

		for {
			log.Println("wait for read")
			n, err := client.Read(data)
			if err != nil {
				log.Println("Fail Stream read, err : ", err)
				break
			}

			rawData := data[:n]
			message := fbMessage.GetRootAsMessage(rawData, 0)
			messageType := message.BodyType()
			handler, ok := msgHandler[messageType]
			log.Println("recv message type : ", messageType)
			if ok {
				ret := handler(message) // calling proper handler function
				if !ret {
					log.Println("Fail handler process", handler)
				}
			} else {
				log.Println("Fail no function defined for type", handler)
				break
			}
		}
	}()

	mw.Run()
}

func ReqLogin(c net.Conn, userUID int64, data []byte) {

	// create flatbuffers
	builder := flatbuffers.NewBuilder(0)

	// create request login message (in message union)
	fbMessage.ReqLoginStart(builder)
	fbMessage.ReqLoginAddUserID(builder, userUID)
	reqBody := fbMessage.ReqLoginEnd(builder)

	// create message
	fbMessage.MessageStart(builder)
	fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyReqLogin)
	fbMessage.MessageAddBody(builder, reqBody)
	msg := fbMessage.MessageEnd(builder)

	// Call `Finish()` to instruct the builder that this message is complete.
	builder.Finish(msg)

	sendData := builder.FinishedBytes()
	_, err := c.Write(sendData)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("ReqLogin client send : %x\n", sendData)
}

func ResLogin(message *fbMessage.Message) bool {

	unionTable := new(flatbuffers.Table)
	if !message.Body(unionTable) {
		return false
	}
	res := new(fbMessage.ResLogin)
	res.Init(unionTable.Bytes, unionTable.Pos)

	log.Println("ResLogin server return : user id : " + lib.Itoa64(res.UserID()))
	log.Println("ResLogin server return : result code : " + lib.Itoa32(res.ResultCode()))

	return true
}

func ReqRoomList(c net.Conn, userUID int64, data []byte) {
	// create flatbuffers
	builder := flatbuffers.NewBuilder(0)

	// create request login message (in message union)
	fbMessage.ReqRoomListStart(builder)
	fbMessage.ReqRoomListAddUserID(builder, userUID)
	reqBody := fbMessage.ReqRoomListEnd(builder)

	// create message
	fbMessage.MessageStart(builder)
	fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyReqRoomList)
	fbMessage.MessageAddBody(builder, reqBody)
	msg := fbMessage.MessageEnd(builder)

	// Call `Finish()` to instruct the builder that this message is complete.
	builder.Finish(msg)

	sendData := builder.FinishedBytes()
	_, err := c.Write(sendData)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("ReqRoomList client send : %x\n", sendData)
}

func ResRoomList(message *fbMessage.Message) bool {

	unionTable := new(flatbuffers.Table)
	if !message.Body(unionTable) {
		return false
	}
	res := new(fbMessage.ResRoomList)
	res.Init(unionTable.Bytes, unionTable.Pos)

	log.Println("ResLogin server return : user id : ", res.UserID())
	log.Println("ResLogin server return : result code : ", res.ResultCode())

	rCount := res.RoomIDsLength()
	for i := 0; i < rCount; i++ {
		log.Println("ResRoomList server return : room id : ", res.RoomIDs(i))
	}
	return true
}

func ReqCreate(c net.Conn, userUID int64, data []byte) {
	// create flatbuffers
	builder := flatbuffers.NewBuilder(0)

	// create request login message (in message union)
	fbMessage.ReqCreateStart(builder)
	fbMessage.ReqCreateAddUserID(builder, userUID)
	reqBody := fbMessage.ReqCreateEnd(builder)

	// create message
	fbMessage.MessageStart(builder)
	fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyReqCreate)
	fbMessage.MessageAddBody(builder, reqBody)
	msg := fbMessage.MessageEnd(builder)

	// Call `Finish()` to instruct the builder that this message is complete.
	builder.Finish(msg)

	sendData := builder.FinishedBytes()
	_, err := c.Write(sendData)
	if err != nil {
		log.Println(err)
		return
	}

	log.Printf("client send : %x\n", sendData)
}

func ResCreate(message *fbMessage.Message) bool {

	unionTable := new(flatbuffers.Table)
	if !message.Body(unionTable) {
		return false
	}
	res := new(fbMessage.ResCreate)
	res.Init(unionTable.Bytes, unionTable.Pos)

	log.Println("ResCreate server return : user id : ", res.UserID())
	log.Println("ResCreate server return : result code : ", res.ResultCode())
	log.Println("ResCreate server return : room id : ", res.RoomID())

	return true
}

func ReqJoin(c net.Conn, userUID int64, data []byte, roomID int64) {
	// Join flatbuffers
	builder := flatbuffers.NewBuilder(0)

	// Join request login message (in message union)
	fbMessage.ReqJoinStart(builder)
	fbMessage.ReqJoinAddUserID(builder, userUID)
	fbMessage.ReqJoinAddRoomID(builder, roomID)
	reqBody := fbMessage.ReqJoinEnd(builder)

	// Join message
	fbMessage.MessageStart(builder)
	fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyReqJoin)
	fbMessage.MessageAddBody(builder, reqBody)
	msg := fbMessage.MessageEnd(builder)

	// Call `Finish()` to instruct the builder that this message is complete.
	builder.Finish(msg)

	sendData := builder.FinishedBytes()
	_, err := c.Write(sendData)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("ReqJoin client send : %x\n", sendData)
}

func ResJoin(message *fbMessage.Message) bool {
	unionTable := new(flatbuffers.Table)
	if !message.Body(unionTable) {
		return false
	}
	res := new(fbMessage.ResJoin)
	res.Init(unionTable.Bytes, unionTable.Pos)

	log.Println("ResLogin server return : user id : ", res.UserID())
	log.Println("ResLogin server return : result code : ", res.ResultCode())

	memberCount := res.MembersLength()
	for i := 0; i < memberCount; i++ {
		log.Println("ResLogin server return : member id : ", res.Members(i))
	}
	return true
}

func ReqAction1(c net.Conn, userUID int64, data []byte) {
	// Join flatbuffers
	builder := flatbuffers.NewBuilder(0)

	// Join request login message (in message union)
	fbMessage.ReqAction1Start(builder)
	fbMessage.ReqAction1AddUserID(builder, userUID)
	reqBody := fbMessage.ReqAction1End(builder)

	// Join message
	fbMessage.MessageStart(builder)
	fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyReqAction1)
	fbMessage.MessageAddBody(builder, reqBody)
	msg := fbMessage.MessageEnd(builder)

	// Call `Finish()` to instruct the builder that this message is complete.
	builder.Finish(msg)

	sendData := builder.FinishedBytes()
	_, err := c.Write(sendData)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("ReqAction1 client send : %x\n", sendData)
}

func ResAction1(message *fbMessage.Message) bool {

	unionTable := new(flatbuffers.Table)
	if !message.Body(unionTable) {
		return false
	}
	res := new(fbMessage.ResAction1)
	res.Init(unionTable.Bytes, unionTable.Pos)

	log.Println("ResAction1 server return : user id : ", res.UserID())
	return true
}

func ReqQuit(c net.Conn, userUID int64, data []byte) {
	// Join flatbuffers
	builder := flatbuffers.NewBuilder(0)

	// Join request login message (in message union)
	fbMessage.ReqQuitStart(builder)
	fbMessage.ReqQuitAddUserID(builder, userUID)
	reqBody := fbMessage.ReqQuitEnd(builder)

	// Join message
	fbMessage.MessageStart(builder)
	fbMessage.MessageAddBodyType(builder, fbMessage.MessageBodyReqQuit)
	fbMessage.MessageAddBody(builder, reqBody)
	msg := fbMessage.MessageEnd(builder)

	// Call `Finish()` to instruct the builder that this message is complete.
	builder.Finish(msg)

	sendData := builder.FinishedBytes()
	_, err := c.Write(sendData)
	if err != nil {
		fmt.Println(err)
		return
	}

	log.Printf("ReqQuit client send : %x\n", sendData)
}

func ResQuit(message *fbMessage.Message) bool {
	unionTable := new(flatbuffers.Table)
	if !message.Body(unionTable) {
		return false
	}
	res := new(fbMessage.ResQuit)
	res.Init(unionTable.Bytes, unionTable.Pos)

	log.Println("ResQuit server return : user id : ", res.UserID())
	log.Println("ResQuit server return : result code : ", res.ResultCode())
	return true
}

func NotifyJoinHandler(message *fbMessage.Message) bool {
	unionTable := new(flatbuffers.Table)
	if !message.Body(unionTable) {
		return false
	}
	res := new(fbMessage.NotifyJoin)
	res.Init(unionTable.Bytes, unionTable.Pos)

	log.Println("NotifyJoin user id : ", res.UserID())

	return true
}

func NotifyAction1Handler(message *fbMessage.Message) bool {
	unionTable := new(flatbuffers.Table)
	if !message.Body(unionTable) {
		return false
	}
	res := new(fbMessage.NotifyAction1)
	res.Init(unionTable.Bytes, unionTable.Pos)

	log.Println("NotifyAction1Handler sender id : ", res.UserID())
	return true
}

func NotifyQuitHandler(message *fbMessage.Message) bool {
	unionTable := new(flatbuffers.Table)
	if !message.Body(unionTable) {
		return false
	}
	res := new(fbMessage.NotifyQuit)
	res.Init(unionTable.Bytes, unionTable.Pos)

	log.Println("NotifyQuit user id : ", res.UserID())

	return true
}

func RunUserIdDialog(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var inDlg *walk.LineEdit

	return Dialog{
		AssignTo:      &dlg,
		Title:         "input User ID",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize: Size{200, 100},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "User ID:",
					},
					LineEdit{
						AssignTo:&inDlg,
						Text: "",
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							inputString = inDlg.Text()
							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(owner)
}

func RunRoomJoinDialog(owner walk.Form) (int, error) {
	var dlg *walk.Dialog
	var acceptPB, cancelPB *walk.PushButton
	var inDlg *walk.LineEdit

	return Dialog{
		AssignTo:      &dlg,
		Title:         "input Room ID",
		DefaultButton: &acceptPB,
		CancelButton:  &cancelPB,
		MinSize: Size{200, 100},
		Layout:  VBox{},
		Children: []Widget{
			Composite{
				Layout: Grid{Columns: 2},
				Children: []Widget{
					Label{
						Text: "room id:",
					},
					LineEdit{
						AssignTo:&inDlg,
						Text: "",
					},
				},
			},
			Composite{
				Layout: HBox{},
				Children: []Widget{
					HSpacer{},
					PushButton{
						AssignTo: &acceptPB,
						Text:     "OK",
						OnClicked: func() {
							inputString = inDlg.Text()
							dlg.Accept()
						},
					},
					PushButton{
						AssignTo:  &cancelPB,
						Text:      "Cancel",
						OnClicked: func() {
							dlg.Cancel()
						},
					},
				},
			},
		},
	}.Run(owner)
}
