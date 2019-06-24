package logging

import (
	"fmt"
	"sync"

	"github.com/mgutz/ansi"
)

type (
	Logger interface {
		Colorized() bool
		Log(prefix *Prefix, format string, args ...interface{})
		Colorize(color Color, format string, args ...interface{}) string
	}

	logger struct {
		colorize   bool
		mutex      sync.Mutex
		touched    bool
		lastPrefix *Prefix
	}
)

func NewLogger(colorize bool) Logger {
	return &logger{
		colorize: colorize,
	}
}

func (l *logger) Colorized() bool {
	return l.colorize
}

func (l *logger) Log(prefix *Prefix, format string, args ...interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	fmt.Printf(formatPrefix(prefix, l.lastPrefix, l.colorize, !l.touched, format, args...))
	l.touched = true
	l.lastPrefix = prefix
}

func (l *logger) Colorize(color Color, format string, args ...interface{}) string {
	if !l.colorize {
		return fmt.Sprintf(format, args...)
	}

	return fmt.Sprintf(
		"%s%s%s",
		colors[color],
		fmt.Sprintf(format, args...),
		ansi.Reset,
	)
}
