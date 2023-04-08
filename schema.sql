CREATE TABLE authors (
  id   BIGSERIAL PRIMARY KEY,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  name text      NOT NULL,
  bio  text
);
