-- +goose Up
CREATE TABLE reading (
  device_id uuid NOT NULL REFERENCES device (id) ON DELETE CASCADE,
  readings json NOT NULL DEFAULT '{}'::json
);
CREATE INDEX reading_device_id ON reading (device_id);

-- +goose Down
DROP TABLE reading;
