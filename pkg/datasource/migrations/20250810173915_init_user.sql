-- +goose Up

CREATE TABLE IF NOT EXISTS "users" (
    login text NOT NULL,
    password text NOT NULL,
    PRIMARY KEY(login)
);
-- +goose StatementBegin

-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS "users";
-- +goose StatementBegin
-- +goose StatementEnd
