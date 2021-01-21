package kafka

import (
	"go-app/server/config"
	"log"
	"os"

	"github.com/Shopify/sarama"
)

// SaramaKafkaImpl implements kafka interface
type SaramaKafkaImpl struct {
	Config *config.KafkaConfig
	Conn   sarama.Client
}

// NewSaramaKafka returns new sarama kafka client instance
func NewSaramaKafka(c *config.KafkaConfig) *SaramaKafkaImpl {
	saramaConfig := sarama.NewConfig()
	saramaConfig.Consumer.Offsets.AutoCommit.Enable = false
	saramaConfig.Consumer.Return.Errors = true
	sarama, err := sarama.NewClient(c.Brokers, saramaConfig)
	if err != nil {
		log.Printf("failed to create sarama client instance: %s", err)
		os.Exit(1)
	}
	return &SaramaKafkaImpl{Config: c, Conn: sarama}
}

// Close terminates client connection with kafka brokers
func (k *SaramaKafkaImpl) Close() {
	k.Conn.Close()
}
