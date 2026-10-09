package notification

import (
	"context"
	"log"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/notification-service/internal/events"
)

type Service struct {
	hub *Hub
}

func NewService(hub *Hub) *Service {
	return &Service{
		hub: hub,
	}
}

func (s *Service) HandleRideEvent(ctx context.Context, event events.RideEvent) error {
	log.Printf(
		"notification queued for ride event type=%s ride_id=%s rider_id=%s status=%s",
		event.Type,
		event.Data.ID,
		event.Data.RiderID,
		event.Data.Status,
	)

	if err := s.hub.Send(event.Data.RiderID, Notification{
		Type: event.Type,
		Payload: map[string]any{
			"ride_id": event.Data.ID,
			"status":  event.Data.Status,
		},
	}); err != nil {
		return err
	}

	if event.Data.DriverID != nil {
		return s.hub.Send(*event.Data.DriverID, Notification{
			Type: event.Type,
			Payload: map[string]any{
				"ride_id": event.Data.ID,
				"status":  event.Data.Status,
			},
		})
	}

	return nil
}

func (s *Service) HandlePaymentEvent(ctx context.Context, event events.PaymentEvent) error {
	log.Printf(
		"notification queued for payment event type=%s payment_id=%s ride_id=%s rider_id=%s status=%s",
		event.Type,
		event.Data.ID,
		event.Data.RideID,
		event.Data.RiderID,
		event.Data.Status,
	)

	payload := map[string]any{
		"payment_id": event.Data.ID,
		"ride_id":    event.Data.RideID,
		"status":     event.Data.Status,
		"amount":     event.Data.Amount,
		"currency":   event.Data.Currency,
	}
	if err := s.hub.Send(event.Data.RiderID, Notification{
		Type:    event.Type,
		Payload: payload,
	}); err != nil {
		return err
	}

	if event.Data.DriverID != nil {
		return s.hub.Send(*event.Data.DriverID, Notification{
			Type:    event.Type,
			Payload: payload,
		})
	}

	return nil
}
