package logging

import (
	"fmt"
)

type (
	Logger interface {
		Debug(format string, args ...interface{})
		Info(format string, args ...interface{})
		Warn(format string, args ...interface{})
		Error(format string, args ...interface{})
	}

	logger struct {
		colorize bool
		verbose  bool
	}

	nilLogger struct{}
)

var NilLogger = &nilLogger{}

func NewLogger(colorize, verbose bool) Logger {
	return &logger{
		colorize: colorize,
		verbose:  verbose,
	}
}

func (l *logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, fmt.Sprintf(format, args...))
}

func (l *logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, fmt.Sprintf(format, args...))
}

func (l *logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, fmt.Sprintf(format, args...))
}

func (l *logger) Error(format string, args ...interface{}) {
	l.log(LevelError, fmt.Sprintf(format, args...))
}

func (l *logger) log(level LogLevel, message string) {
	if level == LevelDebug && !l.verbose {
		return
	}

	if l.colorize {
		message = Colorize(message, level)
	}

	fmt.Println(message)
}

func (l *nilLogger) Debug(string, ...interface{}) {}
func (l *nilLogger) Info(string, ...interface{})  {}
func (l *nilLogger) Warn(string, ...interface{})  {}
func (l *nilLogger) Error(string, ...interface{}) {}
