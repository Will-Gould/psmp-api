package profiles

import (
	"log"
	"log/slog"
	"net/http"
	"slices"
	"sort"
	"time"

	"github.com/Will-Gould/psmp-api/internal/json"
	"github.com/Will-Gould/psmp-api/internal/mapping"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
	"github.com/go-chi/chi"
)

const DATE_FORMAT = "2006-01-02"

type profileHandler struct {
	service     Service
	mappingData *mapping.MappingData
}

func NewHandler(service Service, md *mapping.MappingData) *profileHandler {
	return &profileHandler{
		service:     service,
		mappingData: md,
	}
}

func (ph *profileHandler) GetMobKillChartData(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	data := []responsemodels.SingleDailyChartItem{}
	mobKills, err := ph.service.ListMobKillsByUuid(r.Context(), uuid)
	if err != nil {
		slog.Log(r.Context(), slog.LevelError, err.Error())
		return
	}
	sort.Slice(mobKills, func(i, j int) bool {
		if mobKills[i].Time < mobKills[j].Time {
			return true
		}
		return false
	})

	firstKillTime := time.Unix(int64(mobKills[0].Time), 0).Local()
	earliestMidnight := firstKillTime.Truncate(24 * time.Hour)
	nextDay := time.Now().AddDate(0, 0, 1).Local()
	nextMidnight := nextDay.Truncate(24 * time.Hour)
	// initialise value for each date between now and first kill
	for d := earliestMidnight; !d.After(nextMidnight); d = d.AddDate(0, 0, 1) {
		dateString := d.Local().Format(DATE_FORMAT)
		data = append(data, responsemodels.SingleDailyChartItem{
			Date:  dateString,
			Value: 0,
		})
	}

	for _, k := range mobKills {
		killDate := time.Unix(int64(k.Time), 0).Local().Format(DATE_FORMAT)
		for i, d := range data {
			if d.Date == killDate {
				data[i].Value += 1
				continue
			}
		}
	}

	// calculate trend over last two intervals
	var trend float64 = 0
	if len(data) > 1 {
		cur := data[len(data)-1].Value
		base := data[len(data)-2].Value

		if base > 0 {
			trend = ((float64(cur) - float64(base)) / float64(base)) * 100
		}
	}

	chart := responsemodels.SingleDailyChart{
		Title:        "Mob Kills",
		ChartData:    data,
		TimeDivision: "days",
		Trend:        trend,
	}

	json.Write(w, http.StatusOK, chart)
}

func (ph *profileHandler) GetBlocksBrokenPieChartData(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	data := []responsemodels.BlockChartItem{}

	glUser, err := ph.service.FindGriefLoggerUser(r.Context(), uuid)
	if err != nil {
		json.Write(w, http.StatusNotFound, nil)
		return
	}

	blocksBroken, err := ph.service.GroupCountBlocksByUser(r.Context(), glUser.ID, mapping.BLOCK_BROKEN_ACTION, ph.mappingData.BannedBrokenMaterials)

	// transform into materials
	for _, b := range blocksBroken {
		name := ph.findMaterialName(b.Type)
		data = append(data, responsemodels.BlockChartItem{
			Block: name,
			Value: b.TotalPlaced,
		})
	}

	// consolidate into 'others' category
	if len(data) > 6 {
		sort.Slice(data, func(i, j int) bool {
			if data[i].Value > data[j].Value {
				return true
			}
			return false
		})

		data[5].Block = "others"

		for i := 6; i < len(data); i++ {
			data[5].Value += data[i].Value
		}

		data = slices.Delete(data, 6, len(data))
	}

	blockChart := responsemodels.BlockChart{
		Title:     "Blocks Broken",
		ChartData: data,
	}

	json.Write(w, http.StatusOK, blockChart)
}

func (ph *profileHandler) GetTotalBlocksChart(w http.ResponseWriter, r *http.Request) {
	uuid := chi.URLParam(r, "uuid")
	data := []responsemodels.DoubleDailyChartItem{}

	glUser, err := ph.service.FindGriefLoggerUser(r.Context(), uuid)
	if err != nil {
		json.Write(w, http.StatusNotFound, nil)
		return
	}

	blocksBroken, err := ph.service.ListBlocksByUser(r.Context(), glUser.ID, mapping.BLOCK_BROKEN_ACTION, ph.mappingData.BannedBrokenMaterials)
	if err != nil {
		json.Write(w, http.StatusInternalServerError, nil)
	}
	blocksPlaced, err := ph.service.ListBlocksByUser(r.Context(), glUser.ID, mapping.BLOCK_PLACED_ACTION, ph.mappingData.BannedPlacedMaterials)
	if err != nil {
		json.Write(w, http.StatusInternalServerError, nil)
	}

	// combine block action into one slice & sort
	blocks := append(blocksBroken, blocksPlaced...)

	sort.Slice(blocks, func(i, j int) bool {
		if blocks[i].Time < blocks[j].Time {
			return true
		}
		return false
	})
	firstBlockTime := time.Unix(blocks[0].Time/1000, 0)
	earliestMidnight := firstBlockTime.Truncate(24 * time.Hour)
	nextDay := time.Now().AddDate(0, 0, 1).Local()
	nextMidnight := nextDay.Truncate(24 * time.Hour)
	log.Default().Printf("First Time: %v", earliestMidnight)
	// initialise value for each date between now & first block action
	for d := earliestMidnight; !d.After(nextMidnight); d = d.AddDate(0, 0, 1) {
		dateString := d.Local().Format(DATE_FORMAT)
		data = append(data, responsemodels.DoubleDailyChartItem{
			Date:   dateString,
			Value1: 0,
			Value2: 0,
		})
	}

	for i, b := range blocks {
		blockDate := time.Unix(blocks[i].Time/1000, 0).Format(DATE_FORMAT)
		for j, d := range data {

			if d.Date == blockDate {

				if b.Action == mapping.BLOCK_PLACED_ACTION {
					data[j].Value1 += 1
				}
				if b.Action == mapping.BLOCK_BROKEN_ACTION {
					data[j].Value2 += 1
				}
			}
		}
	}

	// combine date values so each day has a running total
	for i := range data {
		if i == 0 {
			continue
		}
		data[i].Value1 += data[i-1].Value1
		data[i].Value2 += data[i-1].Value2
	}

	chart := responsemodels.DoubleDailyChart{
		Title:     "Total Blocks Placed vs Broken",
		ChartData: data,
		Trend:     0,
	}

	json.Write(w, http.StatusOK, chart)
}

func (ph *profileHandler) findMaterialName(id int32) string {
	for _, m := range ph.mappingData.Materials {
		if id == m.ID {
			return m.Name
		}
	}
	return ""
}
