package user

import "time"

type Role string

const (
	RoleRider  Role = "rider"
	RoleDriver Role = "driver"
)

func (r Role) Valid() bool {
	return r == RoleRider || r == RoleDriver
}

type User struct {
	ID          string    `json:"id"`
	Email       string    `json:"email"`
	DisplayName string    `json:"display_name"`
	PhoneNumber string    `json:"phone_number"`
	Role        Role      `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
