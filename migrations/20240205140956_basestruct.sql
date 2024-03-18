-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
  id uuid PRIMARY KEY NOT NULL UNIQUE,
  login varchar(255),
  name varchar(255),
  surname varchar(255),
  role varchar(255),
  hashed_password text,
  is_deleted bool,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp
);

CREATE INDEX user_login ON users(login);

CREATE TABLE IF NOT EXISTS presets (
  id uuid PRIMARY KEY NOT NULL UNIQUE,
  name varchar(255),
  season varchar(255),
  user_id uuid REFERENCES users (id) ON DELETE CASCADE,
  is_deleted boolean,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp
);

CREATE TABLE IF NOT EXISTS clothes (
  id uuid PRIMARY KEY NOT NULL UNIQUE,
  name varchar(255),
  type varchar(255),
  link text,
  is_deleted boolean,
  created_at timestamp,
  updated_at timestamp,
  deleted_at timestamp
);

CREATE TABLE IF NOT EXISTS clothers_presets (
  preset_id uuid REFERENCES presets (id) ON DELETE CASCADE,
  cloth_id uuid REFERENCES clothes (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS clothes;
DROP TABLE IF EXISTS presets;
DROP TABLE IF EXISTS presets_clothers;
DROP TABLE IF EXISTS users;

-- +goose StatementEnd
