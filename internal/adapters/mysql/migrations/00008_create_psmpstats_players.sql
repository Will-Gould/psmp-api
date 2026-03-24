-- +goose Up
CREATE TABLE IF NOT EXISTS psmpstats_players(
    uuid varchar(36) PRIMARY KEY,
    diamonds_mined INT DEFAULT 0 NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS psmpstats_players;
