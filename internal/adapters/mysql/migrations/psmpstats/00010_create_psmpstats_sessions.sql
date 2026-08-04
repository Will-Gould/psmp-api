-- +goose Up
CREATE TABLE IF NOT EXISTS psmpstats_sessions(
    player_id INT NOT NULL,
    time INT NOT NULL,
    action INT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS psmpstats_sessions;