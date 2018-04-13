-- +migrate Up
CREATE TABLE metric_names (
  id         SERIAL PRIMARY KEY,
  name       VARCHAR(64),
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

ALTER TABLE metric_names
  ADD CONSTRAINT metric_names_name_uniq UNIQUE (name);

INSERT INTO metric_names (name) VALUES('stt');
INSERT INTO metric_names (name) VALUES('duration');

-- +migrate Down
DROP TABLE metric_names;