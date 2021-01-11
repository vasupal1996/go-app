package logger

import (
	"fmt"
	"go-app/server/config"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/pkgerrors"
)

// NewLogger returns logger based on server config
func NewLogger(c *config.LoggerConfig, kw *KafkaLogWriter) *zerolog.Logger {
	var writers []io.Writer

	if c.EnableKafkaLog {
		wr := diode.NewWriter(kw, 1000, 10*time.Millisecond, func(missed int) {
			fmt.Printf("Logger Dropped %d messages", missed)
		})
		writers = append(writers, wr)
	}
	if c.EnableConsoleLog {
		writers = append(writers, zerolog.ConsoleWriter{Out: os.Stderr})
	}
	mw := io.MultiWriter(writers...)
	zlog := zerolog.New(mw).With().Timestamp().Stack().Caller().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return &zlog
}
