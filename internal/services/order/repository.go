package order

import (
	"context"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vhgomes/go-travel/internal/repository"
)

type OrderRepository struct {
	q *repository.Queries
}

func NewOrderRepository(pool *pgxpool.Pool) *OrderRepository {
	return &OrderRepository{q: repository.New(pool)}
}

func (r *OrderRepository) Create(ctx context.Context, order Order) error {
	_, err := r.q.CreateOrder(ctx, repository.CreateOrderParams{
		ID: pgtype.UUID{
			Bytes: order.ID,
			Valid: true,
		},
		UserID: pgtype.UUID{
			Bytes: order.UserID,
			Valid: true,
		},
		FlightID:     order.FlightID,
		HotelID:      order.HotelID,
		PaymentToken: order.PaymentToken,
		TotalAmount: pgtype.Numeric{
			Int:   big.NewInt(order.TotalAmount),
			Valid: true,
		},
		Currency: string(order.Currency),
	})
	return err
}

func (r *OrderRepository) GetOrdersByUserID(ctx context.Context, pgUserID pgtype.UUID) ([]Order, error) {
	orders, err := r.q.GetOrdersByUserID(ctx, repository.GetOrdersByUserIDParams{
		UserID: pgUserID,
	})

	if err != nil {
		return nil, err
	}

	var result []Order
	for _, o := range orders {
		result = append(result, Order{
			ID:           o.ID.Bytes,
			UserID:       o.UserID.Bytes,
			FlightID:     o.FlightID,
			HotelID:      o.HotelID,
			PaymentToken: o.PaymentToken,
			TotalAmount:  o.TotalAmount.Int.Int64(),
			Currency:     Currency(o.Currency),
			Status:       Status(o.Status),
			CreatedAt:    o.CreatedAt.Time,
			UpdatedAt:    o.UpdatedAt.Time,
		})
	}

	return result, nil
}
