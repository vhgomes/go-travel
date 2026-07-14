package flight

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFlightService_GetByID(t *testing.T) {
	repo := NewInMemoryFlightRepository()
	svc := NewFlightService(repo)

	ctx := context.Background()

	// Cria um voo de teste
	flight := &Flight{
		ID:             "FL123",
		Origin:         "GRU",
		Destination:    "JFK",
		Date:           "2025-07-20",
		FlightPrice:    45000,
		AvailableSeats: 10,
	}
	err := repo.Create(ctx, flight)
	require.NoError(t, err)

	t.Run("encontra voo existente", func(t *testing.T) {
		result, err := svc.GetByID(ctx, "FL123")
		assert.NoError(t, err)
		assert.Equal(t, "FL123", result.ID)
		assert.Equal(t, "GRU", result.Origin)
		assert.Equal(t, 10, result.AvailableSeats)
	})

	t.Run("voo não encontrado", func(t *testing.T) {
		_, err := svc.GetByID(ctx, "FL999")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "flight not found")
	})
}

func TestFlightService_ListAvailable(t *testing.T) {
	repo := NewInMemoryFlightRepository()
	svc := NewFlightService(repo)

	ctx := context.Background()

	// Cria voos de teste
	flights := []*Flight{
		{ID: "FL123", Origin: "GRU", Destination: "JFK", Date: "2025-07-20", FlightPrice: 45000, AvailableSeats: 10},
		{ID: "FL124", Origin: "GRU", Destination: "JFK", Date: "2025-07-20", FlightPrice: 35000, AvailableSeats: 0},
		{ID: "FL125", Origin: "GRU", Destination: "LIS", Date: "2025-07-21", FlightPrice: 50000, AvailableSeats: 5},
	}

	for _, f := range flights {
		err := repo.Create(ctx, f)
		require.NoError(t, err)
	}

	t.Run("lista voos disponíveis com filtros", func(t *testing.T) {
		result, err := svc.ListAvailable(ctx, "GRU", "JFK", "2025-07-20")
		assert.NoError(t, err)
		assert.Len(t, result, 1) // apenas FL123 (FL124 tem 0 assentos)
		assert.Equal(t, "FL123", result[0].ID)
	})

	t.Run("nenhum voo disponível", func(t *testing.T) {
		result, err := svc.ListAvailable(ctx, "GRU", "MIA", "2025-07-20")
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestFlightService_ReserveSeats(t *testing.T) {
	repo := NewInMemoryFlightRepository()
	svc := NewFlightService(repo)

	ctx := context.Background()

	// Cria voo com 10 assentos
	flight := &Flight{
		ID:             "FL123",
		Origin:         "GRU",
		Destination:    "JFK",
		Date:           "2025-07-20",
		FlightPrice:    45000,
		AvailableSeats: 10,
	}
	err := repo.Create(ctx, flight)
	require.NoError(t, err)

	t.Run("reserva assentos com sucesso", func(t *testing.T) {
		err := svc.ReserveSeats(ctx, "FL123", 3)
		assert.NoError(t, err)

		// Verifica se os assentos foram reduzidos
		updated, err := repo.GetByID(ctx, "FL123")
		assert.NoError(t, err)
		assert.Equal(t, 7, updated.AvailableSeats)
	})

	t.Run("tenta reservar mais assentos que disponíveis", func(t *testing.T) {
		err := svc.ReserveSeats(ctx, "FL123", 10) // só tem 7
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient seats")
	})

	t.Run("tenta reservar em voo inexistente", func(t *testing.T) {
		err := svc.ReserveSeats(ctx, "FL999", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "flight not found")
	})
}
