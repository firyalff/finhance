BEGIN;

CREATE TYPE CASHFLOW_TYPE_ENUM AS ENUM ('income', 'expense');

CREATE TABLE cashflows (
    id uuid PRIMARY KEY,
    cashflow_type CASHFLOW_TYPE_ENUM NOT NULL,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    amount INT8 NOT NULL,
    proof_document_url varchar(200) DEFAULT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NULL
);

COMMIT;