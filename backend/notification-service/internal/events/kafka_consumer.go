package events

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	rideEventsTopic    = "ride.events"
	paymentEventsTopic = "payment.events"
)

type KafkaConsumer struct {
	reader  *kafka.Reader
	handler MessageHandler
}

type MessageHandler interface {
	HandleMessage(ctx context.Context, value []byte) error
}

func NewKafkaConsumer(brokers string, groupID string, handler RideEventHandler) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: splitBrokers(brokers),
			Topic:   rideEventsTopic,
			GroupID: groupID,
		}),
		handler: rideMessageHandler{handler: handler},
	}
}

func NewPaymentKafkaConsumer(brokers string, groupID string, handler PaymentEventHandler) *KafkaConsumer {
	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: splitBrokers(brokers),
			Topic:   paymentEventsTopic,
			GroupID: groupID,
		}),
		handler: paymentMessageHandler{handler: handler},
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

		if err := c.handleWithRetry(ctx, message.Value); err != nil {
			if errors.Is(err, errInvalidMessage) {
				log.Printf("failed to decode event: %v", err)
				_ = c.reader.CommitMessages(ctx, message)
				continue
			}

			log.Printf("failed to handle event after retries: %v", err)
			_ = c.reader.CommitMessages(ctx, message)
			continue
		}

		if err := c.reader.CommitMessages(ctx, message); err != nil {
			return err
		}
	}
}

func (c *KafkaConsumer) handleWithRetry(ctx context.Context, value []byte) error {
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		err = c.handler.HandleMessage(ctx, value)
		if err == nil || errors.Is(err, errInvalidMessage) {
			return err
		}

		timer := time.NewTimer(time.Duration(attempt) * 250 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil
		case <-timer.C:
		}
	}

	return err
}

type rideMessageHandler struct {
	handler RideEventHandler
}

func (h rideMessageHandler) HandleMessage(ctx context.Context, value []byte) error {
	var event RideEvent
	if err := json.Unmarshal(value, &event); err != nil {
		return errors.Join(errInvalidMessage, err)
	}

	return h.handler.HandleRideEvent(ctx, event)
}

type paymentMessageHandler struct {
	handler PaymentEventHandler
}

func (h paymentMessageHandler) HandleMessage(ctx context.Context, value []byte) error {
	var event PaymentEvent
	if err := json.Unmarshal(value, &event); err != nil {
		return errors.Join(errInvalidMessage, err)
	}

	return h.handler.HandlePaymentEvent(ctx, event)
}

var errInvalidMessage = errors.New("invalid kafka message")

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
