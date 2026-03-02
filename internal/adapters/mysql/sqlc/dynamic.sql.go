package repo

import (
	"context"
	"strings"
)

const CountBlocksBrokenByUser = `-- name: CountBlocksBrokenByUser :many
SELECT
  count(*)
FROM
  blocks
WHERE
  user = ?
AND
  action = 0
`

func (q *Queries) CountBlocksBrokenByUser(ctx context.Context, user int32, banned []int32) (int64, error) {
	// create args
	args := []any{}
	args = append(args, user)
	for i := range banned {
		args = append(args, banned[i])
	}

	// build dynamic query
	var query strings.Builder
	query.WriteString(CountBlocksBrokenByUser)
	for range banned {
		query.WriteString(" AND type != ?")
	}

	row := q.db.QueryRowContext(ctx, query.String(), args...)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const CountBlocksPlacedByUser = `-- name: CountBlocksBrokenByUser :many
SELECT
  count(*)
FROM
  blocks
WHERE
  user = ?
AND
  action = 1
`

func (q *Queries) CountBlocksPlacedByUser(ctx context.Context, user int32, banned []int32) (int64, error) {
	row := q.db.QueryRowContext(ctx, CountBlocksPlacedByUser, user)
	var count int64
	err := row.Scan(&count)
	return count, err
}
