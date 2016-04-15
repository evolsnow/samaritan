package log

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
)

var Logger = logrus.New()

// SetLogLevel sets log level with the given level
func SetLogLevel(debug bool) {
	if debug {
		Logger.Level = logrus.DebugLevel
		Logger.Debug("debug mode")
	} else {
		Logger.Level = logrus.InfoLevel
	}
}

// Panic is Logger.Panic exactly
func Panic(args ...interface{}) {
	Logger.Panic(args...)
}

// Fatal is Logger.Fatal exactly
func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

// Error is Logger.Error exactly
func Error(args ...interface{}) {
	Logger.Error(args...)
}

// Warn is Logger.Warn exactly
func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

// Info is Logger.Info exactly
func Info(args ...interface{}) {
	Logger.Info(args...)
}

// Debug is Logger.Debug exactly
func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

// Println is Logger.Println exactly
func Println(args ...interface{}) {
	Logger.Println(args...)
}

// DebugJson pretty print the request and response
// Should add no more than one description
func DebugJson(args ...interface{}) error {
	if len(args) == 1 {
		//common case
		Logger.Debug("\n", marshalWithIndent(args[0]))
	} else if len(args) == 2 {
		if _, ok := args[0].(string); !ok {
			return fmt.Errorf("first arg should be string")
		}
		Logger.Debug(args[0].(string)+"\n", marshalWithIndent(args[1]))
	} else {
		//should use Debug instead
		Logger.Debug(args...)
	}
	return nil
}

//marshal with four whitespace
func marshalWithIndent(s interface{}) interface{} {
	b, err := json.MarshalIndent(s, "", "    ")
	if err != nil {
		return s
	} else {
		return string(b)
	}
}
