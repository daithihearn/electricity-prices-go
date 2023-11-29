package sync

import (
	"context"
	"electricity-prices/pkg/date"
	"electricity-prices/pkg/price"
	"electricity-prices/pkg/ree"
	"log"
	"time"
)

type Service struct {
	PriceService price.Service
}

func (s *Service) SyncWithAPI(ctx context.Context) {
	log.Println("Starting to sync with API...")

	// Get last day that was synced from database.
	p, err := s.PriceService.GetLatestPrice(ctx)
	if err != nil {
		p = price.Price{DateTime: date.StartOfDay(time.Date(2021, 5, 31, 0, 0, 0, 0, time.Local))}
	}

	log.Println("Last day synced: ", p.DateTime.Format("January 2 2006"))
	currentDate := date.StartOfDay(p.DateTime).AddDate(0, 0, 1)

	// If last day is after tomorrow then exit
	today := time.Now()
	tomorrow := date.StartOfDay(today.AddDate(0, 0, 1))

	// Keep processing until we reach tomorrow
	for {
		if currentDate.After(tomorrow) {
			log.Println("Fully synced. Exiting...")
			break
		}

		// Get the prices from the API
		prices, synced, err := ree.GetPricesFromRee(currentDate)

		if synced {
			log.Println("Fully synced. Exiting...")
			break
		}

		if err != nil {
			panic(err)
		}

		log.Printf("Syncing prices for %s", currentDate.Format("January 2 2006"))

		// Save the prices in the database
		err = s.PriceService.SavePrices(ctx, prices)
		if err != nil {
			panic(err)
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

}
