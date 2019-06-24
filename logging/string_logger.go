package logging

import (
	"fmt"
	"sync"
)

type StringLogger struct {
	content    string
	mutex      sync.Mutex
	lastPrefix *Prefix
	touched    bool
}

func NewStringLogger() *StringLogger {
	return &StringLogger{}
}

func (l *StringLogger) Colorized() bool {
	return false
}

func (l *StringLogger) Log(prefix *Prefix, format string, args ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.content += formatPrefix(prefix, l.lastPrefix, false, !l.touched, format, args...)
	l.touched = true
	l.lastPrefix = prefix
}

func (l *StringLogger) Colorize(color Color, format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}

func (l *StringLogger) String() string {
	return l.content
}
