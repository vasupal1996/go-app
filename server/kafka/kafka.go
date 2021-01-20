package kafka

import (
	"context"
)

// Kafka used by server to implement Kafka.
type Kafka interface {
	Close()
}

// Message defines kafka message implementation
type Message interface{}

// Consumer defines how a kafka consumer should be implemented
type Consumer interface {
	Init(*Kafka)
	Consume(context.Context, func(*Message))
	Commit(context.Context, *Message)
	ConsumeAndCommit(context.Context, func(*Message))
	Close()
}

// Producer defines how kafka producer should be implemented
type Producer interface {
	Init(*Kafka)
	Publish(*Kafka)
	Close()
}
