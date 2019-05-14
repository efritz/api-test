package logging

import (
	"fmt"

	"github.com/mgutz/ansi"
)

type (
	Logger interface {
		Raw(message string)
		Log(format string, args ...interface{})
		Info(format string, args ...interface{})
		Warn(format string, args ...interface{})
		Error(format string, args ...interface{})
		Colorize(message string, color Color) string
	}

	logger struct {
		colorize bool
	}
)

func NewLogger(colorize bool) Logger {
	return &logger{
		colorize: colorize,
	}
}

func (l *logger) Raw(message string) {
	fmt.Print(message)
}

func (l *logger) Log(format string, args ...interface{})   { l.log(ColorNone, format, args...) }
func (l *logger) Info(format string, args ...interface{})  { l.log(ColorInfo, format, args...) }
func (l *logger) Warn(format string, args ...interface{})  { l.log(ColorWarn, format, args...) }
func (l *logger) Error(format string, args ...interface{}) { l.log(ColorError, format, args...) }

func (l *logger) Colorize(message string, color Color) string {
	if !l.colorize {
		return message
	}

	return fmt.Sprintf(
		"%s%s%s",
		colors[color],
		message,
		ansi.Reset,
	)
}

func (l *logger) log(color Color, format string, args ...interface{}) {
	fmt.Println(l.Colorize(fmt.Sprintf(format, args...), color))
}
