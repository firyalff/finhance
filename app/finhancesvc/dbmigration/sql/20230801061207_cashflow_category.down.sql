BEGIN;

ALTER TABLE cashflows DROP COLUMN category_id;
DROP TABLE cashflow_categories CASCADE;

DROP TYPE CASHFLOW_CATEGORY_TYPE_ENUM CASCADE;


COMMIT;