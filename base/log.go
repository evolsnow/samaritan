package base

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
