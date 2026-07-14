package flight

import "context"

type FlightService struct {
	repo FlightRepository
}

func NewFlightService(repo FlightRepository) *FlightService {
	return &FlightService{repo: repo}
}

func (s *FlightService) GetByID(ctx context.Context, id string) (*Flight, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *FlightService) ListAvailable(ctx context.Context, origin string, destination string, date string) ([]*Flight, error) {
	return s.repo.ListAvailable(ctx, origin, destination, date)
}

func (s *FlightService) ReserveSeats(ctx context.Context, flightID string, seats int) error {
	return s.repo.ReserveSeats(ctx, flightID, seats)
}
