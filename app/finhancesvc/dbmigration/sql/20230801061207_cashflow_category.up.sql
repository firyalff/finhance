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

ALTER TABLE cashflows
    ADD category_id uuid REFERENCES cashflow_categories(id) ON DELETE CASCADE ON UPDATE CASCADE; 

COMMIT;