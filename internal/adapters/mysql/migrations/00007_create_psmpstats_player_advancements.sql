-- +goose Up
CREATE TABLE IF NOT EXISTS psmpstats_player_advancements(
    player_uuid varchar(36),
    time INT NOT NULL,
    advancement INT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS psmpstats_player_advancements;
