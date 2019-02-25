package logging

import "github.com/mgutz/ansi"

var levelColors = map[LogLevel]string{
	LevelDebug: ansi.Cyan,
	LevelInfo:  ansi.Green,
	LevelWarn:  ansi.Yellow,
	LevelError: ansi.ColorCode("red+b"),
}
