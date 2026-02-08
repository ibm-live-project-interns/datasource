package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

// NewKafkaProducer initializes the connection
func NewKafkaProducer() (*KafkaProducer, error) {
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "kafka:9092" // Default for Docker
	}

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": broker,
		"client.id":         "datasource-service",
		"acks":              "all",
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create kafka producer: %w", err)
	}

	return &KafkaProducer{
		producer: p,
		topic:    "ingestion-events", // Matches Orchestrator config
	}, nil
}

// SendEventAsync pushes the mapped event to Kafka
func (kp *KafkaProducer) SendEventAsync(event interface{}) error {
	// 1. Serialize
	val, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("json marshal failed: %w", err)
	}

	// 2. Produce
	deliveryChan := make(chan kafka.Event)
	err = kp.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kp.topic, Partition: kafka.PartitionAny},
		Value:          val,
	}, deliveryChan)

	if err != nil {
		return err
	}

	// 3. Handle Delivery Report (Async)
	go func() {
		e := <-deliveryChan
		m := e.(*kafka.Message)
		if m.TopicPartition.Error != nil {
			log.Printf("❌ Kafka delivery failed: %v", m.TopicPartition.Error)
		} else {
			// Optional: Log success (can be noisy)
			// log.Printf("✅ Delivered to %v", m.TopicPartition)
		}
		close(deliveryChan)
	}()

	return nil
}

func (kp *KafkaProducer) Close() {
	kp.producer.Flush(15 * 1000)
	kp.producer.Close()
}