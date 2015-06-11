-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;
CREATE TABLE users (
  id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
  email citext UNIQUE NOT NULL,
  password varchar(64) NOT NULL
);
CREATE INDEX user_email ON users (email);

-- +goose Down
DROP TABLE users;
DROP EXTENSION citext;
