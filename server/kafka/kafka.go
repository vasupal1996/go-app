package kafka

import (
	"fmt"
	"go-app/server/config"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/segmentio/kafka-go"
)

// KImpl used by server to implement Kafka.
type KImpl interface {
	Close()
}

// Kafka has kafka cluster config and connection instance
type Kafka struct {
	Config *config.KafkaConfig
	Conn   *kafka.Conn
}

// Close closes the connection
func (k *Kafka) Close() {
	k.Conn.Close()
}

// NewKafka returns new kafka instance
func NewKafka(c *config.KafkaConfig) *Kafka {
	conn, err := kafka.Dial(c.BrokerDial, fmt.Sprintf("%s:%s", c.BrokerURL, c.BrokerPort))
	if err != nil {
		log.Fatalf("failed to establish kafka connection: %s", err)
		os.Exit(1)
	}
	defer conn.Close()
	controller, err := conn.Controller()
	if err != nil {
		log.Fatalf("failed while establishing connection to controller kafka: %s", err)
		os.Exit(1)
	}
	controllerConn, err := kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		log.Fatalf("failed to establish connection to controller kafka: %s", err)
		os.Exit(1)
	}
	return &Kafka{Config: c, Conn: controllerConn}
}
