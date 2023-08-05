BEGIN;

ALTER TABLE cashflows DROP COLUMN category_id;
DROP TABLE cashflow_categories CASCADE;

COMMIT;