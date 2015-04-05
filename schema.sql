CREATE TABLE IF NOT EXISTS apps (
  id      text        primary key,
  started timestamptz NOT NULL DEFAULT now(),
  log     text        NOT NULL DEFAULT ''
);