package notification

import (
	"context"
	"log"

	"github.com/ErenKarakus1/Ride-Sharing-System/backend/notification-service/internal/events"
)

type Service struct{}

func NewService() *Service {
	return &Service{}
}

func (s *Service) HandleRideEvent(ctx context.Context, event events.RideEvent) error {
	log.Printf(
		"notification queued for ride event type=%s ride_id=%s rider_id=%s status=%s",
		event.Type,
		event.Data.ID,
		event.Data.RiderID,
		event.Data.Status,
	)

	return nil
}
