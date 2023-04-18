-- +goose Up
CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  name text      NOT NULL
);
-- +goose Down
DROP TABLE authors;
