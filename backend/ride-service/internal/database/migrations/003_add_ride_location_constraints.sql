DO $$
BEGIN
    ALTER TABLE rides ADD CONSTRAINT rides_pickup_latitude_range CHECK (pickup_latitude BETWEEN -90 AND 90);
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    ALTER TABLE rides ADD CONSTRAINT rides_pickup_longitude_range CHECK (pickup_longitude BETWEEN -180 AND 180);
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    ALTER TABLE rides ADD CONSTRAINT rides_dropoff_latitude_range CHECK (dropoff_latitude BETWEEN -90 AND 90);
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    ALTER TABLE rides ADD CONSTRAINT rides_dropoff_longitude_range CHECK (dropoff_longitude BETWEEN -180 AND 180);
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;
