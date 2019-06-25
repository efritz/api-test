package logging

import (
	"fmt"
	"strings"

	"github.com/mgutz/ansi"
)

type Prefix struct {
	scenarioName string
	testName     string
}

func NewPrefix(scenarioName, testName string) *Prefix {
	return &Prefix{
		scenarioName: scenarioName,
		testName:     testName,
	}
}

func formatPrefix(prefix, lastPrefix *Prefix, colorize, first bool, format string, args ...interface{}) string {
	message := fmt.Sprintf(format, args...)
	if prefix == nil {
		return fmt.Sprintf("%s\n", strings.TrimSpace(message))
	}

	return fmt.Sprintf(
		"%s%s\n",
		serializePrefix(prefix, lastPrefix, colorize, first),
		indent(message),
	)
}

func serializePrefix(prefix, lastPrefix *Prefix, colorize, first bool) string {
	if prefix == lastPrefix {
		return ""
	}

	if first {
		return colorizePrefix(prefix, colorize)
	}

	return "\n" + colorizePrefix(prefix, colorize)
}

func colorizePrefix(prefix *Prefix, colorize bool) string {
	if !colorize {
		return fmt.Sprintf("%s/%s", prefix.scenarioName, prefix.testName)
	}

	color := chooseColor(prefix.scenarioName)

	return fmt.Sprintf(
		"%s%s%s/%s%s%s: \n",
		color,
		prefix.scenarioName,
		ansi.Reset,
		color,
		prefix.testName,
		ansi.Reset,
	)
}

func indent(message string) string {
	lines := strings.Split(message, "\n")
	for i, line := range lines {
		lines[i] = fmt.Sprintf("  %s", line)
	}

	return strings.Join(lines, "\n")
}
