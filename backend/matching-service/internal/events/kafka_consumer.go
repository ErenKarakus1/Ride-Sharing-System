package events

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

const PaymentEventsTopic = "payment.events"

type PaymentEventHandler interface {
	HandlePaymentEvent(ctx context.Context, event PaymentEvent) error
}

type KafkaConsumer struct {
	reader  *kafka.Reader
	handler PaymentEventHandler
}

func NewKafkaConsumer(brokers string, groupID string, handler PaymentEventHandler) *KafkaConsumer {
	brokerList := splitBrokers(brokers)
	if err := ensureTopic(context.Background(), brokerList, PaymentEventsTopic); err != nil {
		log.Printf("failed to ensure kafka topic %s: %v", PaymentEventsTopic, err)
	}

	return &KafkaConsumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers:                brokerList,
			Topic:                  PaymentEventsTopic,
			GroupID:                groupID,
			WatchPartitionChanges:  true,
			PartitionWatchInterval: 5 * time.Second,
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

		var event PaymentEvent
		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("failed to decode payment event: %v", err)
			_ = c.reader.CommitMessages(ctx, message)
			continue
		}

		log.Printf("received payment event type=%s ride_id=%s status=%s", event.Type, event.Data.RideID, event.Data.Status)
		if err := c.handleWithRetry(ctx, event); err != nil {
			log.Printf("failed to handle payment event after retries type=%s ride_id=%s: %v", event.Type, event.Data.RideID, err)
			_ = c.reader.CommitMessages(ctx, message)
			continue
		}

		if err := c.reader.CommitMessages(ctx, message); err != nil {
			return err
		}
	}
}

func (c *KafkaConsumer) handleWithRetry(ctx context.Context, event PaymentEvent) error {
	var err error
	for attempt := 1; attempt <= 3; attempt++ {
		err = c.handler.HandlePaymentEvent(ctx, event)
		if err == nil {
			return nil
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

func ensureTopic(ctx context.Context, brokers []string, topic string) error {
	if len(brokers) == 0 {
		return nil
	}

	conn, err := kafka.DialContext(ctx, "tcp", brokers[0])
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}

	controllerConn, err := kafka.DialContext(ctx, "tcp", net.JoinHostPort(controller.Host, fmt.Sprint(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	return controllerConn.CreateTopics(kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	})
}
