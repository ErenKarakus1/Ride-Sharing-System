package events

import "context"

type RideEventHandler interface {
	HandleRideEvent(ctx context.Context, event RideEvent) error
}

type RideEvent struct {
	Type       string        `json:"type"`
	OccurredAt string        `json:"occurred_at"`
	Data       RideEventData `json:"data"`
}

type RideEventData struct {
	ID       string  `json:"id"`
	RiderID  string  `json:"rider_id"`
	DriverID *string `json:"driver_id"`
	Status   string  `json:"status"`
}
