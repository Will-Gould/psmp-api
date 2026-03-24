package statistics

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"slices"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
	"github.com/Will-Gould/psmp-api/internal/json"
	"github.com/go-chi/chi"
)

type StatisticsHandler struct {
	Service               Service
	Materials             []repo.Material
	BannedPlacedMaterials []int32
	BannedBrokenMaterials []int32
	CauseMapping          []repo.PsmpstatsCause
	MobMapping            []repo.PsmpstatsMob
}

func NewHandler(service Service) *StatisticsHandler {
	var bannedPlacedMaterials []int32
	var bannedBrokenMaterials []int32

	// get materials mapping
	slog.Log(context.Background(), slog.LevelInfo, "Punching trees...")
	materials, err := service.ListMaterials(context.Background())
	if err != nil {
		log.Default()
		log.Panic("Failed to initialise materials")
	}

	// add banned material IDs
	slog.Log(context.Background(), slog.LevelInfo, "Brewing potions...")
	for _, m := range materials {
		if slices.Contains(BANNED_PLACED_MATERIALS, m.Name) {
			bannedPlacedMaterials = append(bannedPlacedMaterials, m.ID)
		}
		if slices.Contains(BANNED_BROKEN_MATERIALS, m.Name) {
			bannedBrokenMaterials = append(bannedBrokenMaterials, m.ID)
		}
	}

	// get mob mapping
	slog.Log(context.Background(), slog.LevelInfo, "Killing zombies...")
	mobMap, err := service.ListMobs(context.Background())
	if err != nil {
		slog.Log(context.Background(), slog.LevelError, "Failed to get mob mapping")
		log.Panic()
	}

	// get death cause mapping
	slog.Log(context.Background(), slog.LevelInfo, "Jumping over ravines...")
	causeMap, err := service.ListCausesOfDeath(context.Background())
	if err != nil {
		slog.Log(context.Background(), slog.LevelError, "Failed to get causes of death mapping")
		log.Panic()
	}

	return &StatisticsHandler{
		Service:               service,
		Materials:             materials,
		BannedPlacedMaterials: bannedPlacedMaterials,
		BannedBrokenMaterials: bannedBrokenMaterials,
		CauseMapping:          causeMap,
		MobMapping:            mobMap,
	}
}

func (sh StatisticsHandler) ShowPlayerOverview(w http.ResponseWriter, r *http.Request) {
	playerUuid := chi.URLParam(r, "uuid")
	overview := Overview{}

	glUser, err := sh.Service.FindGriefLoggerUser(r.Context(), playerUuid)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	lpPlayer, err := sh.Service.FindLuckpermsPlayer(r.Context(), playerUuid)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	player := Player{
		Uuid:         playerUuid,
		Name:         glUser.Name,
		GlId:         glUser.ID,
		PrimaryGroup: lpPlayer.PrimaryGroup,
	}

	overview.Player = player

	// get combat overview
	overview.CombatOverview, err = GetCombatOverview(r.Context(), sh, playerUuid)

	// get crafting overview
	overview.CraftingOverview, err = GetCraftingOverview(r.Context(), sh, playerUuid, glUser)

	// calculate total score
	overview.Score = overview.CombatOverview.CombatScore + overview.CraftingOverview.CraftingScore

	json.Write(w, http.StatusOK, overview)
}

func (sh StatisticsHandler) ListBlockData(w http.ResponseWriter, r *http.Request) {
	playerUuid := chi.URLParam(r, "uuid")
	blockData := BlockData{}

	glUser, err := sh.Service.FindGriefLoggerUser(r.Context(), playerUuid)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// get block data
	blocksBroken, err := sh.Service.ListBlocksBrokenByUser(r.Context(), glUser.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	blocksPlaced, err := sh.Service.ListBlocksPlacedByUser(r.Context(), glUser.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	blockData.BlocksBroken = blocksBroken
	blockData.BlocksPlaced = blocksPlaced
	json.Write(w, http.StatusOK, blockData)
}

func (sh StatisticsHandler) ListSessionData(w http.ResponseWriter, r *http.Request) {
	playerUuid := chi.URLParam(r, "uuid")

	glUser, err := sh.Service.FindGriefLoggerUser(r.Context(), playerUuid)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	sessions, err := sh.Service.ListSessionDataByUser(r.Context(), glUser.ID)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.Write(w, http.StatusOK, sessions)
}

func (sh StatisticsHandler) ListDeaths(w http.ResponseWriter, r *http.Request) {
	playerUuid := chi.URLParam(r, "uuid")

	deaths, err := sh.Service.ListDeathsByPlayer(r.Context(), playerUuid)
	if err != nil {
		slog.Log(r.Context(), slog.LevelError, "Failed to retrieve deaths")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.Write(w, http.StatusOK, deaths)
}

func (sh StatisticsHandler) GetMappings(w http.ResponseWriter, r *http.Request) {

	maps := struct {
		MobMap   []repo.PsmpstatsMob
		CauseMap []repo.PsmpstatsCause
	}{
		MobMap:   sh.MobMapping,
		CauseMap: sh.CauseMapping,
	}
	json.Write(w, http.StatusOK, maps)
}

func (sh StatisticsHandler) getMobId(mob string) uint64 {
	for _, v := range sh.MobMapping {
		if v.Name == mob {
			return v.ID
		}
	}

	return 0
}
