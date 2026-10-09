CREATE TABLE IF NOT EXISTS payment_outbox_events (
    id UUID PRIMARY KEY,
    event_type TEXT NOT NULL,
    event_key TEXT NOT NULL,
    payload JSONB NOT NULL,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'sent', 'failed')),
    attempts INTEGER NOT NULL DEFAULT 0,
    published_at TIMESTAMPTZ,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS payment_outbox_events_status_idx ON payment_outbox_events (status, created_at);
