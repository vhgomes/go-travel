package hotel

import "context"

type HotelRepository interface {
	GetByID(ctx context.Context, id string) (*Hotel, error)
	ListAvailable(ctx context.Context, location string, date string) ([]*Hotel, error)
	ReserveRooms(ctx context.Context, hotelID string, rooms int) error
	GetAll(ctx context.Context) ([]*Hotel, error)
	Create(ctx context.Context, hotel *Hotel) error
	UpdateAvailableRooms(ctx context.Context, hotelID string, availableRooms int) error
}
