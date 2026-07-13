package order

import (
	"time"

	"github.com/google/uuid"
)

type OrderService struct {
	repo OrderRepository
}

func NewOrderService(repo OrderRepository) *OrderService {
	return &OrderService{repo: repo}
}

func (s *OrderService) Create(order Order) error {
	order.ID = uuid.New()
	order.CreatedAt = time.Now()
	order.UpdatedAt = time.Now()

	return s.repo.Create(order)
}
