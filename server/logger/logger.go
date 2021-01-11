package logger

import (
	"go-app/server/config"
	"io"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

// NewLogger returns logger based on server config
func NewLogger(c *config.LoggerConfig, kw *KafkaLogWriter) *zerolog.Logger {
	var writers []io.Writer

	if c.EnableKafkaLog {
		writers = append(writers, kw)
	}
	if c.EnableConsoleLog {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stdout})
	}

	mw := io.MultiWriter(writers...)
	zlog := zerolog.New(mw).With().Timestamp().Stack().Caller().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return &zlog
}
