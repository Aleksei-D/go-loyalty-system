-- +goose Up

CREATE TABLE IF NOT EXISTS withdrawals (
    login text NOT NULL,
    order_number text NOT NULL PRIMARY KEY,
    sum DECIMAL(10, 2),
    processed_at timestamptz DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS withdrawals;
-- +goose StatementBegin
-- +goose StatementEnd
