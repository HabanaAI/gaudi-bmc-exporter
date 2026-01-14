// Package logger provides convient method to work with our selected logger.
package logger

import (
	"os"

	joonix "github.com/joonix/log"
	"github.com/sirupsen/logrus"
)

const envLogLevel = "LOG_LEVEL"
const envLogFormat = "LOG_FORMAT"
const fluentdFormat = "Fluentd"
const logrusFormat = "Logrus"
const defaultLogLevel = logrus.InfoLevel

func New() *logrus.Logger {
	return &logrus.Logger{
		Out:       os.Stdout,
		Formatter: getLogFormat(),
		Level:     getLogLevel(),
	}
}

func getLogLevel() logrus.Level {
	levelString, exists := os.LookupEnv(envLogLevel)
	if !exists {
		return defaultLogLevel
	}

	level, err := logrus.ParseLevel(levelString)
	if err != nil {
		logrus.Errorf("error parsing %s: %v", envLogLevel, err)
		return defaultLogLevel
	}

	return level
}

func getLogFormat() logrus.Formatter {
	formatString, exists := os.LookupEnv(envLogFormat)
	if !exists || formatString == logrusFormat {
		return &logrus.JSONFormatter{}
	} else if formatString == fluentdFormat {
		return &joonix.Formatter{}
	} else {
		logrus.Errorf("unknown %s: %v", envLogFormat, formatString)
		return &logrus.JSONFormatter{}
	}
}
