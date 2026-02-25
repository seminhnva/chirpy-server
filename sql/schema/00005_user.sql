-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
ADD COLUMN is_chirpy_red BOOL NOT NULL DEFAULT FALSE;

-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
ALTER TABLE users
DROP COLUMN is_chirpy_red;

-- +goose StatementEnd