package order

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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
	var req createOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

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

	if err := h.svc.Create(order); err != nil {
		http.Error(w, "failed to create order", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]string{"id": order.ID.String()})
}
