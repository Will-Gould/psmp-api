-- +goose Up
CREATE TABLE IF NOT EXISTS ObjectIdentifiers(
    id int(11) PRIMARY KEY NOT NULL,
    identifier varchar(191) UNIQUE NOT NULL
);