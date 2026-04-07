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
	combatLeaderboard   map[string]responsemodels.CombatLeaderboardPlayer
	craftingLeaderboard map[string]responsemodels.CraftingLeaderboardPlayer
	storyLeaderboard    map[string]responsemodels.StoryLeaderboardPlayer
}

func NewHandler(service Service) *playerHandler {
	return &playerHandler{
		service:             service,
		players:             make(map[string]responsemodels.ServerPlayer),
		craftingLeaderboard: make(map[string]responsemodels.CraftingLeaderboardPlayer),
		combatLeaderboard:   make(map[string]responsemodels.CombatLeaderboardPlayer),
		storyLeaderboard:    make(map[string]responsemodels.StoryLeaderboardPlayer),
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

func (ph *playerHandler) GetServerPlayer(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	sp, ok := ph.players[uuid]
	if ok {
		json.Write(w, http.StatusOK, sp)
	} else {
		json.Write(w, http.StatusNotFound, nil)
	}
}

func (ph *playerHandler) GetChampionPlayer(w http.ResponseWriter, r *http.Request) {
	for _, p := range ph.players {
		if p.ServerRank == 1 {
			json.Write(w, http.StatusOK, p)
			return
		}
	}
	json.Write(w, http.StatusInternalServerError, nil)
}

func (ph *playerHandler) GetMostDangerousPlayer(w http.ResponseWriter, r *http.Request) {
	kdRanking := slices.Collect(maps.Values(ph.players))
	sort.Slice(kdRanking, func(i, j int) bool {
		if kdRanking[i].CombatOverview.PvpKdRatio < kdRanking[j].CombatOverview.PvpKdRatio {
			return true
		}
		return false
	})

	json.Write(w, http.StatusOK, kdRanking[0])
}

func (ph *playerHandler) GetBiggestBuilder(w http.ResponseWriter, r *http.Request) {
	buildRanking := slices.Collect(maps.Values(ph.players))
	sort.Slice(buildRanking, func(i, j int) bool {
		if buildRanking[i].CraftingOverview.BlocksPlaced < buildRanking[j].CraftingOverview.BlocksPlaced {
			return true
		}
		return false
	})

	json.Write(w, http.StatusOK, buildRanking[0])
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
		lp := responsemodels.CombatLeaderboardPlayer{
			Uuid:         sp.Uuid,
			Name:         sp.Name,
			PrimaryGroup: sp.PrimaryGroup,
			Rank:         int64(i) + 1,
			PvpKdRatio:   sp.CombatOverview.PvpKdRatio,
			PvpKills:     sp.CombatOverview.PvpKills,
			MobKills:     sp.CombatOverview.MobKills,
			Deaths:       sp.CombatOverview.Deaths,
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
		lp := responsemodels.CraftingLeaderboardPlayer{
			Uuid:          sp.Uuid,
			Name:          sp.Name,
			PrimaryGroup:  sp.PrimaryGroup,
			Rank:          int64(i) + 1,
			BlocksPlaced:  sp.CraftingOverview.BlocksPlaced,
			BlocksBroken:  sp.CraftingOverview.BlocksBroken,
			DiamondsMined: sp.CraftingOverview.DiamondsMined,
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
		lp := responsemodels.StoryLeaderboardPlayer{
			Uuid:         sp.Uuid,
			Name:         sp.Name,
			PrimaryGroup: sp.PrimaryGroup,
			Rank:         int64(i) + 1,
		}
		ph.storyLeaderboard[sp.Uuid] = lp
	}
}
