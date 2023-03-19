package logger

import "log"

type LogLevel string

const (
	LogDebug  LogLevel = "Debug"
	LogInfo   LogLevel = "Info"
	LogSilent LogLevel = "Silent"
)

type Log struct {
	Level LogLevel
}

func NewLogger(level LogLevel) *Log {
	return &Log{
		Level: level,
	}
}

func (l *Log) Debug(args ...interface{}) {
	if l.Level != "Debug" {
		return
	}
	log.Println("[DEBUG] ", args)
}

func (l *Log) Info(args ...interface{}) {
	log.Println("[Info] ", args)
}

func (l *Log) Error(args ...interface{}) {
	log.Println("[ERROR] ", args)
}

func (l *Log) Warn(args ...interface{}) {
	log.Println("[WARNING] ", args)
}
