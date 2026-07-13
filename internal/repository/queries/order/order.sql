-- name: CreateOrder :one
INSERT INTO orders (
    id,
    user_id,
    flight_id,
    hotel_id,
    payment_token,
    total_amount,
    currency,
    status
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, 'PENDING'
)
RETURNING *;

-- name: GetOrderByID :one
SELECT * FROM orders
WHERE id = $1;

-- name: GetOrdersByUserID :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateOrderStatus :one
UPDATE orders
SET 
    status = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: UpdateOrderPaymentToken :one
UPDATE orders
SET 
    payment_token = $2,
    updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: CheckPendingOrderExists :one
SELECT COUNT(*) as count
FROM orders
WHERE user_id = $1
  AND flight_id = $2
  AND hotel_id = $3
  AND status IN ('PENDING', 'FLIGHT_OK', 'HOTEL_OK');

-- name: GetOrdersByStatus :many
SELECT * FROM orders
WHERE status = $1
ORDER BY created_at ASC
LIMIT $2 OFFSET $3;