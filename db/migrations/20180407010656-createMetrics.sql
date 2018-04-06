-- +migrate Up
CREATE TABLE metrics (
  id         SERIAL PRIMARY KEY,
  name       VARCHAR(64),
  call       INTEGER REFERENCES calls,
  remote_id  VARCHAR(64),
  data       JSONB,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL
);

-- +migrate Down
DROP TABLE metrics;