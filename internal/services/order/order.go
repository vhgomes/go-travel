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
	ID           uuid.UUID
	UserID       uuid.UUID
	FlightID     string
	HotelID      string
	PaymentToken string
	TotalAmount  int64 // em centavos
	Currency     Currency
	status       Status
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
