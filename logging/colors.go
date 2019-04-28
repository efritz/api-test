package logging

import (
	"fmt"

	"github.com/mgutz/ansi"
)

var levelColors = map[LogLevel]string{
	LevelDebug: ansi.Cyan,
	LevelInfo:  ansi.Green,
	LevelWarn:  ansi.Yellow,
	LevelError: ansi.ColorCode("red+b"),
}

func Colorize(message string, level LogLevel) string {
	return fmt.Sprintf(
		"%s%s%s",
		levelColors[level],
		message,
		ansi.Reset,
	)
}
