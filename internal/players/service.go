package players

import (
	"context"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
)

type Service interface {
	FindGriefLoggerUser(ctx context.Context, uuid string) (repo.User, error)
	FindLedgerPlayer(ctx context.Context, uuid []byte) (repo.Player, error)
	GriefLoggerCountBlocksByUser(ctx context.Context, id int32, action int32, banned []int32) (int64, error)
	LedgerCountBlocksByUser(ctx context.Context, actions map[string]int32, id int32, action int32, banned []int32) (int64, error)
	LedgerListBlocksByUser(ctx context.Context, actions map[string]int32, user int32, action int32, banned []int32) ([]repo.Action, error)
	GriefLoggerListBlocksByUser(ctx context.Context, user int32, action int32, banned []int32) ([]repo.Block, error)
	ListGriefLoggerSessions(ctx context.Context, id int32) ([]repo.Session, error)
	ListPsmpstatsSessions(ctx context.Context, id int32) ([]repo.PsmpstatsSession, error)
	ListDeathsByPlayer(ctx context.Context, id int32) ([]repo.PsmpstatsDeath, error)
	CountDeathsByPlayer(ctx context.Context, id int32) (int64, error)
	CountPvpKillsByPlayer(ctx context.Context, id int32) (int64, error)
	CountMobKillsByPlayer(ctx context.Context, id int32) (int64, error)
	CountSpecificMobKillsByPlayer(ctx context.Context, name string, id int32) (int64, error)
	FindPsmpstatsPlayer(ctx context.Context, uuid string) (repo.PsmpstatsPlayer, error)
	FindLuckpermsPlayer(ctx context.Context, uuid string) (repo.LuckpermsPlayer, error)
	ListAdvancementsByPlayer(ctx context.Context, id int32) ([]repo.PsmpstatsAdvancement, error)
	ListLuckpermsPlayers(ctx context.Context) ([]repo.LuckpermsPlayer, error)
	ListMobKillsByPlayer(ctx context.Context, id int32) ([]repo.PsmpstatsMobKill, error)
	CountDiamondsMinedByPlayer(ctx context.Context, id int32) (int64, error)
	CountFishCaughtByPlayer(ctx context.Context, id int32) (int64, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s svc) FindGriefLoggerUser(ctx context.Context, uuid string) (repo.User, error) {
	return s.repo.FindGriefLoggerUserByUuid(ctx, uuid)
}

func (s svc) FindLedgerPlayer(ctx context.Context, uuid []byte) (repo.Player, error) {
	return s.repo.FindLedgerPlayerByUuid(ctx, uuid)
}

func (s svc) GriefLoggerCountBlocksByUser(ctx context.Context, id int32, action int32, banned []int32) (int64, error) {
	return s.repo.GriefLoggerCountBlocksByUser(ctx, id, action, banned)
}

func (s svc) LedgerCountBlocksByUser(ctx context.Context, actions map[string]int32, id int32, action int32, banned []int32) (int64, error) {
	return s.repo.LedgerCountBlocksByUser(ctx, actions, id, action, banned)
}

func (s svc) LedgerListBlocksByUser(ctx context.Context, actions map[string]int32, user int32, action int32, banned []int32) ([]repo.Action, error) {
	return s.repo.LedgerListBlocksByUser(ctx, actions, user, action, banned)
}

func (s svc) GriefLoggerListBlocksByUser(ctx context.Context, user int32, action int32, banned []int32) ([]repo.Block, error) {
	return s.repo.GriefLoggerListBlocksByUser(ctx, user, action, banned)
}

func (s svc) ListGriefLoggerSessions(ctx context.Context, id int32) ([]repo.Session, error) {
	return s.repo.ListGriefLoggerSessionsByUser(ctx, id)
}

func (s svc) ListPsmpstatsSessions(ctx context.Context, id int32) ([]repo.PsmpstatsSession, error) {
	return s.repo.ListPsmpstatsSessionsById(ctx, id)
}

func (s svc) ListDeathsByPlayer(ctx context.Context, id int32) ([]repo.PsmpstatsDeath, error) {
	return s.repo.ListDeathsById(ctx, id)
}

func (s svc) CountDeathsByPlayer(ctx context.Context, id int32) (int64, error) {
	return s.repo.CountDeathsById(ctx, id)
}

func (s svc) CountPvpKillsByPlayer(ctx context.Context, id int32) (int64, error) {
	return s.repo.CountPvpKillsById(ctx, id)
}

func (s svc) CountMobKillsByPlayer(ctx context.Context, id int32) (int64, error) {
	return s.repo.CountMobsKilledById(ctx, id)
}

func (s svc) CountSpecificMobKillsByPlayer(ctx context.Context, name string, id int32) (int64, error) {
	return s.repo.CountSpecificMobKillsById(ctx,
		repo.CountSpecificMobKillsByIdParams{
			PlayerID: id,
			Name:     name,
		})
}

func (s svc) FindPsmpstatsPlayer(ctx context.Context, uuid string) (repo.PsmpstatsPlayer, error) {
	return s.repo.FindPsmpstatsPlayerByUuid(ctx, uuid)
}

func (s svc) FindLuckpermsPlayer(ctx context.Context, uuid string) (repo.LuckpermsPlayer, error) {
	return s.repo.FindLuckpermsPlayerByUuid(ctx, uuid)
}

func (s svc) ListAdvancementsByPlayer(ctx context.Context, id int32) ([]repo.PsmpstatsAdvancement, error) {
	rows, err := s.repo.ListAdvancementsById(ctx, id)
	advancements := []repo.PsmpstatsAdvancement{}
	if err != nil {
		return advancements, err
	}
	// map result to PsmpstatsAdvancement struct
	for _, a := range rows {
		advancements = append(advancements, repo.PsmpstatsAdvancement{
			ID:   a.ID,
			Name: a.Name,
		})
	}
	return advancements, nil
}

func (s svc) ListLuckpermsPlayers(ctx context.Context) ([]repo.LuckpermsPlayer, error) {
	return s.repo.ListLuckpermsPlayers(ctx)
}

func (s svc) ListMobKillsByPlayer(ctx context.Context, id int32) ([]repo.PsmpstatsMobKill, error) {
	return s.repo.ListMobKillsById(ctx, id)
}

func (s svc) CountDiamondsMinedByPlayer(ctx context.Context, id int32) (int64, error) {
	return s.repo.CountDiamondsMinedById(ctx, id)
}

func (s svc) CountFishCaughtByPlayer(ctx context.Context, id int32) (int64, error) {
	return s.repo.CountFishCaughtByPlayer(ctx, id)
}
