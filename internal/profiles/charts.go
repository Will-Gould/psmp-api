package profiles

import (
	"sort"
	"time"

	repo "github.com/Will-Gould/psmp-api/internal/adapters/mysql/sqlc"
	responsemodels "github.com/Will-Gould/psmp-api/internal/response_models"
	"github.com/Will-Gould/psmp-api/internal/translation"
)

func getAllTimeSingleDailySlice(firstJoin translation.Session) []responsemodels.SingleDailyChartItem {
	data := []responsemodels.SingleDailyChartItem{}

	firstTime := time.Unix(int64(firstJoin.Time/1000), 0).Local()
	earliestMidnight := firstTime.Truncate(24 * time.Hour)

	nextDay := time.Now().AddDate(0, 0, 1).Local()
	nextMidnight := nextDay.Truncate(24 * time.Hour)
	// initialise value for each date between now and first join
	for d := earliestMidnight; !d.After(nextMidnight); d = d.AddDate(0, 0, 1) {
		dateString := d.Local().Format(DATE_FORMAT)
		data = append(data, responsemodels.SingleDailyChartItem{
			Date:  dateString,
			Value: 0,
		})
	}

	return data
}

func GetMobKillChartData(mobKills []repo.PsmpstatsMobKill, firstJoin translation.Session) responsemodels.SingleDailyChart {

	sort.Slice(mobKills, func(i, j int) bool {
		if mobKills[i].Time < mobKills[j].Time {
			return true
		}
		return false
	})

	data := getAllTimeSingleDailySlice(firstJoin)

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
		Title:     "Mob Kills",
		ChartData: data,
		Trend:     trend,
	}

	return chart
}

func GetTotalBlocksChartData(blocks []translation.Block, actions map[string]int32) responsemodels.DoubleDailyChart {
	data := []responsemodels.DoubleDailyChartItem{}

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

				if b.Action == actions["block-break"] {
					data[j].Value1 += 1
				}
				if b.Action == actions["block-place"] {
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

	return chart
}

func GetDeathsChartData(deaths []repo.PsmpstatsDeath, firstJoin translation.Session) responsemodels.SingleDailyChart {
	sort.Slice(deaths, func(i, j int) bool {
		if deaths[i].Time < deaths[j].Time {
			return true
		}
		return false
	})

	data := getAllTimeSingleDailySlice(firstJoin)

	for _, k := range deaths {
		deathDate := time.Unix(int64(k.Time), 0).Local().Format(DATE_FORMAT)
		for i, d := range data {
			if d.Date == deathDate {
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
		Title:     "Deaths",
		ChartData: data,
		Trend:     trend,
	}

	return chart
}
