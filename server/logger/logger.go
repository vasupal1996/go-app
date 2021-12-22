package logger

import (
	"fmt"
	"io"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/diode"
	"github.com/rs/zerolog/pkgerrors"
)

// NewLogger returns logger based on server config
func NewLogger(kw *KafkaLogWriter, cw, fw io.Writer) *zerolog.Logger {
	var writers []io.Writer

	// Setting up kafka writer if True.
	if kw != nil {
		wr := diode.NewWriter(kw, 1000, 10*time.Millisecond, func(missed int) {
			fmt.Printf("Logger Dropped %d messages", missed)
		})
		writers = append(writers, wr)
	}

	// Setting up console writer if True.
	if cw != nil {
		writers = append(writers, cw)
	}

	// Setting up file writer is True.
	if fw != nil {
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
