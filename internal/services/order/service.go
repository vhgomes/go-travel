package order

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/vhgomes/go-travel/internal/services/flight"
	"github.com/vhgomes/go-travel/internal/services/hotel"
	"github.com/vhgomes/go-travel/internal/services/sqs"
)

type OrderService struct {
	hotelService  *hotel.HotelService
	flightService *flight.FlightService
	orderRepo     OrderRepository
	producer      sqs.MessageProducer
}

func NewOrderService(orderRepo OrderRepository, flightService *flight.FlightService, hotelService *hotel.HotelService, producer sqs.MessageProducer) *OrderService {
	return &OrderService{orderRepo: orderRepo, flightService: flightService, hotelService: hotelService, producer: producer}
}

func (s *OrderService) Create(ctx context.Context, order Order) error {
	pgUserID := pgtype.UUID{
		Bytes: order.UserID,
		Valid: true,
	}

	exists, err := s.orderRepo.CheckPendingOrderExists(ctx, pgUserID, order.FlightID, order.HotelID)
	if err != nil {
		return fmt.Errorf("checking pending order: %w", err)
	}

	if exists {
		return fmt.Errorf("order duplicated")
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

	if err := s.orderRepo.Create(ctx, order); err != nil {
		return err
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("marshaling order to JSON: %w", err)
	}

	if _, err := s.producer.SendMessage(ctx, string(orderJSON), nil); err != nil {
		return fmt.Errorf("sending order to SQS: %w", err)
	}

	return nil
}
