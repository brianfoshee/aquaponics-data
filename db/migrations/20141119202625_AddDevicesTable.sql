-- +goose Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE device (
  -- id is for internal use as a primary key
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  -- Device generated identifier
  identifier character varying NOT NULL,
  -- timestamps which set to current time in UTC on record creation
  updated_at timestamp without time zone NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
  created_at timestamp without time zone NOT NULL DEFAULT (now() AT TIME ZONE 'UTC')
);

-- +goose Down
DROP TABLE device;
DROP EXTENSION "uuid-ossp";
