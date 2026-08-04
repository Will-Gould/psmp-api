-- +goose Up
CREATE TABLE IF NOT EXISTS sources(
    id int(11) PRIMARY KEY NOT NULL,
    name varchar(30) UNIQUE NOT NULL
);