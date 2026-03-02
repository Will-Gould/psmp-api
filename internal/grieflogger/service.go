package grieflogger

import (
	"context"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
)

type Service interface {
	ListMaterials(ctx context.Context) ([]repo.Material, error)
	CountBlocksBrokenByUser(ctx context.Context, id int32, banned []int32) (int64, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s svc) ListMaterials(ctx context.Context) ([]repo.Material, error) {
	return s.repo.ListMaterials(ctx)
}

func (s svc) ListBlocksBroken(ctx context.Context, id int32) ([]repo.Block, error) {
	return s.repo.ListBlocksBrokenByUser(ctx, id)
}

func (s svc) ListBlocksPlaced(ctx context.Context, id int32) ([]repo.Block, error) {
	return s.repo.ListBlocksPlacedByUser(ctx, id)
}

func (s svc) CountBlocksBrokenByUser(ctx context.Context, id int32, banned []int32) (int64, error) {
	return s.repo.CountBlocksBrokenByUser(ctx, id, banned)
}
