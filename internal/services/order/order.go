package order

import (
	"time"

	"github.com/google/uuid"
)

type Status string
type Currency string

const (
	Pending   Status = "pending"
	FlightOk  Status = "flight_ok"
	HotelOk   Status = "hotel_ok"
	Confirmed Status = "confirmed"
	Rollback  Status = "rollback"
	Failed    Status = "failed"
)

const (
	USD Currency = "USD"
	EUR Currency = "EUR"
)

type Order struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"user_id"`
	FlightID     string    `json:"flight_id"`
	HotelID      string    `json:"hotel_id"`
	PaymentToken string    `json:"payment_token"`
	TotalAmount  int64     `json:"total_amount"` // em centavos
	Currency     Currency  `json:"currency"`
	Status       Status    `json:"status"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
