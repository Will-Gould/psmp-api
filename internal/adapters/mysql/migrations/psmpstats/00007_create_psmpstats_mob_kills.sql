-- +goose Up
CREATE TABLE IF NOT EXISTS psmpstats_mob_kills(
    player_id INT NOT NULL,
    time INT NOT NULL,
    world varchar(36) NOT NULL,
    x INT NOT NULL,
    y INT NOT NULL,
    z INT NOT NULL,
    mob INT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS psmpstats_mob_kills;
