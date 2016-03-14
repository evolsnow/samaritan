package log

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
)

var Logger = logrus.New()

func SetLogLevel(debug bool) {
	if debug {
		Logger.Level = logrus.DebugLevel
		Logger.Debug("debug mode")
	} else {
		Logger.Level = logrus.InfoLevel
	}
}

func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Println(args ...interface{}) {
	Logger.Println(args...)
}

// DebugJson pretty print the request and response
// Should add no more than one description
func DebugJson(args ...interface{}) {
	if len(args) == 1 {
		//common case
		Logger.Debug("\n", marshalWithIndent(args[0]))
	} else if len(args) == 2 {
		Logger.Debug(args[0].(string)+"\n", marshalWithIndent(args[1]))
	} else {
		//should use Debug instead
		Logger.Debug(args...)
	}

}

func marshalWithIndent(s interface{}) interface{} {
	b, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return s
	} else {
		return string(b)
	}
}
