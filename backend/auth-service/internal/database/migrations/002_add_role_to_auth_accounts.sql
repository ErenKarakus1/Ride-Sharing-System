ALTER TABLE auth_accounts
ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'rider'
CHECK (role IN ('rider', 'driver'));
