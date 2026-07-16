package order

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

type OrderRepository interface {
	Create(ctx context.Context, order Order) error
	GetOrdersByUserID(ctx context.Context, pgUserID pgtype.UUID) ([]Order, error)
	CheckPendingOrderExists(ctx context.Context, userID pgtype.UUID, flightID, hotelID string) (bool, error)
}
