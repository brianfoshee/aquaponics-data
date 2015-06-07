-- +goose Up
ALTER TABLE device ADD COLUMN user_id uuid NOT NULL REFERENCES users (id) ON DELETE CASCADE;

-- +goose Down
ALTER TABLE device DROP COLUMN IF EXISTS user_id;