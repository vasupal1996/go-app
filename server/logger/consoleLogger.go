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

// NewZeroLogConsoleLogger returns new instance zerolog console logger
func NewZeroLogConsoleLogger(cw io.Writer) io.Writer {
	return zerolog.ConsoleWriter{Out: cw}
}
