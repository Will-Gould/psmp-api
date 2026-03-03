-- +goose Up
CREATE TABLE IF NOT EXISTS users(
    id int(11) AUTO_INCREMENT PRIMARY KEY NOT NULL,
    name text NOT NULL,
    uuid text NOT NULL
);
INSERT INTO users VALUES (1, 'The_Juice_Mans', 'd91e1143-7bac-4583-a19a-ad0ea9b18971');
INSERT INTO users VALUES (2, 'monkeywarrior20', '24cefb8b-0154-4203-9c50-a1563d57affd');

-- +goose Down
DROP TABLE IF EXISTS users;
