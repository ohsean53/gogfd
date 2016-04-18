package lib

import (
	"encoding/binary"
	"fmt"
	"github.com/artyom/scribe"
	"github.com/artyom/thrift"
	"math/rand"
	"os"
	"runtime"
	"strconv"
)

func Itoa64(i int64) string {
	return strconv.FormatInt(i, 10)
}

func Itoa32(i int32) string {
	return strconv.Itoa(int(i))
}

func Itoa(i int) string {
	return strconv.Itoa(i)
}

func Log(a ...interface{}) {
	fmt.Println(a...)
}

func Logf(format string, a ...interface{}) {
	fmt.Printf(format, a...)
}

func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
	}
}

// http://stackoverflow.com/questions/16888357/convert-an-integer-to-a-byte-array
func ReadInt32(data []byte) (ret int32) {
	ret = int32(binary.BigEndian.Uint32(data)) // fastest convert method, do not use "binary.Read"
	return
}

// After benchmarking the "encoding/binary" way, it takes almost 4 times longer than int -> string -> byte
func WriteInt32(n int32) (buf []byte) {
	buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(n)) // fastest convert method, do not use "binary.Write"
	return
}

func RandInt64(min int64, max int64) int64 {
	return min + rand.Int63n(max - min)
}

func RandInt32(min int32, max int32) int32 {
	return min + rand.Int31n(max - min)
}

func WriteScribe(category string, message string) {

	// currently available on linux platform
	if runtime.GOOS != "linux" {
		Log(category + " : " + message)
		return
	}
	entry := scribe.NewLogEntry()
	entry.Category = category
	entry.Message = message
	messages := []*scribe.LogEntry{entry}
	socket, err := thrift.NewTSocket("localhost:1463")
	CheckError(err)

	transport := thrift.NewTFramedTransport(socket)
	protocol := thrift.NewTBinaryProtocol(transport, false, false)
	client := scribe.NewScribeClientProtocol(transport, protocol, protocol)

	transport.Open()
	result, err := client.Log(messages)
	CheckError(err)
	transport.Close()
	Log(result.String())
}

func Int64SliceToString(set []int64) (str string) {
	str += "[";
	for _, value := range set {
		str += "," + Itoa64(value)
	}
	str += "]";
	return str
}