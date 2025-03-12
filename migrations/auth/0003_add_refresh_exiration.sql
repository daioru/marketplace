-- +goose Up
ALTER TABLE users ADD COLUMN refresh_expires_at TIMESTAMP;

-- +goose Down
ALTER TABLE users DROP COLUMN refresh_expires_at;
