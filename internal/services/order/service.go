package order

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vhgomes/go-travel/internal/services/flight"
	"github.com/vhgomes/go-travel/internal/services/hotel"
)

type OrderService struct {
	hotelService  *hotel.HotelService
	flightService *flight.FlightService
	repo          OrderRepository
	sqs           *sqs.SQS
}

func NewOrderService(repo OrderRepository, flightService *flight.FlightService, hotelService *hotel.HotelService, sqs *sqs.SQS) *OrderService {
	return &OrderService{repo: repo, flightService: flightService, hotelService: hotelService, sqs: sqs}
}

func (s *OrderService) Create(ctx context.Context, order Order) error {
	pgUserID := pgtype.UUID{
		Bytes: order.UserID,
		Valid: true,
	}

	orders, err := s.repo.GetOrdersByUserID(ctx, pgUserID)

	if err != nil {
		return fmt.Errorf("get orders by user id failed")
	}

	for _, o := range orders {
		if o.FlightID == order.FlightID && o.HotelID == order.HotelID && o.Status == Pending {
			return fmt.Errorf("order duplicated")
		}
	}

	flight, err := s.flightService.GetByID(ctx, order.FlightID)
	if err != nil {
		return err
	}

	hotel, err := s.hotelService.GetByID(ctx, order.HotelID)
	if err != nil {
		return err
	}

	if flight.AvailableSeats <= 0 || hotel.RoomsAvailable <= 0 {
		return fmt.Errorf("flight or hotel not available")
	}

	if err := s.flightService.ReserveSeats(ctx, order.FlightID, 1); err != nil {
		return err
	}

	if err := s.hotelService.ReserveRooms(ctx, order.HotelID, 1); err != nil {
		return err
	}

	totalAmount := flight.FlightPrice + hotel.HotelPrice

	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.Status = Pending
	order.TotalAmount = totalAmount

	if err := s.repo.Create(ctx, order); err != nil {
		return err
	}

	// TODO: enviar para o SQS
	if err := s.sqs.SendOrder(ctx, order); err != nil {
		return err
	}

	return nil
}
