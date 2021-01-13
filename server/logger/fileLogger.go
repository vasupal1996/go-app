package logger

import (
	"fmt"
	"go-app/server/config"
	"io"
	"os"

	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

// FileWriter is a type io.Writer
type FileWriter io.Writer

// NewFileWriter returns file writer
func NewFileWriter(fn, path string, c *config.FileLoggerConfig) FileWriter {
	if err := os.MkdirAll(path, 0744); err != nil {
		log.Error().Err(err).Str("path", path).Msg("can't create log directory")
		return nil
	}
	return &lumberjack.Logger{
		Filename:   fmt.Sprintf("%s/%s.log", path, fn),
		MaxBackups: c.MaxBackupsFile, // files
		MaxSize:    c.MaxSize,        // megabytes
		MaxAge:     c.MaxAge,         // days
		Compress:   c.Compress,
	}
}
