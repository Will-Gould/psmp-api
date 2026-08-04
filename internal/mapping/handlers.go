package mapping

import (
	"context"
	"log"
	"log/slog"
	"slices"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
	"github.com/Will-Gould/psmp-api/internal/cache"
)

type mappingHandler struct {
	service    Service
	dataStore  *cache.DataStore
	gameLogger string
}

func NewHandler(service Service, ds *cache.DataStore, gameLogger string) *mappingHandler {
	return &mappingHandler{
		service:    service,
		dataStore:  ds,
		gameLogger: gameLogger,
	}
}

func (mh mappingHandler) LoadMappingData(ctx context.Context) {
	objects := mh.GetObjects(ctx)
	bannedPlacedObjects, bannedBrokenObjects := mh.GetBannedObjects(ctx, objects)
	actionIdentifiers := mh.GetActionIdentifiers(ctx)
	causeMapping := mh.GetDeathCauseMapping(ctx)
	mobMapping := mh.GetMobMapping(ctx)
	advancementMapping := mh.GetAdvancementMapping(ctx)

	mh.dataStore.MappingData.Objects = objects
	mh.dataStore.MappingData.BannedBrokenObjects = bannedBrokenObjects
	mh.dataStore.MappingData.BannedPlacedObjects = bannedPlacedObjects
	mh.dataStore.MappingData.Actions = actionIdentifiers
	mh.dataStore.MappingData.CauseMapping = causeMapping
	mh.dataStore.MappingData.MobMapping = mobMapping
	mh.dataStore.MappingData.AdvancementMapping = advancementMapping
}

func (mh mappingHandler) GetObjects(ctx context.Context) []cache.Object {
	// get materials mapping
	// slog.Log(ctx, slog.LevelInfo, "Punching trees...")

	var objects []cache.Object
	switch mh.gameLogger {
	// GriefLogger
	case "grieflogger":
		materials, err := mh.service.ListGriefLoggerMaterials(ctx)
		if err != nil {
			log.Default()
			log.Panic("Failed to initialise materials")
		}
		for _, m := range materials {
			// for the case of grieflogger "minecraft:" will be prepended for consitency
			objects = append(objects, cache.Object{ID: m.ID, Name: "minecraft:" + m.Name})
		}
	// Ledger
	default:
		ledgerObjects, err := mh.service.ListLedgerObjects(ctx)
		if err != nil {
			log.Default()
			log.Panic("Failed to initialise materials")
		}
		for _, o := range ledgerObjects {
			objects = append(objects, cache.Object{ID: o.ID, Name: o.Identifier})
		}
	}

	return objects
}

func (mh mappingHandler) GetBannedObjects(ctx context.Context, objects []cache.Object) ([]int32, []int32) {
	bannedPlacedObjects, bannedBrokenObjects := []int32{}, []int32{}
	// add banned material IDs
	// slog.Log(ctx, slog.LevelInfo, "Brewing potions...")
	for _, o := range objects {
		if slices.Contains(BANNED_PLACED_MATERIALS, o.Name) {
			bannedPlacedObjects = append(bannedPlacedObjects, o.ID)
		}
		if slices.Contains(BANNED_BROKEN_MATERIALS, o.Name) {
			bannedBrokenObjects = append(bannedBrokenObjects, o.ID)
		}
	}

	return bannedPlacedObjects, bannedBrokenObjects
}

func (mh mappingHandler) GetAdvancementMapping(ctx context.Context) []repo.PsmpstatsAdvancement {
	// get advancement mapping
	// slog.Log(ctx, slog.LevelInfo, "Shearing sheep...")
	advancementMap, err := mh.service.ListAdvancements(ctx)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Failed to get advancement mapping")
		log.Panic()
	}
	return advancementMap
}

func (mh mappingHandler) GetMobMapping(ctx context.Context) []repo.PsmpstatsMob {
	// get mob mapping
	// slog.Log(ctx, slog.LevelInfo, "Killing zombies...")
	mobMap, err := mh.service.ListMobs(ctx)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Failed to get mob mapping")
		log.Panic()
	}
	return mobMap
}

func (mh mappingHandler) GetDeathCauseMapping(ctx context.Context) []repo.PsmpstatsCause {
	// get death cause mapping
	// slog.Log(ctx, slog.LevelInfo, "Jumping over ravines...")
	causeMap, err := mh.service.ListCausesOfDeath(ctx)
	if err != nil {
		slog.Log(ctx, slog.LevelError, "Failed to get causes of death mapping")
		log.Panic()
	}
	return causeMap
}

func (mh mappingHandler) GetActionIdentifiers(ctx context.Context) map[string]int32 {

	actions := make(map[string]int32)
	switch mh.gameLogger {
	case "grieflogger":
		actions["block-break"] = 0
		actions["block-place"] = 1
		actions["player-join"] = 0
		actions["player-leave"] = 1
	default:
		actionIdentifiers, err := mh.service.ListActionIdentifiers(ctx)
		if err != nil {
			slog.Log(ctx, slog.LevelError, "Failed to load Ledger action identifiers")
			log.Panic()
		}
		for _, a := range actionIdentifiers {
			//actions = append(actions, cache.Action{ID: a.ID, Name: a.ActionIdentifier})
			actions[a.ActionIdentifier] = a.ID
		}
	}

	return actions
}
