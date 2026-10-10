package events

type PaymentEvent struct {
	Type string           `json:"type"`
	Data PaymentEventData `json:"data"`
}

type PaymentEventData struct {
	ID     string `json:"id"`
	RideID string `json:"ride_id"`
	Status string `json:"status"`
}
