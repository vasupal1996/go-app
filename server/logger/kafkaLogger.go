package logger

import (
	"context"
	"go-app/server/config"
	"io"
	"log"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"github.com/segmentio/kafka-go"
)

// KafkaLogWriter extends the existing kafka.Writer functionality by implementing io.Writer`s Write method.
type KafkaLogWriter struct {
	*kafka.Writer
}

// NewKafkaLogWriter returns new instance of KafkaLogWriter
func NewKafkaLogWriter(topic string, dialer *kafka.Dialer, c *config.KafkaConfig) *KafkaLogWriter {
	kw := kafka.NewWriter(kafka.WriterConfig{
		Brokers:  c.Brokers,
		Topic:    topic,
		Async:    true,
		Balancer: &kafka.LeastBytes{},
		Dialer:   dialer,
	})
	return &KafkaLogWriter{
		kw,
	}
}

// Write method implements io.Writer`s Write method.
func (kw *KafkaLogWriter) Write(p []byte) (n int, err error) {
	m := kafka.Message{
		Value: p,
	}
	err = kw.WriteMessages(context.Background(), m)
	if err != nil {
		log.Printf("failed to write log message `%s` to kafka topic: %s", p, err)
		return 0, err
	}
	return len(p), nil
}

// NewKafkaLogger returns new instance of Kafka Logger
func NewKafkaLogger(kw *KafkaLogWriter) *zerolog.Logger {
	mw := io.Writer(kw)
	zlog := zerolog.New(mw).With().Timestamp().Stack().Caller().Logger()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	return &zlog
}
