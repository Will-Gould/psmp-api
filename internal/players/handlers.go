package players

import (
	"context"
	"log/slog"
	"maps"
	"net/http"
	"slices"
	"sort"
	"sync"

	"github.com/Will-Gould/psmp-api/internal/json"
	"github.com/Will-Gould/psmp-api/internal/mapping"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
	"github.com/go-chi/chi"
)

type playerHandler struct {
	mu                  sync.RWMutex
	service             Service
	players             map[string]responsemodels.ServerPlayer
	combatLeaderboard   map[string]responsemodels.LeaderboardPlayer
	craftingLeaderboard map[string]responsemodels.LeaderboardPlayer
	storyLeaderboard    map[string]responsemodels.LeaderboardPlayer
}

func NewHandler(service Service) *playerHandler {
	return &playerHandler{
		service:             service,
		players:             make(map[string]responsemodels.ServerPlayer),
		craftingLeaderboard: make(map[string]responsemodels.LeaderboardPlayer),
		combatLeaderboard:   make(map[string]responsemodels.LeaderboardPlayer),
		storyLeaderboard:    make(map[string]responsemodels.LeaderboardPlayer),
	}
}

func (ph *playerHandler) Load(ctx context.Context, md *mapping.MappingData) {
	ph.mu.Lock()
	ph.loadServerLeaderboard(ctx, md)
	ph.formServerRanks()
	ph.formCombatRanks()
	ph.formCraftingRanks()
	ph.formStoryRanks()
	ph.mu.Unlock()
}

func (ph *playerHandler) ListServerLeaderboard(w http.ResponseWriter, r *http.Request) {
	json.Write(w, http.StatusOK, slices.Collect(maps.Values(ph.players)))
}

func (ph *playerHandler) ListCombatLeaderboard(w http.ResponseWriter, r *http.Request) {
	json.Write(w, http.StatusOK, slices.Collect(maps.Values(ph.combatLeaderboard)))
}

func (ph *playerHandler) ListCraftingLeaderboard(w http.ResponseWriter, r *http.Request) {
	json.Write(w, http.StatusOK, slices.Collect(maps.Values(ph.craftingLeaderboard)))
}

func (ph *playerHandler) ListStoryLeaderboard(w http.ResponseWriter, r *http.Request) {
	json.Write(w, http.StatusOK, slices.Collect(maps.Values(ph.storyLeaderboard)))
}

func (ph *playerHandler) GetServerPlayer(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	sp, ok := ph.players[uuid]
	if ok {
		json.Write(w, http.StatusOK, sp)
	} else {
		json.Write(w, http.StatusNotFound, nil)
	}
}

func (ph *playerHandler) loadServerLeaderboard(ctx context.Context, md *mapping.MappingData) {
	luckpermsPlayers, err := ph.service.ListLuckpermsPlayers(ctx)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Failed to retrieve Luckperms players")
	}

	// build overviews
	for _, p := range luckpermsPlayers {
		glUser, err := ph.service.FindGriefLoggerUser(ctx, p.Uuid)
		if err != nil {
			continue
		}
		player := &responsemodels.ServerPlayer{
			Uuid:         p.Uuid,
			Name:         glUser.Name,
			GlId:         glUser.ID,
			PrimaryGroup: p.PrimaryGroup,
		}
		ph.getPlayerData(ctx, player, md)

		ph.players[player.Uuid] = *player
	}
}

func (ph *playerHandler) getPlayerData(ctx context.Context, player *responsemodels.ServerPlayer, md *mapping.MappingData) {

	// get combat overview
	combatOverview, err := ph.getCombatOverview(ctx, player.Uuid)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
	}

	// get crafting overview
	craftingOverview, err := ph.getCraftingOverview(ctx, player.Uuid, player.GlId, md)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
	}

	// get story overview
	storyOverview, err := ph.getStoryOverview(ctx, player.Uuid)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
	}

	player.CombatOverview = combatOverview
	player.CraftingOverview = craftingOverview
	player.StoryOverview = storyOverview

	// calculate total score
	player.Score = player.CombatOverview.CombatScore + player.CraftingOverview.CraftingScore + player.StoryOverview.StoryScore
}

func (ph *playerHandler) formServerRanks() {
	// get player slice and sort by score
	l := slices.Collect(maps.Values(ph.players))
	sort.Slice(l, func(i, j int) bool {
		if l[i].Score < l[j].Score {
			return true
		}
		return false
	})
	for i, sp := range l {
		lp := ph.players[sp.Uuid]
		lp.ServerRank = int64(i) + 1
		ph.players[sp.Uuid] = lp
	}
}

func (ph *playerHandler) formCombatRanks() {
	// copy player to slice and sort by combat score
	l := slices.Collect(maps.Values(ph.players))
	sort.Slice(l, func(i, j int) bool {
		if l[i].CombatOverview.CombatScore < l[j].CombatOverview.CombatScore {
			return true
		}
		return false
	})

	for i, sp := range l {
		lp := responsemodels.LeaderboardPlayer{
			Uuid:    sp.Uuid,
			Ranking: int64(i) + 1,
			Value:   sp.CombatOverview.CombatScore,
		}
		ph.combatLeaderboard[sp.Uuid] = lp
	}
}

func (ph *playerHandler) formCraftingRanks() {
	// copy player to slice and sort by crafting score
	l := slices.Collect(maps.Values(ph.players))
	sort.Slice(l, func(i, j int) bool {
		if l[i].CraftingOverview.CraftingScore < l[j].CraftingOverview.CraftingScore {
			return true
		}
		return false
	})

	for i, sp := range l {
		lp := responsemodels.LeaderboardPlayer{
			Uuid:    sp.Uuid,
			Ranking: int64(i) + 1,
			Value:   sp.CraftingOverview.CraftingScore,
		}
		ph.craftingLeaderboard[sp.Uuid] = lp
	}
}

func (ph *playerHandler) formStoryRanks() {
	// copy player to slice and sort by story score
	l := slices.Collect(maps.Values(ph.players))
	sort.Slice(l, func(i, j int) bool {
		if l[i].StoryOverview.StoryScore < l[j].StoryOverview.StoryScore {
			return true
		}
		return false
	})

	for i, sp := range l {
		lp := responsemodels.LeaderboardPlayer{
			Uuid:    sp.Uuid,
			Ranking: int64(i) + 1,
			Value:   sp.StoryOverview.StoryScore,
		}
		ph.storyLeaderboard[sp.Uuid] = lp
	}
}
