package price

import "time"

type Client interface {
	GetPrices(t time.Time) ([]Price, bool, error)
}
