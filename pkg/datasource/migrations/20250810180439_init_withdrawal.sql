-- +goose Up

CREATE TABLE IF NOT EXISTS withdrawals (
    login text NOT NULL,
    order_number text NOT NULL,
    sum DECIMAL(10, 2),
    processed_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY(login, order_number)
);
-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS withdrawals;
-- +goose StatementBegin
-- +goose StatementEnd
