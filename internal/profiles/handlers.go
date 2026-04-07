package profiles

import (
	"log/slog"
	"net/http"
	"sort"
	"time"

	"github.com/Will-Gould/psmp-api/internal/json"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
	"github.com/go-chi/chi"
)

const DATE_FORMAT = "2006-01-02"

type profileHandler struct {
	service Service
}

func NewHandler(service Service) *profileHandler {
	return &profileHandler{
		service: service,
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

	startDate := time.Unix(int64(mobKills[0].Time), 0)
	endDate := time.Now().AddDate(0, 0, 1)
	// initialise value for each date between now and first kill
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
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
