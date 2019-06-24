package logging

type VerbosityLevel int

const (
	VerbosityLevelNone VerbosityLevel = iota
	VerbosityLevelProgress
	VerbosityLevelRequestResponse
)

func ToVerbosityLevel(level int) VerbosityLevel {
	if level == 0 {
		return VerbosityLevelNone
	}

	if level == 1 {
		return VerbosityLevelProgress
	}

	return VerbosityLevelRequestResponse
}
