-- +goose Up
CREATE TABLE IF NOT EXISTS players(
    id int(11) PRIMARY KEY NOT NULL,
    player_id binary(16) UNIQUE NOT NULL,
    player_name varchar(16) NOT NULL,
    first_join timestamp(6) NOT NULL,
    last_join timestamp(6) NOT NULL
);