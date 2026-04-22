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

-- name: GroupCountBlocksPlacedByUser :many
SELECT
  type, COUNT(*) AS total_placed
FROM
  blocks
WHERE
  user = ?
AND
  action = 1
GROUP BY
  type;

-- name: FindFirstJoinByUser :one
SELECT
  *
FROM
  sessions
WHERE
  user = ?
ORDER BY
  time
ASC
LIMIT 1;