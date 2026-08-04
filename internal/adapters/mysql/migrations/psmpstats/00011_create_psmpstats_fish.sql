-- +goose Up
CREATE TABLE IF NOT EXISTS psmpstats_fish(
    player_id INT NOT NULL,
    time INT NOT NULL,
    fish varchar(42) NOT NULL,
    size INT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS psmpstats_fish;