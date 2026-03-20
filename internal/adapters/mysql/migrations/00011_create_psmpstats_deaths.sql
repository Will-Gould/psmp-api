-- +goose Up
CREATE TABLE IF NOT EXISTS psmpstats_deaths(
    player_uuid varchar(36) NOT NULL,
    time INT NOT NULL,
    world varchar(36) NOT NULL,
    x INT NOT NULL,
    y INT NOT NULL,
    z INT NOT NULL,
    cause INT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS psmpstats_deaths;
