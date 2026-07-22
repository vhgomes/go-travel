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
	"github.com/vhgomes/go-travel/pkg/logger"
	"go.uber.org/zap"
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

	logger.Info("create_order_start", zap.String("order_id", order.ID.String()), zap.String("user_id", order.UserID.String()), zap.String("flight_id", order.FlightID), zap.String("hotel_id", order.HotelID))

	exists, err := s.orderRepo.CheckPendingOrderExists(ctx, pgUserID, order.FlightID, order.HotelID)
	if err != nil {
		logger.Error("checking pending order", fmt.Errorf("repo error: %w", err), zap.String("user_id", order.UserID.String()))
		return fmt.Errorf("checking pending order: %w", err)
	}

	if exists {
		logger.Info("order duplicated", zap.String("user_id", order.UserID.String()), zap.String("flight_id", order.FlightID), zap.String("hotel_id", order.HotelID))
		return fmt.Errorf("order duplicated")
	}

	flight, err := s.flightService.GetByID(ctx, order.FlightID)
	if err != nil {
		logger.Error("failed to get flight", fmt.Errorf("flight error: %w", err), zap.String("flight_id", order.FlightID))
		return err
	}

	hotel, err := s.hotelService.GetByID(ctx, order.HotelID)
	if err != nil {
		logger.Error("failed to get hotel", fmt.Errorf("hotel error: %w", err), zap.String("hotel_id", order.HotelID))
		return err
	}

	if flight.AvailableSeats <= 0 || hotel.RoomsAvailable <= 0 {
		return fmt.Errorf("flight or hotel not available")
	}

	if err := s.flightService.ReserveSeats(ctx, order.FlightID, 1); err != nil {
		logger.Error("failed to reserve seats", fmt.Errorf("reserve flight error: %w", err), zap.String("flight_id", order.FlightID), zap.String("order_id", order.ID.String()))
		return err
	}

	if err := s.hotelService.ReserveRooms(ctx, order.HotelID, 1); err != nil {
		logger.Error("failed to reserve rooms", fmt.Errorf("reserve hotel error: %w", err), zap.String("hotel_id", order.HotelID), zap.String("order_id", order.ID.String()))
		return err
	}

	totalAmount := flight.FlightPrice + hotel.HotelPrice

	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()
	order.Status = Pending
	order.TotalAmount = totalAmount

	if err := s.orderRepo.Create(ctx, order); err != nil {
		logger.Error("failed to create order in repo", fmt.Errorf("repo create error: %w", err), zap.String("order_id", order.ID.String()))
		return err
	}

	orderJSON, err := json.Marshal(order)
	if err != nil {
		return fmt.Errorf("marshaling order to JSON: %w", err)
	}

	if _, err := s.producer.SendMessage(ctx, string(orderJSON), nil); err != nil {
		logger.Error("failed to send order to SQS", fmt.Errorf("sqs send error: %w", err), zap.String("order_id", order.ID.String()))
		return fmt.Errorf("sending order to SQS: %w", err)
	}

	logger.Info("order_processed", zap.String("order_id", order.ID.String()), zap.String("user_id", order.UserID.String()))

	return nil
}
