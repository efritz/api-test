package logging

import "io"

type (
	logWriter struct {
		logger Logger
	}
)

func Writer(logger Logger) io.Writer {
	return &logWriter{logger: logger}
}

func (w *logWriter) Write(p []byte) (int, error) {
	w.logger.Raw(string(p))
	return len(p), nil
}
