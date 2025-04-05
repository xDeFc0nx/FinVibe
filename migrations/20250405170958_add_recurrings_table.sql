-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE recurrings (
    id UUID PRIMARY KEY,
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    amount DECIMAL(15,2) NOT NULL,
    frequency TEXT NOT NULL,
    start_date TIMESTAMPTZ NOT NULL,
    next_date TIMESTAMPTZ NOT NULL,
    end_date TIMESTAMPTZ  NOT NULL,
    created_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
