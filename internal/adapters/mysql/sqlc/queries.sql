-- name: ListLuckpermsPlayers :many
SELECT
  *
FROM
  luckperms_players;

-- name: FindLuckpermsPlayerByUuid :one
SELECT
  *
FROM
  luckperms_players
WHERE
  uuid = ?;

-- name: ListMaterials :many
SELECT
  *
FROM
  materials;

-- name: ListBlocksBrokenByUser :many
SELECT
  *
FROM
  blocks
WHERE
  user = ?
AND
  action = 0;

-- name: ListBlocksPlacedByUser :many
SELECT
  *
FROM
  blocks
WHERE
  user = ?
AND
  action = 1;

-- name: FindGriefLoggerUserByUuid :one
SELECT
  *
FROM
  users
WHERE
  uuid = ?;

-- name: ListSessionDataByUser :many
SELECT
  *
FROM
  sessions
WHERE
  user = ?;