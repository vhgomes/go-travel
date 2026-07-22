package order

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/vhgomes/go-travel/pkg/logger"
	"go.uber.org/zap"
)

type OrderHandler struct {
	svc *OrderService
}

func NewOrderHandler(svc *OrderService) *OrderHandler {
	return &OrderHandler{svc: svc}
}

// TODO: ver se esse dto está no lugar correto, ou se deveria estar em outro lugar, tipo no internal/services/order/dto.go
type createOrderRequest struct {
	UserID       string `json:"user_id"`
	FlightID     string `json:"flight_id"`
	HotelID      string `json:"hotel_id"`
	PaymentToken string `json:"payment_token"`
	TotalAmount  int64  `json:"total_amount"`
	Currency     string `json:"currency"`
}

func (h *OrderHandler) CreateOrder(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error("invalid request body", fmt.Errorf("decode error: %w", err))
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	logger.Info("create_order_request_received", zap.String("user_id", req.UserID), zap.String("flight_id", req.FlightID), zap.String("hotel_id", req.HotelID))

	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	if req.FlightID == "" || req.HotelID == "" || req.PaymentToken == "" || req.TotalAmount <= 0 || req.Currency == "" {
		http.Error(w, "missing required fields", http.StatusBadRequest)
		return
	}

	if req.Currency != string(USD) && req.Currency != string(EUR) {
		http.Error(w, "invalid currency", http.StatusBadRequest)
		return
	}

	order := Order{
		ID:           uuid.New(),
		UserID:       userID,
		FlightID:     req.FlightID,
		HotelID:      req.HotelID,
		PaymentToken: req.PaymentToken,
		TotalAmount:  req.TotalAmount,
		Currency:     Currency(req.Currency),
	}

	if err := h.svc.Create(ctx, order); err != nil {
		logger.Error("failed to create order", fmt.Errorf("service error: %w", err), zap.String("order_id", order.ID.String()), zap.String("user_id", order.UserID.String()))
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	logger.Info("order_created", zap.String("order_id", order.ID.String()), zap.String("user_id", order.UserID.String()))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(map[string]string{"id": order.ID.String()}); err != nil {
		logger.Error("failed to encode response", fmt.Errorf("encode error: %w", err), zap.String("order_id", order.ID.String()))
	}
}
