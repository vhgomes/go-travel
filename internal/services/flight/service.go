package flight

import "context"

type Service struct {
	repo FlightRepository
}

func NewService(repo FlightRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(ctx context.Context, id string) (*Flight, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListAvailable(ctx context.Context, origin string, destination string, date string) ([]*Flight, error) {
	return s.repo.ListAvailable(ctx, origin, destination, date)
}

func (s *Service) ReserveSeats(ctx context.Context, flightID string, seats int) error {
	return s.repo.ReserveSeats(ctx, flightID, seats)
}
