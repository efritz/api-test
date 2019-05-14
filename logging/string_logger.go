package logging

import "fmt"

type StringLogger struct {
	content string
}

func NewStringLogger() *StringLogger {
	return &StringLogger{}
}

func (l *StringLogger) Raw(message string) {
	l.content += message
}

func (l *StringLogger) Log(format string, args ...interface{}) {
	l.Raw(fmt.Sprintf(format+"\n", args...))
}

func (l *StringLogger) Info(format string, args ...interface{}) {
	l.Raw(fmt.Sprintf(format+"\n", args...))
}

func (l *StringLogger) Warn(format string, args ...interface{}) {
	l.Raw(fmt.Sprintf(format+"\n", args...))
}

func (l *StringLogger) Error(format string, args ...interface{}) {
	l.Raw(fmt.Sprintf(format+"\n", args...))
}

func (l *StringLogger) Colorize(message string, color Color) string {
	return message
}

func (l *StringLogger) String() string {
	return l.content
}
