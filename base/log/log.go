package log

import (
	"github.com/Sirupsen/logrus"
)

var Logger = logrus.New()

func SetLogLevel(debug bool) {
	if debug {
		Logger.Level = logrus.DebugLevel
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
