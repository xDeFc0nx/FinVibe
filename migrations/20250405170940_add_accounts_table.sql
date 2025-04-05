-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE accounts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type TEXT NOT NULL,
    income DECIMAL(15,2),
    expense DECIMAL(15,2),
    balance DECIMAL(15,2),
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);


-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
