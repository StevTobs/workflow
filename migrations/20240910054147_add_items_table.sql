-- +goose Up
CREATE TABLE items (
    id bigserial NOT NULL PRIMARY KEY,
    title text NOT NULL,
    amount integer NOT NULL,
    quantity integer NOT NULL,
    status text NOT NULL,
    owner_id bigserial NOT NULL
);
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS items;
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd