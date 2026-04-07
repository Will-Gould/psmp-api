package players

import (
	"maps"
	"net/http"
	"slices"
	"sort"

	"github.com/Will-Gould/psmp-api/internal/json"
	"github.com/go-chi/chi"
)

func (ph *playerHandler) ListServerLeaderboard(w http.ResponseWriter, r *http.Request) {
	lb := slices.Collect(maps.Values(ph.players))
	sort.Slice(lb, func(i, j int) bool {
		if lb[i].ServerRank < lb[j].ServerRank {
			return true
		}
		return false
	})
	json.Write(w, http.StatusOK, lb)
}

func (ph *playerHandler) ListCombatLeaderboard(w http.ResponseWriter, r *http.Request) {
	lb := slices.Collect(maps.Values(ph.combatLeaderboard))
	sort.Slice(lb, func(i, j int) bool {
		if lb[i].Rank < lb[j].Rank {
			return true
		}
		return false
	})
	json.Write(w, http.StatusOK, lb)
}

func (ph *playerHandler) ListCraftingLeaderboard(w http.ResponseWriter, r *http.Request) {
	lb := slices.Collect(maps.Values(ph.craftingLeaderboard))
	sort.Slice(lb, func(i, j int) bool {
		if lb[i].Rank < lb[j].Rank {
			return true
		}
		return false
	})
	json.Write(w, http.StatusOK, lb)
}

func (ph *playerHandler) ListStoryLeaderboard(w http.ResponseWriter, r *http.Request) {
	lb := slices.Collect(maps.Values(ph.storyLeaderboard))
	sort.Slice(lb, func(i, j int) bool {
		if lb[i].Rank < lb[j].Rank {
			return true
		}
		return false
	})
	json.Write(w, http.StatusOK, lb)
}

func (ph *playerHandler) GetCombatLeaderboardPlayer(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	lp, ok := ph.combatLeaderboard[uuid]
	if ok {
		json.Write(w, http.StatusOK, lp)
	} else {
		json.Write(w, http.StatusNotFound, nil)
	}
}

func (ph *playerHandler) GetCraftingLeaderboardPlayer(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	lp, ok := ph.craftingLeaderboard[uuid]
	if ok {
		json.Write(w, http.StatusOK, lp)
	} else {
		json.Write(w, http.StatusNotFound, nil)
	}
}

func (ph *playerHandler) GetStoryLeaderboardPlayer(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	lp, ok := ph.storyLeaderboard[uuid]
	if ok {
		json.Write(w, http.StatusOK, lp)
	} else {
		json.Write(w, http.StatusNotFound, nil)
	}
}
