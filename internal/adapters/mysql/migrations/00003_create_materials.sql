-- +goose Up
CREATE TABLE IF NOT EXISTS materials(
    id int(11) AUTO_INCREMENT PRIMARY KEY NOT NULL,
    name text UNIQUE NOT NULL
);
INSERT INTO materials VALUES (1, 'dirt');
INSERT INTO materials VALUES (2, 'oak_log');
INSERT INTO materials VALUES (3, 'grass');
INSERT INTO materials VALUES (4, 'netherrack');
INSERT INTO materials VALUES (5, 'grass_block');

-- +goose Down
DROP TABLE IF EXISTS materials;
