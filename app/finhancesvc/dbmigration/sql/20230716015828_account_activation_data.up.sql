BEGIN;

CREATE TABLE user_activations(
    activation_token uuid PRIMARY KEY,
    user_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE,
    expired_at TIMESTAMPTZ NOT NULL
);

COMMIT;