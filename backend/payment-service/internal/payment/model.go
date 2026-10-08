package payment

import "time"

type Status string

const (
	StatusAuthorized Status = "authorized"
	StatusCaptured   Status = "captured"
	StatusRefunded   Status = "refunded"
	StatusFailed     Status = "failed"
)

type Payment struct {
	ID        string    `json:"id"`
	RideID    string    `json:"ride_id"`
	RiderID   string    `json:"rider_id"`
	DriverID  *string   `json:"driver_id,omitempty"`
	Amount    float64   `json:"amount"`
	Currency  string    `json:"currency"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
