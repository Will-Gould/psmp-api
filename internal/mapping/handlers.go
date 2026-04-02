package mapping

import (
	"context"
	"log"
	"log/slog"
	"slices"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
)

type mappingHandler struct {
	service Service
}

func NewHandler(service Service) *mappingHandler {
	return &mappingHandler{
		service: service,
	}
}

func (mh mappingHandler) GetMaterials(ctx context.Context) []repo.Material {
	// get materials mapping
	slog.Log(ctx, slog.LevelInfo, "Punching trees...")
	materials, err := mh.service.ListMaterials(ctx)
	if err != nil {
		log.Default()
		log.Panic("Failed to initialise materials")
	}

	return materials
}

func (mh mappingHandler) GetBannedMaterials(ctx context.Context, materials []repo.Material) ([]int32, []int32) {
	bannedPlacedMaterials, bannedBrokenMaterials := []int32{}, []int32{}
	// add banned material IDs
	slog.Log(ctx, slog.LevelInfo, "Brewing potions...")
	for _, m := range materials {
		if slices.Contains(BANNED_PLACED_MATERIALS, m.Name) {
			bannedPlacedMaterials = append(bannedPlacedMaterials, m.ID)
		}
		if slices.Contains(BANNED_BROKEN_MATERIALS, m.Name) {
			bannedBrokenMaterials = append(bannedBrokenMaterials, m.ID)
		}
	}
	return bannedPlacedMaterials, bannedBrokenMaterials
}

func (mh mappingHandler) GetAdvancementMapping(ctx context.Context) []repo.PsmpstatsAdvancement {
	// get advancement mapping
	slog.Log(ctx, slog.LevelInfo, "Shearing sheep...")
	advancementMap, err := mh.service.ListAdvancements(ctx)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Failed to get advancement mapping")
		log.Panic()
	}
	return advancementMap
}

func (mh mappingHandler) GetMobMapping(ctx context.Context) []repo.PsmpstatsMob {
	// get mob mapping
	slog.Log(ctx, slog.LevelInfo, "Killing zombies...")
	mobMap, err := mh.service.ListMobs(ctx)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Failed to get mob mapping")
		log.Panic()
	}
	return mobMap
}

func (mh mappingHandler) GetDeathCauseMapping(ctx context.Context) []repo.PsmpstatsCause {
	// get death cause mapping
	slog.Log(ctx, slog.LevelInfo, "Jumping over ravines...")
	causeMap, err := mh.service.ListCausesOfDeath(ctx)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Failed to get causes of death mapping")
		log.Panic()
	}
	return causeMap
}
