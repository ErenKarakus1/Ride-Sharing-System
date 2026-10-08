package events

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

const rideEventsTopic = "ride.events"

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(splitBrokers(brokers)...),
			Topic:        rideEventsTopic,
			Balancer:     &kafka.Hash{},
			RequiredAcks: kafka.RequireOne,
			Async:        false,
		},
	}
}

func (p *KafkaPublisher) Publish(ctx context.Context, event Event) error {
	payload := map[string]any{
		"type":        event.Type,
		"occurred_at": time.Now().UTC(),
		"data":        event.Data,
	}

	value, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.Key),
		Value: value,
	})
}

func (p *KafkaPublisher) Close() error {
	return p.writer.Close()
}

func splitBrokers(brokers string) []string {
	parts := strings.Split(brokers, ",")
	cleaned := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			cleaned = append(cleaned, part)
		}
	}
	if len(cleaned) == 0 {
		return []string{"localhost:9092"}
	}

	return cleaned
}
