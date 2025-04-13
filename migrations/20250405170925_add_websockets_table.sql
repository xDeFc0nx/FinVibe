-- +goose Up
-- +goose StatementBegin
SELECT 'up sql query';
CREATE TABLE web_sockets (
    id uuid PRIMARY KEY,
    connection_id text NOT NULL,
    user_id uuid NOT NULL REFERENCES users("id") ON DELETE CASCADE,
    is_active boolean DEFAULT false,
    last_ping TIMESTAMPTZ, 
    created_at TIMESTAMPTZ
  );
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
