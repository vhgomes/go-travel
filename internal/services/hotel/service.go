package hotel

import "context"

type Service struct {
	repo HotelRepository
}

func NewService(repo HotelRepository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByID(ctx context.Context, id string) (*Hotel, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *Service) ListAvailable(ctx context.Context, location string, date string) ([]*Hotel, error) {
	return s.repo.ListAvailable(ctx, location, date)
}

func (s *Service) ReserveRooms(ctx context.Context, hotelID string, rooms int) error {
	return s.repo.ReserveRooms(ctx, hotelID, rooms)
}
