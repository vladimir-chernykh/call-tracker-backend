
-- +migrate Up
CREATE TABLE phones(
  id serial PRIMARY KEY,
  phone varchar(16) NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL
);

-- +migrate Down
DROP TABLE phones;