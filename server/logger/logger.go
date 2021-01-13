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
func NewLogger(c *config.LoggerConfig, kw *KafkaLogWriter, cw, fw io.Writer) *zerolog.Logger {
	var writers []io.Writer

	// Setting up kafka writer if True.
	if c.EnableKafkaLogger {
		if kw == nil {
			fmt.Println("failed to create kafka lagger: kafka writer cannot be nil")
			os.Exit(1)
		}
		wr := diode.NewWriter(kw, 1000, 10*time.Millisecond, func(missed int) {
			fmt.Printf("Logger Dropped %d messages", missed)
		})
		writers = append(writers, wr)
	}

	// Setting up console writer if True.
	if c.EnableConsoleLogger {
		writers = append(writers, zerolog.ConsoleWriter{Out: cw})
	}

	// Setting up file writer is True.
	if c.EnableFileLogger {
		if kw == nil {
			fmt.Println("failed to create file lagger: file writer cannot be nil")
			os.Exit(1)
		}
		wr := diode.NewWriter(fw, 1000, 10*time.Millisecond, func(missed int) {
			fmt.Printf("Logger Dropped %d messages", missed)
		})
		writers = append(writers, wr)
	}

	mw := io.MultiWriter(writers...)
	zlog := zerolog.New(mw).With().Timestamp().Stack().Caller().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return &zlog
}
