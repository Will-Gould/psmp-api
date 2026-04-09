package repo

import (
	"context"
	"strings"
)

const CountBlocksByUser = `-- name: CountBlocksBrokenByUser :one
SELECT
  count(*)
FROM
  blocks
WHERE
  user = ?
AND
  action = ?
`

func (q *Queries) CountBlocksByUser(ctx context.Context, user int32, action int32, banned []int32) (int64, error) {
	// create args
	args := []any{}
	args = append(args, user)
	args = append(args, action)
	for i := range banned {
		args = append(args, banned[i])
	}

	// build dynamic query
	var query strings.Builder
	query.WriteString(CountBlocksByUser)
	for range banned {
		query.WriteString(" AND type != ?")
	}

	row := q.db.QueryRowContext(ctx, query.String(), args...)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const GroupCountBlocksByUser = `-- name: GroupCountBlocksByUSer :many
SELECT
  type, COUNT(*) AS total_placed
FROM
  blocks
WHERE
  user = ?
AND
  action = ?
`

func (q *Queries) GroupCountBlocksByUser(ctx context.Context, user int32, action int32, banned []int32) ([]GroupCountBlocksPlacedByUserRow, error) {
	// create args
	args := []any{}
	args = append(args, user)
	args = append(args, action)
	for i := range banned {
		args = append(args, banned[i])
	}

	// build dynamic query
	var query strings.Builder
	query.WriteString(GroupCountBlocksByUser)
	for range banned {
		query.WriteString(" AND type != ?")
	}

	query.WriteString(" GROUP BY type")

	rows, err := q.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GroupCountBlocksPlacedByUserRow
	for rows.Next() {
		var i GroupCountBlocksPlacedByUserRow
		if err := rows.Scan(&i.Type, &i.TotalPlaced); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const ListBlocksByUser = `-- name: ListBlocksByUser :many
SELECT
  time, user, level, x, y, z, type, action
FROM
  blocks
WHERE
  user = ?
AND
  action = ?
`

func (q *Queries) ListBlocksByUser(ctx context.Context, user int32, action int32, banned []int32) ([]Block, error) {
	// create args
	args := []any{}
	args = append(args, user)
	args = append(args, action)
	for i := range banned {
		args = append(args, banned[i])
	}

	// build dynamic query
	var query strings.Builder
	query.WriteString(ListBlocksByUser)
	for range banned {
		query.WriteString(" AND type != ?")
	}

	rows, err := q.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Block
	for rows.Next() {
		var i Block
		if err := rows.Scan(
			&i.Time,
			&i.User,
			&i.Level,
			&i.X,
			&i.Y,
			&i.Z,
			&i.Type,
			&i.Action,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
