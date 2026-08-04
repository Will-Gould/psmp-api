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