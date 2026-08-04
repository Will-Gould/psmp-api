-- +goose Up
CREATE TABLE IF NOT EXISTS ActionIdentifiers(
    id int(11) PRIMARY KEY NOT NULL,
    action_identifier varchar(16) UNIQUE NOT NULL
);