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

// SegmentioKafkaImpl has kafka cluster config and connection instance
type SegmentioKafkaImpl struct {
	Config *config.KafkaConfig
	Conn   *kafka.Conn
}

// Close closes the connection
func (k *SegmentioKafkaImpl) Close() {
	k.Conn.Close()
}

// NewSegmentioKafka returns new segmentio kafka client instance
func NewSegmentioKafka(c *config.KafkaConfig) *SegmentioKafkaImpl {
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
	return &SegmentioKafkaImpl{Config: c, Conn: controllerConn}
}
