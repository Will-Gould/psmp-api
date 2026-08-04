-- +goose Up
CREATE TABLE IF NOT EXISTS luckperms_players (
    uuid varchar(36) PRIMARY KEY NOT NULL,
    username varchar(16) NOT NULL,
    primary_group varchar(36) NOT NULL
);
INSERT INTO luckperms_players VALUES ('d91e1143-7bac-4583-a19a-ad0ea9b18971', 'the_juice_mans', 'owner');
INSERT INTO luckperms_players VALUES ('24cefb8b-0154-4203-9c50-a1563d57affd', 'monkeywarrior20', 'admin');

-- +goose Down
DROP TABLE IF EXISTS luckperms_players;
