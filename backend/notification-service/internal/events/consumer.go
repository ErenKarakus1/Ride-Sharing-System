package events

import "context"

type RideEventHandler interface {
	HandleRideEvent(ctx context.Context, event RideEvent) error
}

type PaymentEventHandler interface {
	HandlePaymentEvent(ctx context.Context, event PaymentEvent) error
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

type PaymentEvent struct {
	Type       string           `json:"type"`
	OccurredAt string           `json:"occurred_at"`
	Data       PaymentEventData `json:"data"`
}

type PaymentEventData struct {
	ID       string  `json:"id"`
	RideID   string  `json:"ride_id"`
	RiderID  string  `json:"rider_id"`
	DriverID *string `json:"driver_id"`
	Status   string  `json:"status"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}
