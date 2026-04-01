package statistics

import (
	"context"
	"errors"
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
	Leaderboard           map[string]LeaderboardRanks
	Materials             []repo.Material
	BannedPlacedMaterials []int32
	BannedBrokenMaterials []int32
	CauseMapping          []repo.PsmpstatsCause
	MobMapping            []repo.PsmpstatsMob
	AdvancementMapping    []repo.PsmpstatsAdvancement
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

	// get advancement mapping
	slog.Log(context.Background(), slog.LevelInfo, "Shearing sheep...")
	advancementMap, err := service.ListAdvancements(context.Background())
	if err != nil {
		slog.Log(context.Background(), slog.LevelError, "Failed to get advancement mapping")
		log.Panic()
	}

	return &StatisticsHandler{
		Service:               service,
		Materials:             materials,
		BannedPlacedMaterials: bannedPlacedMaterials,
		BannedBrokenMaterials: bannedBrokenMaterials,
		CauseMapping:          causeMap,
		MobMapping:            mobMap,
		AdvancementMapping:    advancementMap,
	}
}

func (sh StatisticsHandler) InitialiseLeaderboard() {
	l, err := sh.UpdateLeaderboard(context.Background())
	if err != nil {
		panic("Failed to initialise leaderboard!")
	}
	sh.Leaderboard = l
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

	// get overview
	overview, err = sh.GetOverview(r.Context(), player)
	if err != nil {
		slog.Log(r.Context(), slog.LevelError, err.Error())
	}

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

func (sh StatisticsHandler) ListAdvancements(w http.ResponseWriter, r *http.Request) {
	playerUuid := chi.URLParam(r, "uuid")

	advancements, err := sh.Service.ListAdvancementsByPlayer(r.Context(), playerUuid)
	if err != nil {
		slog.Log(r.Context(), slog.LevelError, "Failed to retrieve advancements")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.Write(w, http.StatusOK, advancements)
}

func (sh StatisticsHandler) ListPlayers(w http.ResponseWriter, r *http.Request) {
	players, err := sh.getPlayers(r.Context())
	if err != nil {
		slog.Log(r.Context(), slog.LevelError, "Failed to retrieve player overviews")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, players)
}

func (sh StatisticsHandler) ListFeaturedPlayers(w http.ResponseWriter, r *http.Request) {
	featuredPlayers := FeaturedPlayers{}
	players, err := sh.getPlayers(r.Context())
	if err != nil {
		slog.Log(r.Context(), slog.LevelError, "Failed to retrieve player overviews")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	featuredPlayers.ChampionPlayer = GetChampionPlayer(sh, r.Context(), players)
	featuredPlayers.BiggestBuilder = GetBiggestBuilder(sh, r.Context(), players)
	featuredPlayers.MostDangerousPlayer = GetMostDangerousPlayer(sh, r.Context(), players)

	json.Write(w, http.StatusOK, featuredPlayers)
}

func (sh StatisticsHandler) ListMobKillChartData(w http.ResponseWriter, r *http.Request) {
	playerUuid := chi.URLParam(r, "uuid")
	mobKillData := GetMobKillChartData(r.Context(), sh, playerUuid)

	json.Write(w, http.StatusOK, mobKillData)
}

func (sh StatisticsHandler) GetOverview(ctx context.Context, player Player) (Overview, error) {
	overview := Overview{}
	overview.Player = player

	// get combat overview
	combatOverview, err := sh.GetCombatOverview(ctx, player.Uuid)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
		return overview, errors.New("Failed to get combat overview for player: " + player.Uuid)
	}

	// get crafting overview
	craftingOverview, err := sh.GetCraftingOverview(ctx, player.Uuid, player.GlId)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
		return overview, errors.New("Failed to get crafting overview for player: " + player.Uuid)
	}

	// get story overview
	storyOverview, err := GetStoryOverview(ctx, sh, player.Uuid)
	if err != nil {
		slog.Log(ctx, slog.LevelError, err.Error())
		return overview, errors.New("Failed to get story overview for player: " + player.Uuid)
	}

	overview.CombatOverview = combatOverview
	overview.CraftingOverview = craftingOverview
	overview.StoryOverview = storyOverview

	// calculate total score
	overview.Score = overview.CombatOverview.CombatScore + overview.CraftingOverview.CraftingScore + overview.StoryOverview.StoryScore

	return overview, nil
}

func (sh StatisticsHandler) getPlayers(ctx context.Context) ([]Overview, error) {
	var players []Overview

	luckpermsPlayers, err := sh.Service.ListLuckpermsPlayers(ctx)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Failed to retrieve Luckperms players")
		return players, errors.New("Failed to retrieve Luckperms players")
	}

	// build overviews
	for _, p := range luckpermsPlayers {
		glUser, err := sh.Service.FindGriefLoggerUser(ctx, p.Uuid)
		if err != nil {
			continue
		}
		player := Player{
			Uuid:         p.Uuid,
			Name:         glUser.Name,
			GlId:         glUser.ID,
			PrimaryGroup: p.PrimaryGroup,
		}
		overview, err := sh.GetOverview(ctx, player)
		if err != nil {
			continue
		}
		players = append(players, overview)
	}

	return players, nil
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
