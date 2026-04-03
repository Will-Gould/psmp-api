package leaderboards

import (
	"context"
	"log/slog"
	"maps"
	"net/http"
	"slices"
	"sort"

	"github.com/Will-Gould/psmp-api/internal/json"
	"github.com/Will-Gould/psmp-api/internal/mapping"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
	"github.com/go-chi/chi"
)

type leaderboardHandler struct {
	service             Service
	serverLeaderboard   map[string]responsemodels.ServerPlayer
	combatLeaderboard   map[string]responsemodels.LeaderboardPlayer
	craftingLeaderboard map[string]responsemodels.LeaderboardPlayer
	storyLeaderboard    map[string]responsemodels.LeaderboardPlayer
}

func NewHandler(service Service) *leaderboardHandler {
	return &leaderboardHandler{
		service:             service,
		serverLeaderboard:   make(map[string]responsemodels.ServerPlayer),
		craftingLeaderboard: make(map[string]responsemodels.LeaderboardPlayer),
		combatLeaderboard:   make(map[string]responsemodels.LeaderboardPlayer),
		storyLeaderboard:    make(map[string]responsemodels.LeaderboardPlayer),
	}
}

func (lh leaderboardHandler) Initialise(ctx context.Context, md *mapping.MappingData) {
	lh.loadServerLeaderboard(ctx, md)
	lh.formRanks()
}

func (lh leaderboardHandler) ListServerLeaderboard(w http.ResponseWriter, r *http.Request) {
	json.Write(w, http.StatusOK, slices.Collect(maps.Values(lh.serverLeaderboard)))
}

func (lh leaderboardHandler) GetServerPlayer(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	sp, ok := lh.serverLeaderboard[uuid]
	if ok {
		json.Write(w, http.StatusOK, sp)
	} else {
		json.Write(w, http.StatusNotFound, nil)
	}
}

func (lh leaderboardHandler) loadServerLeaderboard(ctx context.Context, md *mapping.MappingData) {
	luckpermsPlayers, err := lh.service.ListLuckpermsPlayers(ctx)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Failed to retrieve Luckperms players")
	}

	// build overviews
	for _, p := range luckpermsPlayers {
		glUser, err := lh.service.FindGriefLoggerUser(ctx, p.Uuid)
		if err != nil {
			continue
		}
		player := &responsemodels.ServerPlayer{
			Uuid:         p.Uuid,
			Name:         glUser.Name,
			GlId:         glUser.ID,
			PrimaryGroup: p.PrimaryGroup,
		}
		lh.getPlayerData(ctx, player, md)

		lh.serverLeaderboard[player.Uuid] = *player
	}
}

func (lh leaderboardHandler) getPlayerData(ctx context.Context, player *responsemodels.ServerPlayer, md *mapping.MappingData) {

	// get combat overview
	combatOverview, err := lh.getCombatOverview(ctx, player.Uuid)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
	}

	// get crafting overview
	craftingOverview, err := lh.getCraftingOverview(ctx, player.Uuid, player.GlId, md)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
	}

	// get story overview
	storyOverview, err := lh.getStoryOverview(ctx, player.Uuid)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
	}

	player.CombatOverview = combatOverview
	player.CraftingOverview = craftingOverview
	player.StoryOverview = storyOverview

	// calculate total score
	player.Score = player.CombatOverview.CombatScore + player.CraftingOverview.CraftingScore + player.StoryOverview.StoryScore
}

func (lh leaderboardHandler) formRanks() {
	// copy player to slice and sort by score
	l := []responsemodels.ServerPlayer{}
	for _, sp := range lh.serverLeaderboard {
		l = append(l, sp)
	}
	sort.Slice(l, func(i, j int) bool {
		if l[i].Score < l[j].Score {
			return true
		}
		return false
	})

	for i, sp := range l {
		lp := lh.serverLeaderboard[sp.Uuid]
		lp.ServerRank = int64(i) + 1
		lh.serverLeaderboard[sp.Uuid] = lp
	}
}
