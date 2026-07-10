-- +goose Up
CREATE TABLE IF NOT EXISTS sessions(
    time bigint NOT NULL,
    user int(11) NOT NULL,
    level int(11) NOT NULL,
    x int(11) NOT NULL,
    y int(11) NOT NULL,
    z int(11) NOT NULL,
    action int(11) NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS sessions;
