CREATE TABLE IF NOT EXISTS rides (
    id UUID PRIMARY KEY,
    rider_id UUID NOT NULL,
    driver_id UUID,
    pickup_latitude DOUBLE PRECISION NOT NULL,
    pickup_longitude DOUBLE PRECISION NOT NULL,
    pickup_address TEXT NOT NULL,
    dropoff_latitude DOUBLE PRECISION NOT NULL,
    dropoff_longitude DOUBLE PRECISION NOT NULL,
    dropoff_address TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN ('requested', 'accepted', 'started', 'completed', 'cancelled')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS rides_rider_id_idx ON rides (rider_id);
CREATE INDEX IF NOT EXISTS rides_driver_id_idx ON rides (driver_id);
CREATE INDEX IF NOT EXISTS rides_status_idx ON rides (status);
