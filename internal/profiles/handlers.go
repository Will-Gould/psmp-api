package profiles

import (
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

func (ph *profileHandler) findMaterialName(id int32) string {
	for _, m := range ph.mappingData.Materials {
		if id == m.ID {
			return m.Name
		}
	}
	return ""
}
