package players

import (
	"context"
	"fmt"
	"log"
	"net/http"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
	"github.com/Will-Gould/psmp-api/internal/grieflogger"
	"github.com/Will-Gould/psmp-api/internal/json"
	"github.com/go-chi/chi"
)

type handler struct {
	service Service
	glh     *grieflogger.GriefLoggerHandler
}

type Player struct {
	Uuid          string
	PrimaryGroup  string
	GriefLoggerId int32
	Username      string
}

type Overview struct {
	BlocksBroken int64
	BlocksPlaced int64
}

type BlockData struct {
	BlocksBroken []repo.Block
	BlocksPlaced []repo.Block
}

func NewHandler(service Service, glh *grieflogger.GriefLoggerHandler) *handler {
	return &handler{
		service: service,
		glh:     glh,
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

func (h handler) ListPlayerBlockData(w http.ResponseWriter, r *http.Request) {
	playerUuid := chi.URLParam(r, "uuid")
	blockData := BlockData{}

	// get player
	player, err := h.getPlayer(r.Context(), playerUuid)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// get block data
	blocksBroken, err := h.service.FindBlocksBrokenByPlayer(r.Context(), player)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	blocksPlaced, err := h.service.FindBlocksPlacedByPlayer(r.Context(), player)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	blockData.BlocksBroken = blocksBroken
	blockData.BlocksPlaced = blocksPlaced

	json.Write(w, http.StatusOK, blockData)
}

func (h handler) GetPlayerOverview(w http.ResponseWriter, r *http.Request) {
	playerUuid := chi.URLParam(r, "uuid")
	overview := Overview{}

	// get player
	player, err := h.getPlayer(r.Context(), playerUuid)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// blocksBroken, err := h.service.CountBlocksBroken(r.Context(), player.GriefLoggerId, h.glh.BannedBrokenMaterials)
	blocksBroken, err := h.glh.Service.CountBlocksBrokenByUser(r.Context(), player.GriefLoggerId, h.glh.BannedBrokenMaterials)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	overview.BlocksBroken = blocksBroken

	json.Write(w, http.StatusOK, overview)
}
