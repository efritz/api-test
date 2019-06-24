package logging

import "fmt"

type nilLogger struct{}

var NilLogger = &nilLogger{}

func (l *nilLogger) Colorized() bool {
	return false
}

func (l *nilLogger) Log(prefix *Prefix, format string, args ...interface{}) {
	// no-op
}

func (l *nilLogger) Colorize(color Color, format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
