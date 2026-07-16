package order

import (
	"context"
	"math/big"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vhgomes/go-travel/internal/repository"
)

// PostgresOrderRepository é a implementação de OrderRepository baseada em
// Postgres (via sqlc + pgx). Para testes, use InMemoryOrderRepository.
type PostgresOrderRepository struct {
	q *repository.Queries
}

func NewOrderRepository(pool *pgxpool.Pool) *PostgresOrderRepository {
	return &PostgresOrderRepository{q: repository.New(pool)}
}

var _ OrderRepository = (*PostgresOrderRepository)(nil)

func (r *PostgresOrderRepository) Create(ctx context.Context, order Order) error {
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
			Exp:   -2,
		},
		Currency: string(order.Currency),
	})
	return err
}

func (r *PostgresOrderRepository) GetOrdersByUserID(ctx context.Context, pgUserID pgtype.UUID) ([]Order, error) {
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

func (r *PostgresOrderRepository) CheckPendingOrderExists(ctx context.Context, userID pgtype.UUID, flightID, hotelID string) (bool, error) {
	count, err := r.q.CheckPendingOrderExists(ctx, repository.CheckPendingOrderExistsParams{
		UserID:   userID,
		FlightID: flightID,
		HotelID:  hotelID,
	})

	if err != nil {
		return false, err
	}

	return count > 0, nil
}
