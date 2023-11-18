package price

import (
	"time"
)

type Price struct {
	ID       string    `bson:"_id,omitempty" json:"-"`
	DateTime time.Time `bson:"dateTime" json:"dateTime"`
	Price    float64   `bson:"price" json:"price"`
}
