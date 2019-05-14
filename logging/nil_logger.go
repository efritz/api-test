package logging

type nilLogger struct{}

var NilLogger = &nilLogger{}

func (l *nilLogger) Raw(message string)                       {}
func (l *nilLogger) Log(format string, args ...interface{})   {}
func (l *nilLogger) Info(format string, args ...interface{})  {}
func (l *nilLogger) Warn(format string, args ...interface{})  {}
func (l *nilLogger) Error(format string, args ...interface{}) {}

func (l *nilLogger) Colorize(message string, color Color) string {
	return message
}
