package statistics

import (
	"context"
	"log"
	"net/http"
	"slices"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
	"github.com/Will-Gould/psmp-api/internal/json"
	"github.com/go-chi/chi"
)

const BLOCK_BROKEN_ACTION = 0
const BLOCK_PLACED_ACTION = 1

var BannedPlacedMaterialsList = []string{
	"short_grass",
	"tall_grass",
	"fire",
	"large_fern",
	"seagrass",
	"shulker_box",
	"allium",
	"azure_bluet",
	"red_tulip",
	"orange_tulip",
	"white_tulip",
	"lily_of_the_valley",
	"wildflowers",
	"pink_petals",
	"peony",
	"rose_bush",
	"lilac",
	"moss_carpet",
	"sculk_vein",
	"dandelion",
	"sunflower",
	"oxeye_daisy",
	"poppy",
	"cornflower",
	"pink_tulip",
	"blue_orchid",
	"leaf_litter",
}

var BannedBrokenMaterialsList = []string{
	"netherrack",
	"short_grass",
	"tall_grass",
	"fire",
	"oak_leaves",
	"fern",
	"large_fern",
	"birch_leaves",
	"spruce_leaves",
	"seagrass",
	"flowering_azalea_leaves",
	"dark_oak_leaves",
	"cherry_leaves",
	"soul_soil",
	"soul_sand",
	"jungle_leaves",
	"azalea_leaves",
	"mangrove_leaves",
	"pale_oak_leaves",
	"dead_bush",
	"shulker_box",
	"allium",
	"azure_bluet",
	"red_tulip",
	"orange_tulip",
	"white_tulip",
	"lily_of_the_valley",
	"firefly_bush",
	"wildflowers",
	"pink_petals",
	"wither_rose",
	"open_eyeblossom",
	"closed_eyeblossom",
	"cactus_flower",
	"bamboo_sapling",
	"crimson_roots",
	"warped_roots",
	"twisting_vines",
	"large_fern",
	"hanging_roots",
	"pitcher_plant",
	"peony",
	"rose_bush",
	"lilac",
	"moss_carpet",
	"sculk_vein",
	"seagrass",
	"sea_pickle",
	"tube_coral",
	"brain_coral",
	"bubble_coral",
	"fire_coral",
	"horn_coral",
	"dead_tube_coral",
	"dead_brain_coral",
	"dead_bubble_coral",
	"dead_tube_coral_fan",
	"dead_brain_coral_fan",
	"horn_coral_fan",
	"fire_coral_fan",
	"bubble_coral_fan",
	"brain_coral_fan",
	"tube_coral_fan",
	"dead_horn_coral",
	"dead_fire_coral",
	"dead_bubble_coral_fan",
	"dead_fire_coral_fan",
	"dead_horn_coral_fan",
	"pale_hanging_moss",
	"weeping_vines",
	"dandelion",
	"sunflower",
	"oxeye_daisy",
	"dead_bush",
	"poppy",
	"cornflower",
	"pink_tulip",
	"blue_orchid",
	"torchflower",
	"bamboo",
	"vine",
	"sugar_cane",
	"kelp",
	"leaf_litter",
	"snow",
}

type StatisticsHandler struct {
	Service               Service
	Materials             []repo.Material
	BannedPlacedMaterials []int32
	BannedBrokenMaterials []int32
}

type Overview struct {
	BlocksBroken int64
	BlocksPlaced int64
}

type BlockData struct {
	BlocksBroken []repo.Block
	BlocksPlaced []repo.Block
}

func NewHandler(service Service) *StatisticsHandler {
	var bannedPlacedMaterials []int32
	var bannedBrokenMaterials []int32

	// get materials mapping
	materials, err := service.ListMaterials(context.Background())
	if err != nil {
		log.Default()
		log.Panic("Failed to initialise materials")
	}

	// add banned material IDs
	for _, m := range materials {
		if slices.Contains(BannedPlacedMaterialsList, m.Name) {
			bannedPlacedMaterials = append(bannedPlacedMaterials, m.ID)
		}
		if slices.Contains(BannedBrokenMaterialsList, m.Name) {
			bannedBrokenMaterials = append(bannedBrokenMaterials, m.ID)
		}
	}

	return &StatisticsHandler{
		Service:               service,
		Materials:             materials,
		BannedPlacedMaterials: bannedPlacedMaterials,
		BannedBrokenMaterials: bannedBrokenMaterials,
	}
}

func (sh StatisticsHandler) ShowPlayerOverview(w http.ResponseWriter, r *http.Request) {
	playerUuid := chi.URLParam(r, "uuid")
	overview := Overview{}

	glUser, err := sh.Service.FindGriefLoggerUser(r.Context(), playerUuid)
	if err != nil {
		log.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	blocksBroken, err := sh.Service.CountBlocksByUser(r.Context(), glUser.ID, BLOCK_BROKEN_ACTION, sh.BannedBrokenMaterials)
	blocksPlaced, err := sh.Service.CountBlocksByUser(r.Context(), glUser.ID, BLOCK_PLACED_ACTION, sh.BannedPlacedMaterials)

	overview.BlocksBroken = blocksBroken
	overview.BlocksPlaced = blocksPlaced

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
