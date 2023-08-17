BEGIN;

CREATE TABLE cashflows (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    name TEXT NOT NULL,
    notes TEXT,
    amount INT8 NOT NULL,
    proof_document_url TEXT DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NULL
);

COMMIT;