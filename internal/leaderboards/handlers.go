package leaderboards

import (
	"net/http"

	"github.com/Will-Gould/psmp-api/internal/json"
	"github.com/Will-Gould/psmp-api/internal/mapping"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
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

func (lh leaderboardHandler) Initialise(*mapping.MappingData) {

}

func (lh leaderboardHandler) ListServerLeaderboard(w http.ResponseWriter, r *http.Request) {
	json.Write(w, http.StatusOK, lh.serverLeaderboard)
}
