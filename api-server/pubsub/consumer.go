package pubsub

import (
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/IBM/sarama"
)

// KafkaConsumer represents a Kafka consumer.
type KafkaConsumer struct {
	consumer  sarama.Consumer
	messages  chan *sarama.ConsumerMessage
	closeChan chan struct{}
	wg        sync.WaitGroup
}

// NewKafkaConsumer creates a new KafkaConsumer instance.
func NewKafkaConsumer(brokers []string, topics []string) (*KafkaConsumer, error) {
	config := sarama.NewConfig()
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		return nil, err
	}

	return &KafkaConsumer{
		consumer:  consumer,
		messages:  make(chan *sarama.ConsumerMessage),
		closeChan: make(chan struct{}),
	}, nil
}

// Consume starts consuming messages from Kafka topics.
func (c *KafkaConsumer) Consume(topics []string) {
	partitions := make(map[int32]struct{})
	for _, topic := range topics {
		partitionIDs, err := c.consumer.Partitions(topic)
		if err != nil {
			log.Println("Error getting partitions for topic", topic, ":", err)
			continue
		}

		for _, partitionID := range partitionIDs {
			partitions[partitionID] = struct{}{}
		}
	}

	for partitionID := range partitions {
		partitionConsumer, err := c.consumer.ConsumePartition(topics[0], partitionID, sarama.OffsetNewest)
		if err != nil {
			log.Println("Error creating partition consumer:", err)
			continue
		}

		c.wg.Add(1)
		go func(pc sarama.PartitionConsumer) {
			defer c.wg.Done()
			for {
				select {
				case msg := <-pc.Messages():
					c.messages <- msg
				case <-c.closeChan:
					return
				}
			}
		}(partitionConsumer)
	}

	go c.handleSignals()
}

// handleSignals waits for a termination signal and closes the consumer gracefully.
func (c *KafkaConsumer) handleSignals() {
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	select {
	case <-signalCh:
		close(c.closeChan)
		c.wg.Wait()
		if err := c.consumer.Close(); err != nil {
			log.Println("Error closing Kafka consumer:", err)
		}
		close(c.messages)
	}
}

// Messages returns the channel for consuming messages.
func (c *KafkaConsumer) Messages() <-chan *sarama.ConsumerMessage {
	return c.messages
}

// Close closes the Kafka consumer.
func (c *KafkaConsumer) Close() {
	close(c.closeChan)
	c.wg.Wait()
	if err := c.consumer.Close(); err != nil {
		log.Println("Error closing Kafka consumer:", err)
	}
	close(c.messages)
}
