ALTER TABLE payments
ADD CONSTRAINT payments_amount_positive CHECK (amount > 0),
ADD CONSTRAINT payments_currency_length CHECK (char_length(currency) = 3);
