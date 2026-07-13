package order

import (
	"time"

	"github.com/google/uuid"
)

// TODO: isso aqui parece meio errado, creio que irei mudar esses uuids depois
// estou pensando a melhor forma para isso aqui funcionar
type Order struct {
	ID        int64
	UserID    uuid.UUID
	FlightID  uuid.UUID
	HotelID   int64
	PaymentID uuid.UUID
	CreatedAt time.Time
}
