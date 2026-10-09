package events

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

const PaymentEventsTopic = "payment.events"

type Publisher interface {
	Publish(ctx context.Context, event Event) error
}

type KafkaPublisher struct {
	writer *kafka.Writer
}

func NewKafkaPublisher(brokers string) *KafkaPublisher {
	return &KafkaPublisher{
		writer: &kafka.Writer{
			Addr:         kafka.TCP(splitBrokers(brokers)...),
			Topic:        PaymentEventsTopic,
			Balancer:     &kafka.LeastBytes{},
			RequiredAcks: kafka.RequireOne,
			Async:        false,
		},
	}
}

func (p *KafkaPublisher) Publish(ctx context.Context, event Event) error {
	payload, err := json.Marshal(struct {
		Type       string `json:"type"`
		OccurredAt string `json:"occurred_at"`
		Data       any    `json:"data"`
	}{
		Type:       event.Type,
		OccurredAt: time.Now().UTC().Format(time.RFC3339Nano),
		Data:       event.Data,
	})
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(event.Key),
		Value: payload,
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
