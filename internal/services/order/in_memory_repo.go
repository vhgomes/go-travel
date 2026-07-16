package order

import (
	"context"
	"sync"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type InMemoryOrderRepository struct {
	mu     sync.RWMutex
	orders []Order
}

func NewInMemoryOrderRepository() *InMemoryOrderRepository {
	return &InMemoryOrderRepository{}
}

var _ OrderRepository = (*InMemoryOrderRepository)(nil)

func (r *InMemoryOrderRepository) Create(ctx context.Context, order Order) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.orders = append(r.orders, order)
	return nil
}

func (r *InMemoryOrderRepository) GetOrdersByUserID(ctx context.Context, pgUserID pgtype.UUID) ([]Order, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []Order
	for _, o := range r.orders {
		if o.UserID == uuid.UUID(pgUserID.Bytes) {
			result = append(result, o)
		}
	}
	return result, nil
}

func (r *InMemoryOrderRepository) CheckPendingOrderExists(ctx context.Context, userID pgtype.UUID, flightID, hotelID string) (bool, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, o := range r.orders {
		if o.UserID == uuid.UUID(userID.Bytes) &&
			o.FlightID == flightID &&
			o.HotelID == hotelID &&
			(o.Status == Pending || o.Status == FlightOk || o.Status == HotelOk) {
			return true, nil
		}
	}
	return false, nil
}
