-- +goose Up
CREATE TABLE IF NOT EXISTS psmpstats_combat(
    player_id INT NOT NULL,
    victim_id INT NOT NULL,
    time INT NOT NULL,
    world varchar(36) NOT NULL,
    x INT NOT NULL,
    y INT NOT NULL,
    z INT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS psmpstats_combat;
