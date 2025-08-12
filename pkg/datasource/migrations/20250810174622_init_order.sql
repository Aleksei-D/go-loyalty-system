-- +goose Up
CREATE TYPE status_type AS ENUM('NEW', 'PROCESSING', 'INVALID', 'PROCESSED');
CREATE TABLE IF NOT EXISTS orders (
    number text PRIMARY KEY NOT NULL,
    login text NOT NULL,
    status status_type DEFAULT 'NEW',
    accrual DECIMAL(10, 2),
    uploaded_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    is_update_status bool DEFAULT false
);
-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
DROP TYPE IF EXISTS "status_type";
DROP TABLE IF EXISTS orders;
-- +goose StatementBegin
-- +goose StatementEnd
