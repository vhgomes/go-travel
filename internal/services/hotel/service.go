package hotel

import "context"

type HotelService struct {
	repo HotelRepository
}

func NewHotelService(repo HotelRepository) *HotelService {
	return &HotelService{repo: repo}
}

func (s *HotelService) GetByID(ctx context.Context, id string) (*Hotel, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *HotelService) ListAvailable(ctx context.Context, location string, date string) ([]*Hotel, error) {
	return s.repo.ListAvailable(ctx, location, date)
}

func (s *HotelService) ReserveRooms(ctx context.Context, hotelID string, rooms int) error {
	return s.repo.ReserveRooms(ctx, hotelID, rooms)
}
