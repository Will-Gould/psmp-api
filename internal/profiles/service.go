package profiles

import (
	"context"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
)

type Service interface {
	FindGriefLoggerUser(ctx context.Context, uuid string) (repo.User, error)
	LedgerListBlocksByUser(ctx context.Context, actions map[string]int32, id int32, action int32, banned []int32) ([]repo.Action, error)
	GriefLoggerListBlocksByUser(ctx context.Context, id int32, action int32, banned []int32) ([]repo.Block, error)
	ListGriefLoggerSessions(ctx context.Context, id int32) ([]repo.Session, error)
	ListPsmpstatsSessions(ctx context.Context, id int32) ([]repo.PsmpstatsSession, error)
	ListDeathsByPlayer(ctx context.Context, id int32) ([]repo.PsmpstatsDeath, error)
	CountSpecificMobKillsByPlayer(ctx context.Context, name string, id int32) (int64, error)
	FindLedgerPlayer(ctx context.Context, uuid []byte) (repo.Player, error)
	FindPsmpstatsPlayerByUuid(ctx context.Context, uuid string) (repo.PsmpstatsPlayer, error)
	FindLuckpermsPlayer(ctx context.Context, uuid string) (repo.LuckpermsPlayer, error)
	ListAdvancementsByPlayer(ctx context.Context, id int32) ([]repo.PsmpstatsAdvancement, error)
	ListMobKillsByPlayer(ctx context.Context, id int32) ([]repo.PsmpstatsMobKill, error)
	GriefLoggerGroupCountBlocksByUser(ctx context.Context, id int32, action int32, banned []int32) ([]repo.GroupCountBlocksPlacedByUserRow, error)
	LedgerGroupCountBlocksByUser(ctx context.Context, actions map[string]int32, id int32, action int32, banned []int32) ([]repo.GroupCountBlocksPlacedByUserRow, error)
	FindGriefLoggerFirstJoin(ctx context.Context, id int32) (repo.Session, error)
	FindPsmpstatsFirstJoin(ctx context.Context, id int32) (repo.PsmpstatsSession, error)
	GroupCountMobKillsByPlayer(ctx context.Context, id int32) ([]repo.GroupCountMobKillsByPlayerRow, error)
}

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier, gameLogger string) Service {
	return &svc{repo: repo}
}

// CountSpecificMobKillsByUuid implements [Service].
func (s *svc) CountSpecificMobKillsByPlayer(ctx context.Context, name string, id int32) (int64, error) {
	return s.repo.CountSpecificMobKillsById(ctx, repo.CountSpecificMobKillsByIdParams{
		PlayerID: id,
		Name:     name,
	})
}

// FindGriefLoggerUser implements [Service].
func (s *svc) FindGriefLoggerUser(ctx context.Context, uuid string) (repo.User, error) {
	return s.repo.FindGriefLoggerUserByUuid(ctx, uuid)
}

func (s *svc) FindLedgerPlayer(ctx context.Context, uuid []byte) (repo.Player, error) {
	return s.repo.FindLedgerPlayerByUuid(ctx, uuid)
}

// FindLuckpermsPlayer implements [Service].
func (s *svc) FindLuckpermsPlayer(ctx context.Context, uuid string) (repo.LuckpermsPlayer, error) {
	return s.repo.FindLuckpermsPlayerByUuid(ctx, uuid)
}

// FindPsmpstatsPlayerByUuid implements [Service].
func (s *svc) FindPsmpstatsPlayerByUuid(ctx context.Context, uuid string) (repo.PsmpstatsPlayer, error) {
	return s.repo.FindPsmpstatsPlayerByUuid(ctx, uuid)
}

// ListAdvancementsByPlayer implements [Service].
func (s *svc) ListAdvancementsByPlayer(ctx context.Context, id int32) ([]repo.PsmpstatsAdvancement, error) {
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

// ListDeathsByPlayer implements [Service].
func (s *svc) ListDeathsByPlayer(ctx context.Context, id int32) ([]repo.PsmpstatsDeath, error) {
	return s.repo.ListDeathsById(ctx, id)
}

// ListMobKillsByUuid implements [Service].
func (s *svc) ListMobKillsByPlayer(ctx context.Context, id int32) ([]repo.PsmpstatsMobKill, error) {
	return s.repo.ListMobKillsById(ctx, id)
}

// ListSessionDataByUser implements [Service].
func (s *svc) ListGriefLoggerSessions(ctx context.Context, id int32) ([]repo.Session, error) {
	return s.repo.ListGriefLoggerSessionsByUser(ctx, id)
}

func (s *svc) ListPsmpstatsSessions(ctx context.Context, id int32) ([]repo.PsmpstatsSession, error) {
	return s.repo.ListPsmpstatsSessionsById(ctx, id)
}

func (s *svc) GriefLoggerGroupCountBlocksByUser(ctx context.Context, id int32, action int32, banned []int32) ([]repo.GroupCountBlocksPlacedByUserRow, error) {
	return s.repo.GriefLoggerGroupCountBlocksByUser(ctx, id, action, banned)
}

func (s *svc) LedgerGroupCountBlocksByUser(ctx context.Context, actions map[string]int32, id int32, action int32, banned []int32) ([]repo.GroupCountBlocksPlacedByUserRow, error) {
	return s.repo.LedgerGroupCountBlocksByUser(ctx, actions, id, action, banned)
}

func (s *svc) GriefLoggerListBlocksByUser(ctx context.Context, id int32, action int32, banned []int32) ([]repo.Block, error) {
	return s.repo.GriefLoggerListBlocksByUser(ctx, id, action, banned)
}

func (s *svc) LedgerListBlocksByUser(ctx context.Context, actions map[string]int32, id int32, action int32, banned []int32) ([]repo.Action, error) {
	return s.repo.LedgerListBlocksByUser(ctx, actions, id, action, banned)
}

func (s *svc) FindGriefLoggerFirstJoin(ctx context.Context, id int32) (repo.Session, error) {
	return s.repo.FindGriefLoggerFirstJoinByUser(ctx, id)
}

func (s *svc) FindPsmpstatsFirstJoin(ctx context.Context, id int32) (repo.PsmpstatsSession, error) {
	return s.repo.FindPsmpstatsFirstJoinById(ctx, id)
}

func (s *svc) GroupCountMobKillsByPlayer(ctx context.Context, id int32) ([]repo.GroupCountMobKillsByPlayerRow, error) {
	return s.repo.GroupCountMobKillsByPlayer(ctx, id)
}
