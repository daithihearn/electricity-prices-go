package price

import "time"

type PriceClient interface {
	GetPrices(t time.Time) ([]Price, bool, error)
}
