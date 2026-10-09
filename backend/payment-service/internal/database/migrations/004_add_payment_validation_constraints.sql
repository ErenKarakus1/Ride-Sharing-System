DO $$
BEGIN
    ALTER TABLE payments ADD CONSTRAINT payments_amount_positive CHECK (amount > 0);
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

DO $$
BEGIN
    ALTER TABLE payments ADD CONSTRAINT payments_currency_length CHECK (char_length(currency) = 3);
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;
