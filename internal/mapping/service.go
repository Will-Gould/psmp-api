package mapping

import (
	"context"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
)

type Service interface {
	ListMaterials(ctx context.Context) ([]repo.Material, error)
	ListMobs(ctx context.Context) ([]repo.PsmpstatsMob, error)
	ListCausesOfDeath(ctx context.Context) ([]repo.PsmpstatsCause, error)
	ListAdvancements(ctx context.Context) ([]repo.PsmpstatsAdvancement, error)
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

func (s svc) ListMobs(ctx context.Context) ([]repo.PsmpstatsMob, error) {
	return s.repo.ListMobs(ctx)
}

func (s svc) ListCausesOfDeath(ctx context.Context) ([]repo.PsmpstatsCause, error) {
	return s.repo.ListCausesOfDeath(ctx)
}

func (s svc) ListAdvancements(ctx context.Context) ([]repo.PsmpstatsAdvancement, error) {
	return s.repo.ListAdvancements(ctx)
}
