package mapping

import (
	"context"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
)

type Service interface {
	ListLedgerObjects(ctx context.Context) ([]repo.Objectidentifier, error)
	ListGriefLoggerMaterials(ctx context.Context) ([]repo.Material, error)
	ListMobs(ctx context.Context) ([]repo.PsmpstatsMob, error)
	ListCausesOfDeath(ctx context.Context) ([]repo.PsmpstatsCause, error)
	ListAdvancements(ctx context.Context) ([]repo.PsmpstatsAdvancement, error)
	ListActionIdentifiers(ctx context.Context) ([]repo.Actionidentifier, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s svc) ListLedgerObjects(ctx context.Context) ([]repo.Objectidentifier, error) {
	return s.repo.ListObjects(ctx)
}

func (s svc) ListGriefLoggerMaterials(ctx context.Context) ([]repo.Material, error) {
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

func (s svc) ListActionIdentifiers(ctx context.Context) ([]repo.Actionidentifier, error) {
	return s.repo.ListActionsIdentifiers(ctx)
}
