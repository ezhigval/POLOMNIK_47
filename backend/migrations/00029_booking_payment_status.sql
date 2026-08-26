-- +goose Up
ALTER TABLE bookings
    ADD COLUMN payment_status TEXT NOT NULL DEFAULT 'NOT_REQUIRED';

UPDATE bookings
SET payment_status = CASE
    WHEN total_price > 0 THEN 'UNPAID'
    ELSE 'NOT_REQUIRED'
END;

ALTER TABLE bookings
    ADD CONSTRAINT bookings_payment_status_check CHECK (
        payment_status IN ('UNPAID', 'AWAITING_PAYMENT', 'PAID', 'NOT_REQUIRED')
    );

-- +goose Down
ALTER TABLE bookings DROP CONSTRAINT IF EXISTS bookings_payment_status_check;
ALTER TABLE bookings DROP COLUMN IF EXISTS payment_status;
