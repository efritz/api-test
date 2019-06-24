package logging

import (
	"hash/fnv"

	"github.com/mgutz/ansi"
)

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

var allColors = []string{
	ansi.Red,
	ansi.Yellow,
	ansi.Blue,
	ansi.Magenta,
	ansi.White,
	ansi.LightBlack,
	ansi.LightRed,
	ansi.LightGreen,
	ansi.LightYellow,
	ansi.LightBlue,
	ansi.LightMagenta,
	ansi.LightCyan,
	ansi.LightWhite,
}

func chooseColor(text string) string {
	hash := fnv.New32a()
	hash.Write([]byte(text))

	return allColors[int(hash.Sum32())%len(allColors)]
}
