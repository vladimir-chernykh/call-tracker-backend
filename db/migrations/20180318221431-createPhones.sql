
-- +migrate Up
CREATE TABLE phones(
  id serial PRIMARY KEY,
  number varchar(16) NOT NULL UNIQUE,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL
);

-- +migrate Down
DROP TABLE phones;