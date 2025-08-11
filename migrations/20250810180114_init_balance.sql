-- +goose Up

CREATE TABLE IF NOT EXISTS balance (
    login text PRIMARY KEY NOT NULL,
    current DECIMAL(10, 2) DEFAULT 0.00
);
-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS balance;
-- +goose StatementBegin
-- +goose StatementEnd
