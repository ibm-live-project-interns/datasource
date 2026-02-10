//go:build kafka
// +build kafka

// Package client provides clients for communicating with external services.
//
// This file implements a Kafka producer for publishing normalized events to
// the ingestion pipeline. It uses the confluent-kafka-go library which
// requires CGO and librdkafka.
//
// NOTE: The current main.go uses the HTTP IngestorClient (client/ingestor.go)
// for event delivery. This Kafka producer is available as an alternative
// transport for environments where Kafka is deployed. The two approaches
// can coexist — HTTP for direct ingestor communication, Kafka for
// event-driven architectures.
package client

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

// KafkaProducer wraps a confluent-kafka-go producer with topic configuration.
type KafkaProducer struct {
	producer *kafka.Producer
	topic    string
}

// NewKafkaProducer creates a Kafka producer configured via the KAFKA_BROKER
// environment variable. Defaults to "kafka:9092" for Docker environments.
//
// TODO: Make the topic name configurable via environment variable or
// constructor parameter instead of hardcoding "ingestion-events".
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

// SendEventAsync serializes the event as JSON and produces it to the configured
// Kafka topic asynchronously. Delivery success/failure is handled in a background
// goroutine via the delivery channel.
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
			log.Printf("kafka delivery failed: %v", m.TopicPartition.Error)
		} else {
			// Optional: Log success (can be noisy)
			// log.Printf("✅ Delivered to %v", m.TopicPartition)
		}
		close(deliveryChan)
	}()

	return nil
}

// Close flushes any pending messages (up to 15 seconds) and closes the
// underlying Kafka producer. Always call Close before program exit to
// ensure all buffered messages are delivered.
func (kp *KafkaProducer) Close() {
	kp.producer.Flush(15 * 1000)
	kp.producer.Close()
}