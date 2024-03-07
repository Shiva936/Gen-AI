package pubsub

import (
	"log"

	"github.com/IBM/sarama"
)

// Message represents a Kafka message.
type Message struct {
	Key, Value string
}

// KafkaProducer represents a Kafka producer.
type KafkaProducer struct {
	producer sarama.SyncProducer
}

// NewKafkaProducer creates a new KafkaProducer instance.
func NewKafkaProducer(brokers []string) (*KafkaProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true

	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaProducer{producer: producer}, nil
}

// Produce sends a message to the Kafka topic.
func (p *KafkaProducer) Produce(topic string, message *Message) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Key:   sarama.StringEncoder(message.Key),
		Value: sarama.StringEncoder(message.Value),
	}

	_, _, err := p.producer.SendMessage(msg)
	return err
}

// Close closes the Kafka producer.
func (p *KafkaProducer) Close() {
	if err := p.producer.Close(); err != nil {
		log.Println("Error closing Kafka producer:", err)
	}
}
