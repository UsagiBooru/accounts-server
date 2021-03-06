package server

import (
	_ "github.com/bshuster-repo/logrus-logstash-hook"
	"github.com/sirupsen/logrus"
)

// Log - log client.
var Log *logrus.Logger

func init() {
	logger := logrus.New()
	/*
		TODO: Test logstash later.
		logstashHook, err := logrustash.NewHook("tcp", "logstash:5959", "from_logstash")
		if err != nil {
			panic(err)
		}
		logger.Hooks.Add(logstashHook)
	*/
	logger.SetLevel(logrus.DebugLevel)
	Log = logger
}

// Debug - shorthand Log.Debug
func Debug(msg string) {
	Log.Debug(msg)
}

// Info - shorthand Log.Info
func Info(msg string) {
	Log.Info(msg)
}

// Warn - shorthand Log.Warn
func Warn(msg string) {
	Log.Warn(msg)
}

// Error - shorthand Log.Error
func Error(msg string) {
	Log.Error(msg)
}

// Fatal - shorthand Log.Fatal
func Fatal(msg string) {
	Log.Fatal(msg)
}
