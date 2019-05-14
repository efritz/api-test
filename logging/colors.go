package logging

import "github.com/mgutz/ansi"

type Color int

const (
	ColorNone Color = iota
	ColorInfo
	ColorWarn
	ColorError
)

var colors = map[Color]string{
	ColorNone:  "",
	ColorInfo:  ansi.Green,
	ColorWarn:  ansi.Yellow,
	ColorError: ansi.ColorCode("red+b"),
}
