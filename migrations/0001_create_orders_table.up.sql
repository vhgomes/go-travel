CREATE TABLE orders (
    id              UUID PRIMARY KEY,
    user_id         UUID NOT NULL,
    flight_id       TEXT NOT NULL,
    hotel_id        TEXT NOT NULL,
    payment_token   TEXT NOT NULL,
    total_amount    DECIMAL(10,2) NOT NULL,
    currency        TEXT NOT NULL,
    status          TEXT NOT NULL CHECK (status IN ('PENDING', 'FLIGHT_OK', 'HOTEL_OK', 'CONFIRMED', 'ROLLBACK', 'FAILED')),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id ON orders(user_id);
CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);