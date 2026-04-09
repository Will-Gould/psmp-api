package responsemodels

type SingleDailyChartItem struct {
	Date  string
	Value int64
}

type DoubleDailyChartItem struct {
	Date   string
	Value1 int64
	Value2 int64
}

type SingleDailyChart struct {
	Title     string
	ChartData []SingleDailyChartItem
	Trend     float64
}

type DoubleDailyChart struct {
	Title     string
	ChartData []DoubleDailyChartItem
	Trend     float64
}

type BlockChartItem struct {
	Block string
	Value int64
}

type BlockChart struct {
	Title     string
	ChartData []BlockChartItem
}
