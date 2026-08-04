-- +goose Up
CREATE TABLE IF NOT EXISTS sources(
    id int(11) PRIMARY KEY NOT NULL,
    identifier varchar(191) UNIQUE NOT NULL
);