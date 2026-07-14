package flight

import "context"

type FlightRepository interface {
	GetByID(ctx context.Context, id string) (*Flight, error)

	ListAvailable(ctx context.Context, origin string, destination string, date string) ([]*Flight, error)

	ReserveSeats(ctx context.Context, flightID string, seats int) error

	GetAll(ctx context.Context) ([]*Flight, error)

	Create(ctx context.Context, flight *Flight) error

	UpdateAvailableSeats(ctx context.Context, flightID string, availableSeats int) error
}
