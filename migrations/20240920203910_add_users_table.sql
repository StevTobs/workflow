-- +goose Up
CREATE TABLE users (
    id bigserial NOT NULL PRIMARY KEY,
    username text NOT NULL UNIQUE,
    password text NOT NULL
);

ALTER TABLE items
ADD CONSTRAINT fk_owner
FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE;

-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

-- +goose Down
ALTER TABLE items
DROP CONSTRAINT fk_owner;

DROP TABLE IF EXISTS users;

-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
