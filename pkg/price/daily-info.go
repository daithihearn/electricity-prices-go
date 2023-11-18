package price

type DayRating string

const (
	Nil    DayRating = ""
	Good             = "GOOD"
	Normal           = "NORMAL"
	Bad              = "BAD"
)

type DailyPriceInfo struct {
	DayRating        DayRating `json:"dayRating"`
	DayAverage       float64   `json:"dayAverage"`
	ThirtyDayAverage float64   `json:"thirtyDayAverage"`
	Prices           []Price   `json:"prices"`
	CheapPeriods     [][]Price `json:"cheapestPeriods"`
	ExpensivePeriods [][]Price `json:"expensivePeriods"`
}
