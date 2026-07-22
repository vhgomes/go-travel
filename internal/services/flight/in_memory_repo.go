package flight

import (
	"context"
	"fmt"
	"sync"

	"github.com/vhgomes/go-travel/pkg/logger"
	"go.uber.org/zap"
)

type InMemoryFlightRepository struct {
	mu      sync.RWMutex
	flights map[string]*Flight
}

func NewInMemoryFlightRepository() *InMemoryFlightRepository {
	return &InMemoryFlightRepository{
		flights: make(map[string]*Flight),
	}
}

func (r *InMemoryFlightRepository) GetByID(ctx context.Context, id string) (*Flight, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	flight, exists := r.flights[id]
	if !exists {
		logger.Error("flight not found", fmt.Errorf("not found: %s", id), zap.String("flight_id", id))
		return nil, fmt.Errorf("flight not found: %s", id)
	}

	flightCopy := *flight
	return &flightCopy, nil
}

func (r *InMemoryFlightRepository) ListAvailable(ctx context.Context, origin string, destination string, date string) ([]*Flight, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var available []*Flight

	for _, flight := range r.flights {
		if flight.Origin == origin && flight.Destination == destination && flight.Date == date && flight.AvailableSeats > 0 {
			flightCopy := *flight
			available = append(available, &flightCopy)
		}
	}

	return available, nil
}

func (r *InMemoryFlightRepository) ReserveSeats(ctx context.Context, flightID string, seats int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	flight, exists := r.flights[flightID]
	if !exists {
		logger.Error("flight not found for reserve", fmt.Errorf("not found: %s", flightID), zap.String("flight_id", flightID))
		return fmt.Errorf("flight not found: %s", flightID)
	}

	if flight.AvailableSeats < seats {
		logger.Warn("insufficient seats available", zap.Int("requested", seats), zap.Int("available", flight.AvailableSeats), zap.String("flight_id", flightID))
		return fmt.Errorf("insufficient seats available: requested %d, available %d", seats, flight.AvailableSeats)
	}

	flight.AvailableSeats -= seats
	return nil
}

func (r *InMemoryFlightRepository) GetAll(ctx context.Context) ([]*Flight, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var flights []*Flight
	for _, flight := range r.flights {
		flightCopy := *flight
		flights = append(flights, &flightCopy)
	}

	return flights, nil
}

func (r *InMemoryFlightRepository) Create(ctx context.Context, flight *Flight) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.flights[flight.ID]; exists {
		return fmt.Errorf("flight already exists: %s", flight.ID)
	}

	flightCopy := *flight
	r.flights[flight.ID] = &flightCopy
	return nil
}

func (r *InMemoryFlightRepository) UpdateAvailableSeats(ctx context.Context, flightID string, availableSeats int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	flight, exists := r.flights[flightID]
	if !exists {
		logger.Error("flight not found for update", fmt.Errorf("not found: %s", flightID), zap.String("flight_id", flightID))
		return fmt.Errorf("flight not found: %s", flightID)
	}

	flight.AvailableSeats = availableSeats
	return nil
}
