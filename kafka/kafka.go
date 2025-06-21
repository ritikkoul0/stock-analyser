package kafka

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
)

var KafkaWriter *kafka.Writer
var KafkaReader *kafka.Reader

func Initialise() {
	// Create SCRAM mechanism with SHA-512 hash
	mechanism, err := scram.Mechanism(scram.SHA512, "stockanalyser1", "Controlstock@123")
	if err != nil {
		panic(err)
	}

	dialer := &kafka.Dialer{
		Timeout:       10 * time.Second,
		DualStack:     true,
		TLS:           &tls.Config{},
		SASLMechanism: mechanism,
	}

	KafkaWriter = kafka.NewWriter(kafka.WriterConfig{
		Brokers:      []string{"d125td00ht3usuv2qhfg.any.ap-south-1.mpx.prd.cloud.redpanda.com:9092"},
		Topic:        "stock-analyser",
		Balancer:     &kafka.LeastBytes{},
		Dialer:       dialer,
		RequiredAcks: int(kafka.RequireAll),
	})

	KafkaReader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"d125td00ht3usuv2qhfg.any.ap-south-1.mpx.prd.cloud.redpanda.com:9092"},
		Topic:    "stock-analyser",
		GroupID:  "stock-analyser-group",
		Dialer:   dialer,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
}

func SendMessage(ctx context.Context, key string, value interface{}) error {
	// Serialize the value to JSON
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// Create the Kafka message
	msg := kafka.Message{
		Key:   []byte(key),
		Value: jsonValue,
	}

	// Write the message using the global KafkaWriter
	err = KafkaWriter.WriteMessages(ctx, msg)
	if err != nil {
		log.Printf("failed to send message to Kafka: %v", err)
	}
	return err
}

func ReadMessage[T any](ctx context.Context) (T, error) {
	var result T

	// Read a message
	msg, err := KafkaReader.ReadMessage(ctx)
	if err != nil {
		log.Printf("Error reading message from Kafka: %v", err)
		return result, err
	}

	// Unmarshal the message into the specified struct
	if err := json.Unmarshal(msg.Value, &result); err != nil {
		log.Printf("Error unmarshalling Kafka message: %v", err)
		return result, err
	}

	log.Printf("Received message from Kafka (offset %d)", msg.Offset)

	return result, nil
}
