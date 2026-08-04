-- +goose Up
CREATE TABLE IF NOT EXISTS psmpstats_players(
    id INT PRIMARY KEY,
    uuid varchar(36) UNIQUE NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS psmpstats_players;
