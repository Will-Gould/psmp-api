package players

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Will-Gould/psmp-api/internal/json"
	"github.com/go-chi/chi"
)

type handler struct {
	service Service
}

type Player struct {
	Uuid          string
	PrimaryGroup  string
	GriefLoggerId int32
	Username      string
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
}

func (h handler) ListPlayersHandler(w http.ResponseWriter, r *http.Request) {
	players, err := h.service.ListPlayers(r.Context())
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, players)
}

func (h handler) getPlayer(ctx context.Context, uuid string) (Player, error) {
	player := Player{}

	// get luckperms data
	luckpermsPlayer, err := h.service.FindLuckpermsPlayer(ctx, uuid)
	if err != nil {
		fmt.Printf("Error retrieving luckperms player with uuid: %v\n", uuid)
		return player, err
	}

	// get grief logger user
	griefLoggerUser, err := h.service.FindGriefLoggerUser(ctx, uuid)
	if err != nil {
		fmt.Println("Error retrieving grief logger user")
		return player, err
	}

	player = Player{uuid, luckpermsPlayer.PrimaryGroup, griefLoggerUser.ID, griefLoggerUser.Name}

	return player, nil
}

func (h handler) ListPlayer(w http.ResponseWriter, r *http.Request) {
	playerUuid := chi.URLParam(r, "uuid")

	player, err := h.getPlayer(r.Context(), playerUuid)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.Write(w, http.StatusOK, player)
}
