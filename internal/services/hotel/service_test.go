package hotel

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHotelService_GetByID(t *testing.T) {
	repo := NewInMemoryHotelRepository()
	svc := NewHotelService(repo)

	ctx := context.Background()

	hotel := &Hotel{
		ID:             "HT123",
		Name:           "Hilton Garden Inn",
		Location:       "New York",
		Date:           "2025-07-20",
		HotelPrice:     25000,
		RoomsAvailable: 15,
	}
	err := repo.Create(ctx, hotel)
	require.NoError(t, err)

	t.Run("encontra hotel existente", func(t *testing.T) {
		result, err := svc.GetByID(ctx, "HT123")
		assert.NoError(t, err)
		assert.Equal(t, "HT123", result.ID)
		assert.Equal(t, "Hilton Garden Inn", result.Name)
		assert.Equal(t, "New York", result.Location)
		assert.Equal(t, 15, result.RoomsAvailable)
	})

	t.Run("hotel não encontrado", func(t *testing.T) {
		_, err := svc.GetByID(ctx, "HT999")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "hotel not found")
	})
}

func TestHotelService_ListAvailable(t *testing.T) {
	repo := NewInMemoryHotelRepository()
	svc := NewHotelService(repo)

	ctx := context.Background()

	hotels := []*Hotel{
		{
			ID:             "HT123",
			Name:           "Hilton Garden Inn",
			Location:       "New York",
			Date:           "2025-07-20",
			HotelPrice:     25000,
			RoomsAvailable: 15,
		},
		{
			ID:             "HT124",
			Name:           "Marriott Downtown",
			Location:       "New York",
			Date:           "2025-07-20",
			HotelPrice:     32000,
			RoomsAvailable: 0,
		},
		{
			ID:             "HT125",
			Name:           "Hilton London",
			Location:       "London",
			Date:           "2025-07-21",
			HotelPrice:     28000,
			RoomsAvailable: 10,
		},
		{
			ID:             "HT126",
			Name:           "Holiday Inn NYC",
			Location:       "New York",
			Date:           "2025-07-22",
			HotelPrice:     18000,
			RoomsAvailable: 8,
		},
	}

	for _, h := range hotels {
		err := repo.Create(ctx, h)
		require.NoError(t, err)
	}

	t.Run("lista hotéis disponíveis com filtros", func(t *testing.T) {
		result, err := svc.ListAvailable(ctx, "New York", "2025-07-20")
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "HT123", result[0].ID)
		assert.Equal(t, "Hilton Garden Inn", result[0].Name)
	})

	t.Run("lista hotéis em localização diferente", func(t *testing.T) {
		result, err := svc.ListAvailable(ctx, "London", "2025-07-21")
		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "HT125", result[0].ID)
		assert.Equal(t, "Hilton London", result[0].Name)
	})

	t.Run("nenhum hotel disponível", func(t *testing.T) {
		result, err := svc.ListAvailable(ctx, "Miami", "2025-07-20")
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("nenhum hotel disponível na data", func(t *testing.T) {
		result, err := svc.ListAvailable(ctx, "New York", "2025-07-21")
		assert.NoError(t, err)
		assert.Empty(t, result)
	})
}

func TestHotelService_ReserveRooms(t *testing.T) {
	repo := NewInMemoryHotelRepository()
	svc := NewHotelService(repo)

	ctx := context.Background()

	hotel := &Hotel{
		ID:             "HT123",
		Name:           "Hilton Garden Inn",
		Location:       "New York",
		Date:           "2025-07-20",
		HotelPrice:     25000,
		RoomsAvailable: 10,
	}
	err := repo.Create(ctx, hotel)
	require.NoError(t, err)

	t.Run("reserva quartos com sucesso", func(t *testing.T) {
		err := svc.ReserveRooms(ctx, "HT123", 3)
		assert.NoError(t, err)

		updated, err := repo.GetByID(ctx, "HT123")
		assert.NoError(t, err)
		assert.Equal(t, 7, updated.RoomsAvailable)
	})

	t.Run("tenta reservar mais quartos que disponíveis", func(t *testing.T) {
		err := svc.ReserveRooms(ctx, "HT123", 10) // só tem 7
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "insufficient rooms")
	})

	t.Run("tenta reservar em hotel inexistente", func(t *testing.T) {
		err := svc.ReserveRooms(ctx, "HT999", 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "hotel not found")
	})

	t.Run("reserva quartos após outras reservas", func(t *testing.T) {
		err := svc.ReserveRooms(ctx, "HT123", 2)
		assert.NoError(t, err)

		updated, err := repo.GetByID(ctx, "HT123")
		assert.NoError(t, err)
		assert.Equal(t, 5, updated.RoomsAvailable) // 7 - 2 = 5

		err = svc.ReserveRooms(ctx, "HT123", 3)
		assert.NoError(t, err)

		final, err := repo.GetByID(ctx, "HT123")
		assert.NoError(t, err)
		assert.Equal(t, 2, final.RoomsAvailable) // 5 - 3 = 2
	})
}

func TestHotelService_CreateHotel(t *testing.T) {
	repo := NewInMemoryHotelRepository()
	svc := NewHotelService(repo)

	ctx := context.Background()

	t.Run("cria hotel com sucesso", func(t *testing.T) {
		hotel := &Hotel{
			ID:             "HT999",
			Name:           "Novotel",
			Location:       "Paris",
			Date:           "2025-08-01",
			HotelPrice:     30000,
			RoomsAvailable: 20,
		}

		err := repo.Create(ctx, hotel)
		assert.NoError(t, err)

		result, err := svc.GetByID(ctx, "HT999")
		assert.NoError(t, err)
		assert.Equal(t, "Novotel", result.Name)
		assert.Equal(t, "Paris", result.Location)
		assert.Equal(t, 20, result.RoomsAvailable)
	})

	t.Run("cria hotel com ID duplicado", func(t *testing.T) {
		hotel := &Hotel{
			ID:             "HT123", // ID já existe
			Name:           "Duplicate Hotel",
			Location:       "Boston",
			Date:           "2025-08-01",
			HotelPrice:     20000,
			RoomsAvailable: 10,
		}

		err := repo.Create(ctx, hotel)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "hotel already exists")
	})
}

func TestHotelService_UpdateAvailableRooms(t *testing.T) {
	repo := NewInMemoryHotelRepository()

	ctx := context.Background()

	hotel := &Hotel{
		ID:             "HT123",
		Name:           "Hilton Garden Inn",
		Location:       "New York",
		Date:           "2025-07-20",
		HotelPrice:     25000,
		RoomsAvailable: 15,
	}
	err := repo.Create(ctx, hotel)
	require.NoError(t, err)

	t.Run("atualiza quartos disponíveis com sucesso", func(t *testing.T) {
		err := repo.UpdateAvailableRooms(ctx, "HT123", 25)
		assert.NoError(t, err)

		updated, err := repo.GetByID(ctx, "HT123")
		assert.NoError(t, err)
		assert.Equal(t, 25, updated.RoomsAvailable)
	})

	t.Run("atualiza quartos em hotel inexistente", func(t *testing.T) {
		err := repo.UpdateAvailableRooms(ctx, "HT999", 10)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "hotel not found")
	})
}

func TestHotelService_GetAll(t *testing.T) {
	repo := NewInMemoryHotelRepository()

	ctx := context.Background()

	hotels := []*Hotel{
		{ID: "HT123", Name: "Hotel A", Location: "City A", Date: "2025-07-20", HotelPrice: 10000, RoomsAvailable: 10},
		{ID: "HT124", Name: "Hotel B", Location: "City B", Date: "2025-07-21", HotelPrice: 15000, RoomsAvailable: 5},
		{ID: "HT125", Name: "Hotel C", Location: "City C", Date: "2025-07-22", HotelPrice: 20000, RoomsAvailable: 8},
	}

	for _, h := range hotels {
		err := repo.Create(ctx, h)
		require.NoError(t, err)
	}

	t.Run("lista todos os hotéis", func(t *testing.T) {
		result, err := repo.GetAll(ctx)
		assert.NoError(t, err)
		assert.Len(t, result, 3)
	})

	t.Run("lista todos os hotéis - verifica dados", func(t *testing.T) {
		result, err := repo.GetAll(ctx)
		assert.NoError(t, err)

		ids := make(map[string]bool)
		for _, h := range result {
			ids[h.ID] = true
		}
		assert.True(t, ids["HT123"])
		assert.True(t, ids["HT124"])
		assert.True(t, ids["HT125"])
	})
}
