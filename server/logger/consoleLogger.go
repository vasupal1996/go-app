package logger

import (
	"io"
	"os"

	"github.com/rs/zerolog"
)

// ConsoleWriter is a type io.Writer
type ConsoleWriter io.Writer

// NewStandardConsoleWriter returns standard output as writer
func NewStandardConsoleWriter() ConsoleWriter {
	cw := os.Stdout
	return cw
}

// NewZeroLogConsoleWriter returns new instance zerolog console logger
func NewZeroLogConsoleWriter(cw io.Writer) io.Writer {
	return zerolog.ConsoleWriter{Out: cw}
}
