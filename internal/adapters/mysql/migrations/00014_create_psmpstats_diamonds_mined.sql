-- +goose Up
CREATE TABLE IF NOT EXISTS psmpstats_diamonds_mined(
    player_id INT NOT NULL,
    time INT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS psmpstats_diamonds_mined;
