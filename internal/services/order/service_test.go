package order

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/vhgomes/go-travel/internal/services/flight"
	"github.com/vhgomes/go-travel/internal/services/hotel"
)

type mockProducer struct {
	sentBody string
	called   bool
	err      error
}

func (m *mockProducer) SendMessage(ctx context.Context, body string, attributes map[string]string) (string, error) {
	m.called = true
	m.sentBody = body
	return "mock-id", m.err
}

func TestOrderService_Create_Success(t *testing.T) {
	ctx := context.Background()

	flightRepo := flight.NewInMemoryFlightRepository()
	hotelRepo := hotel.NewInMemoryHotelRepository()

	flightService := flight.NewFlightService(flightRepo)
	hotelService := hotel.NewHotelService(hotelRepo)
	orderRepo := NewInMemoryOrderRepository()
	producer := &mockProducer{}

	flight := &flight.Flight{
		ID:             "FL123",
		Origin:         "GRU",
		Destination:    "JFK",
		Date:           "2025-07-20",
		FlightPrice:    45000,
		AvailableSeats: 10,
	}
	hotel := &hotel.Hotel{
		ID:             "HT123",
		Name:           "Hilton Garden Inn",
		Location:       "New York",
		Date:           "2025-07-20",
		HotelPrice:     25000,
		RoomsAvailable: 10,
	}

	err := flightRepo.Create(ctx, flight)
	require.NoError(t, err)
	err = hotelRepo.Create(ctx, hotel)
	require.NoError(t, err)

	orderService := NewOrderService(orderRepo, flightService, hotelService, producer)
	order := Order{
		ID:           uuid.New(),
		UserID:       uuid.New(),
		FlightID:     "FL123",
		HotelID:      "HT123",
		PaymentToken: "token-123",
		Currency:     USD,
	}

	err = orderService.Create(ctx, order)
	require.NoError(t, err)
	assert.True(t, producer.called)

	var createdOrder Order
	assert.NoError(t, json.Unmarshal([]byte(producer.sentBody), &createdOrder))
	assert.Equal(t, order.ID, createdOrder.ID)
	assert.Equal(t, Pending, createdOrder.Status)
	assert.Equal(t, int64(70000), createdOrder.TotalAmount)

	updatedFlight, err := flightRepo.GetByID(ctx, "FL123")
	require.NoError(t, err)
	assert.Equal(t, 9, updatedFlight.AvailableSeats)

	updatedHotel, err := hotelRepo.GetByID(ctx, "HT123")
	require.NoError(t, err)
	assert.Equal(t, 9, updatedHotel.RoomsAvailable)
}

func TestOrderService_Create_DuplicateOrder(t *testing.T) {
	ctx := context.Background()

	flightRepo := flight.NewInMemoryFlightRepository()
	hotelRepo := hotel.NewInMemoryHotelRepository()

	flightService := flight.NewFlightService(flightRepo)
	hotelService := hotel.NewHotelService(hotelRepo)
	orderRepo := NewInMemoryOrderRepository()
	producer := &mockProducer{}

	flight := &flight.Flight{
		ID:             "FL123",
		Origin:         "GRU",
		Destination:    "JFK",
		Date:           "2025-07-20",
		FlightPrice:    45000,
		AvailableSeats: 10,
	}
	hotel := &hotel.Hotel{
		ID:             "HT123",
		Name:           "Hilton Garden Inn",
		Location:       "New York",
		Date:           "2025-07-20",
		HotelPrice:     25000,
		RoomsAvailable: 10,
	}

	err := flightRepo.Create(ctx, flight)
	require.NoError(t, err)
	err = hotelRepo.Create(ctx, hotel)
	require.NoError(t, err)

	orderService := NewOrderService(orderRepo, flightService, hotelService, producer)
	order := Order{
		ID:           uuid.New(),
		UserID:       uuid.MustParse("00000000-0000-0000-0000-000000000001"),
		FlightID:     "FL123",
		HotelID:      "HT123",
		PaymentToken: "token-123",
		Currency:     USD,
	}

	err = orderService.Create(ctx, order)
	require.NoError(t, err)

	err = orderService.Create(ctx, order)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "order duplicated")
}

func TestOrderService_Create_MessageFailure(t *testing.T) {
	ctx := context.Background()

	flightRepo := flight.NewInMemoryFlightRepository()
	hotelRepo := hotel.NewInMemoryHotelRepository()

	flightService := flight.NewFlightService(flightRepo)
	hotelService := hotel.NewHotelService(hotelRepo)
	orderRepo := NewInMemoryOrderRepository()
	producer := &mockProducer{err: errors.New("send failed")}

	flight := &flight.Flight{
		ID:             "FL123",
		Origin:         "GRU",
		Destination:    "JFK",
		Date:           "2025-07-20",
		FlightPrice:    45000,
		AvailableSeats: 10,
	}
	hotel := &hotel.Hotel{
		ID:             "HT123",
		Name:           "Hilton Garden Inn",
		Location:       "New York",
		Date:           "2025-07-20",
		HotelPrice:     25000,
		RoomsAvailable: 10,
	}

	err := flightRepo.Create(ctx, flight)
	require.NoError(t, err)
	err = hotelRepo.Create(ctx, hotel)
	require.NoError(t, err)

	orderService := NewOrderService(orderRepo, flightService, hotelService, producer)
	order := Order{
		ID:           uuid.New(),
		UserID:       uuid.New(),
		FlightID:     "FL123",
		HotelID:      "HT123",
		PaymentToken: "token-123",
		Currency:     USD,
	}

	err = orderService.Create(ctx, order)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "sending order to SQS")
	assert.True(t, producer.called)
}
