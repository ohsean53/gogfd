package logger

import (
	"github.com/artyom/scribe"
	"github.com/artyom/thrift"
	"gogfd/config"
	"gogfd/lib"
	"fmt"
	"runtime"
	"os"
)

type LogLevel int

const (
	DEBUG LogLevel = 0
	INFO LogLevel = 1
	NOTIFY LogLevel = 2
	WARNING LogLevel = 3
	ERROR LogLevel = 4
	CRITICAL LogLevel = 5
)

func getLogTypeName(lv LogLevel) string {
	switch lv {
	case DEBUG :
		return "DEBU"
	case INFO :
		return "INFO"
	case NOTIFY :
		return "NOTI"
	case WARNING :
		return "WARN"
	case ERROR :
		return "ERRO"
	case CRITICAL :
		return "CRIT"
	}
	return "NDEF"
}

func checkDebugMode(lv LogLevel) bool {
	if lv == DEBUG {
		if config.DEBUG == false {
			return false
		}
	}
	return true
}

func Log(lv LogLevel, a ...interface{}) {

	if checkDebugMode(lv) == false {
		return
	}
	pc := make([]uintptr, 10)  // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	fmt.Fprintf(os.Stderr, "%s %s ▶ %s ", lib.GetDateTime(), f.Name(), getLogTypeName(lv))
	fmt.Fprint(os.Stderr, a...)
	fmt.Fprintln(os.Stderr)
}

func Logf(lv LogLevel, format string, a ...interface{}) {
	if checkDebugMode(lv) == false {
		return
	}
	pc := make([]uintptr, 10)  // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, line := f.FileLine(pc[0])
	fmt.Fprintf(os.Stderr, "%s %s:%d %s ▶ %s ", lib.GetDateTime(), file, line, f.Name(), getLogTypeName(lv))
	fmt.Fprintf(os.Stderr, format, a...)
	fmt.Fprintln(os.Stderr)
}

func WriteScribe(category string, message string) {

	// currently available on linux platform
	if runtime.GOOS != "linux" {
		Log(DEBUG, category + " : " + message)
		return
	}
	entry := scribe.NewLogEntry()
	entry.Category = category
	entry.Message = message
	messages := []*scribe.LogEntry{entry}
	socket, err := thrift.NewTSocket(config.SERVER_IP + ":1463")
	lib.CheckError(err)

	transport := thrift.NewTFramedTransport(socket)
	protocol := thrift.NewTBinaryProtocol(transport, false, false)
	client := scribe.NewScribeClientProtocol(transport, protocol, protocol)

	transport.Open()
	result, err := client.Log(messages)
	lib.CheckError(err)
	transport.Close()
	Log(DEBUG, result.String())
}
