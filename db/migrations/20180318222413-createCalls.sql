
-- +migrate Up
CREATE TABLE calls(
  id serial PRIMARY KEY,
  phone_id integer REFERENCES phones,
  record bytea,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL
);

-- +migrate Down
DROP TABLE calls;