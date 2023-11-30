package sync

import (
	"context"
	"electricity-prices/pkg/date"
	"electricity-prices/pkg/price"
	"log"
	"time"
)

type Service struct {
	PriceService  price.Service
	PrimaryClient price.PriceClient
	BackupClient  price.PriceClient
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

	// Keep processing until we reach tomorrow
	for {

		// Get the prices from the primary API
		prices, synced, err := s.PrimaryClient.GetPrices(currentDate)

		// If there is an error or the primary API is synced, try the backup API
		if err != nil && synced {
			prices, synced, err = s.BackupClient.GetPrices(currentDate)
		}

		// If there is an error exit
		if err != nil {
			panic(err)
		}

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
