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

type BlockChartItem struct {
	Block string
	Value int64
}

type BlockChart struct {
	Title     string
	ChartData []BlockChartItem
}
