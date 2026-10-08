package events

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
)

const rideEventsTopic = "ride.events"

type KafkaConsumer struct {
	reader  *kafka.Reader
	handler RideEventHandler
}

func NewKafkaConsumer(brokers string, groupID string, handler RideEventHandler) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: splitBrokers(brokers),
			Topic:   rideEventsTopic,
			GroupID: groupID,
		}),
		handler: handler,
	}
}

func (c *KafkaConsumer) Run(ctx context.Context) error {
	for {
		message, err := c.reader.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, context.Canceled) {
				return nil
			}

			return err
		}

		var event RideEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("failed to decode ride event: %v", err)
			_ = c.reader.CommitMessages(ctx, message)
			continue
		}

		if err := c.handler.HandleRideEvent(ctx, event); err != nil {
			return err
		}

		if err := c.reader.CommitMessages(ctx, message); err != nil {
			return err
		}
	}
}

func (c *KafkaConsumer) Close() error {
	return c.reader.Close()
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
