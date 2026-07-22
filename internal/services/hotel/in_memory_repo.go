package hotel

import (
	"context"
	"fmt"
	"sync"

	"github.com/vhgomes/go-travel/pkg/logger"
	"go.uber.org/zap"
)

type InMemoryHotelRepository struct {
	mu     sync.RWMutex
	hotels map[string]*Hotel
}

func NewInMemoryHotelRepository() *InMemoryHotelRepository {
	return &InMemoryHotelRepository{
		hotels: make(map[string]*Hotel),
	}
}

func (r *InMemoryHotelRepository) GetByID(ctx context.Context, id string) (*Hotel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	hotel, exists := r.hotels[id]
	if !exists {
		logger.Error("hotel not found", fmt.Errorf("not found: %s", id), zap.String("hotel_id", id))
		return nil, fmt.Errorf("hotel not found: %s", id)
	}

	hotelCopy := *hotel
	return &hotelCopy, nil
}

func (r *InMemoryHotelRepository) ListAvailable(ctx context.Context, location string, date string) ([]*Hotel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var available []*Hotel

	for _, hotel := range r.hotels {
		if hotel.Location == location && hotel.Date == date && hotel.RoomsAvailable > 0 {
			hotelCopy := *hotel
			available = append(available, &hotelCopy)
		}
	}

	return available, nil
}

func (r *InMemoryHotelRepository) ReserveRooms(ctx context.Context, hotelID string, rooms int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	hotel, exists := r.hotels[hotelID]
	if !exists {
		logger.Error("hotel not found for reserve", fmt.Errorf("not found: %s", hotelID), zap.String("hotel_id", hotelID))
		return fmt.Errorf("hotel not found: %s", hotelID)
	}

	if hotel.RoomsAvailable < rooms {
		logger.Warn("insufficient rooms available", zap.Int("requested", rooms), zap.Int("available", hotel.RoomsAvailable), zap.String("hotel_id", hotelID))
		return fmt.Errorf("insufficient rooms available: requested %d, available %d", rooms, hotel.RoomsAvailable)
	}

	hotel.RoomsAvailable -= rooms
	return nil
}

func (r *InMemoryHotelRepository) GetAll(ctx context.Context) ([]*Hotel, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var hotels []*Hotel
	for _, hotel := range r.hotels {
		hotelCopy := *hotel
		hotels = append(hotels, &hotelCopy)
	}

	return hotels, nil
}

func (r *InMemoryHotelRepository) Create(ctx context.Context, hotel *Hotel) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.hotels[hotel.ID]; exists {
		logger.Warn("hotel already exists", zap.String("hotel_id", hotel.ID))
		return fmt.Errorf("hotel already exists: %s", hotel.ID)
	}

	hotelCopy := *hotel
	r.hotels[hotel.ID] = &hotelCopy
	return nil
}

func (r *InMemoryHotelRepository) UpdateAvailableRooms(ctx context.Context, hotelID string, availableRooms int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	hotel, exists := r.hotels[hotelID]
	if !exists {
		logger.Error("hotel not found for update", fmt.Errorf("not found: %s", hotelID), zap.String("hotel_id", hotelID))
		return fmt.Errorf("hotel not found: %s", hotelID)
	}

	hotel.RoomsAvailable = availableRooms
	return nil
}
