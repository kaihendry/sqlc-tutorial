-- +goose Up
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  name text      NOT NULL,
  bio  text -- why would we want this to be nullable?
);
-- +goose Down
DROP TABLE authors;
