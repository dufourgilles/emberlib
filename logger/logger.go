package logger

import (
	"fmt"

	"github.com/dufourgilles/emberlib/errors"
)


type LogLevel uint8
const (

	ErrorLevel LogLevel = 3
	WarnLevel LogLevel = iota
	InfoLevel LogLevel = iota
	DebugLevel LogLevel = iota
)

type Logger interface {
	Error(errors.Error)
	Warn(format string, a ...interface{})
	Info(format string, a ...interface{})
	Debug(format string, a ...interface{})
	Debugln(a ...interface{})
	SetLogLevel(level LogLevel)
	GetLogLevel() LogLevel
}

type ConsoleLogger struct {
	logLevel LogLevel
}

type NullLogger struct {
	logLevel LogLevel
}

func NewNullLogger() Logger {
	return &NullLogger{}
}

func NewConsoleLogger(logLevel LogLevel) Logger {
	return &ConsoleLogger{logLevel: logLevel}
}

func (logger *ConsoleLogger)Error(e errors.Error) {
	if logger.logLevel >= ErrorLevel {
		fmt.Println(e)
	}
}

func (logger *ConsoleLogger)Warn(format string, a ...interface{}) {
	if logger.logLevel >= ErrorLevel {
		fmt.Printf(format, a...)
	}
}

func (logger *ConsoleLogger)Info(format string, a ...interface{}) {
	if logger.logLevel >= ErrorLevel {
		fmt.Printf(format, a...)
	}
}

func (logger *ConsoleLogger)Debug(format string, a ...interface{}) {
	if logger.logLevel >= ErrorLevel {
		fmt.Printf(format, a...)
	}
}

func (logger *ConsoleLogger)Debugln(a ...interface{}) {
	if logger.logLevel >= ErrorLevel {
		fmt.Println(a...)
	}
}

func (logger *ConsoleLogger)SetLogLevel(level LogLevel) {
	logger.logLevel = level
}

func (logger *ConsoleLogger)GetLogLevel() LogLevel {
	return logger.logLevel
}




func (logger *NullLogger)Error(e errors.Error) {
}
func (logger *NullLogger)Warn(format string, a ...interface{}) {
}
func (logger *NullLogger)Info(format string, a ...interface{}) {
}
func (logger *NullLogger)Debug(format string, a ...interface{}) {
}
func (logger *NullLogger)Debugln(a ...interface{}) {
}
func (logger *NullLogger)SetLogLevel(level LogLevel) {
	logger.logLevel = level
}

func (logger *NullLogger)GetLogLevel() LogLevel {
	return logger.logLevel
}

