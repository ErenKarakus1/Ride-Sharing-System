ALTER TABLE auth_accounts
ADD COLUMN IF NOT EXISTS role TEXT NOT NULL DEFAULT 'rider';

DO $$
BEGIN
    ALTER TABLE auth_accounts ADD CONSTRAINT auth_accounts_role_check CHECK (role IN ('rider', 'driver'));
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;
