package sync

import (
	"context"
	"electricity-prices/pkg/date"
	"electricity-prices/pkg/price"
	"log"
	"time"
)

type Syncer struct {
	PriceService    price.Service
	PrimaryClient   price.Client
	SecondaryClient price.Client
}

// Sync syncs the prices from the API to the database
// It returns a boolean indicating whether the sync was successful or not
// and an error if there was one
// It takes a context and a time indicating the end date of the sync
func (s *Syncer) Sync(ctx context.Context, end time.Time) (bool, error) {
	log.Println("Starting to sync with API...")

	// Get last day that was synced from database.
	p, notFound, err := s.PriceService.GetLatestPrice(ctx)
	if notFound {
		p = price.Price{DateTime: date.StartOfDay(time.Date(2021, 5, 31, 0, 0, 0, 0, time.Local))}
	}
	if err != nil {
		return false, err
	}

	log.Println("Last day synced: ", p.DateTime.Format("January 2 2006"))
	currentDate := date.StartOfDay(p.DateTime).AddDate(0, 0, 1)

	// Keep processing until we reach tomorrow
	for {
		// If we reach the end date, exit
		if currentDate.After(date.StartOfDay(end).Add(time.Hour)) {
			break
		}

		// Get the prices from the primary API
		prices, synced, err := s.PrimaryClient.GetPrices(currentDate)

		// If there is an error or the primary API is synced, try the backup API
		if err != nil || synced || len(prices) == 0 {
			err = nil
			prices, synced, err = s.SecondaryClient.GetPrices(currentDate)
		}

		// If there is an error exit
		if err != nil {
			return false, err
		}

		if synced {
			break
		}

		if len(prices) == 0 {
			log.Printf("No prices for %s. Exiting...", currentDate.Format("January 2 2006"))
			return false, nil
		}

		log.Printf("Syncing prices for %s", currentDate.Format("January 2 2006"))

		// Save the prices in the database
		err = s.PriceService.SavePrices(ctx, prices)
		if err != nil {
			return false, err
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	log.Println("Fully Synced. Exiting...")
	return true, nil
}
