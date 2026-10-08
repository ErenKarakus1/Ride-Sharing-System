package ride

import "time"

type Status string

const (
	StatusRequested Status = "requested"
	StatusAccepted  Status = "accepted"
	StatusStarted   Status = "started"
	StatusCompleted Status = "completed"
	StatusCancelled Status = "cancelled"
)

type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Address   string  `json:"address"`
}

type Ride struct {
	ID        string    `json:"id"`
	RiderID   string    `json:"rider_id"`
	DriverID  *string   `json:"driver_id,omitempty"`
	Pickup    Location  `json:"pickup"`
	Dropoff   Location  `json:"dropoff"`
	Status    Status    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
