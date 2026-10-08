CREATE TABLE IF NOT EXISTS payments (
    id UUID PRIMARY KEY,
    ride_id UUID NOT NULL,
    rider_id UUID NOT NULL,
    driver_id UUID,
    amount NUMERIC(12, 2) NOT NULL,
    currency TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('authorized', 'captured', 'refunded', 'failed')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS payments_ride_id_idx ON payments (ride_id);
CREATE INDEX IF NOT EXISTS payments_rider_id_idx ON payments (rider_id);
CREATE INDEX IF NOT EXISTS payments_driver_id_idx ON payments (driver_id);
CREATE INDEX IF NOT EXISTS payments_status_idx ON payments (status);
