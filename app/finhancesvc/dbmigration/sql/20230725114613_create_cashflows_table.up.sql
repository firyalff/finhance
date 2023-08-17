BEGIN;

CREATE TYPE CASHFLOW_CATEGORY_TYPE_ENUM AS ENUM ('income', 'expense');

CREATE TABLE cashflow_categories (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    cashflow_category_type CASHFLOW_CATEGORY_TYPE_ENUM NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NULL,
    deleted_at TIMESTAMPTZ DEFAULT NULL,
    UNIQUE (user_id, name)
);

CREATE TABLE cashflows (
    id uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    category_id uuid NOT NULL REFERENCES cashflow_categories(id) ON DELETE CASCADE ON UPDATE CASCADE,
    name TEXT NOT NULL,
    notes TEXT,
    amount INT8 NOT NULL,
    proof_document_url TEXT DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NULL
);

COMMIT;