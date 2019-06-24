package logging

import (
	"fmt"
	"os"
)

func EmergencyLog(format string, args ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, args...))
}
