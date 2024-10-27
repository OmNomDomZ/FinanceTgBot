-- +goose Up
-- +goose StatementBegin
CREATE TABLE transactions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount NUMERIC(10, 2) NOT NULL,
    type VARCHAR(10) NOT NULL CHECK ( type IN ('доход', 'расход')),
    crated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transactions;
-- +goose StatementEnd
