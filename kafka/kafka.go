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
