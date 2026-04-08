package responsemodels

type SingleDailyChartItem struct {
	Date  string
	Value int64
}

type SingleDailyChart struct {
	Title        string
	ChartData    []SingleDailyChartItem
	TimeDivision string
	Trend        float64
}
