package service

import (
	"electricity-prices/pkg/client"
	"electricity-prices/pkg/db"
	"electricity-prices/pkg/model"
	"electricity-prices/pkg/utils"
	"log"
	"time"
)

func SyncWithAPI() {
	log.Println("Starting to sync with API...")

	// Get last day that was synced from database.
	p, err := db.GetLatestPrice()
	if err != nil {
		p = model.Price{DateTime: utils.StartOfDay(time.Date(2014, 3, 31, 0, 0, 0, 0, time.Local))}
	}
	currentDate := utils.StartOfDay(p.DateTime).AddDate(0, 0, 1)
	log.Println("Last day synced: ", currentDate)

	// If last day is after tomorrow then exit
	today := time.Now()
	tomorrow := utils.StartOfDay(today.AddDate(0, 0, 1))

	// Keep processing until we reach tomorrow
	for {
		if currentDate.After(tomorrow) {
			log.Println("Fully synced. Exiting...")
			break
		}

		// Get the prices from the API
		prices, synced, err := client.GetPricesFromRee(currentDate)

		if synced {
			log.Println("Fully synced. Exiting...")
			break
		}

		if err != nil {
			panic(err)
		}

		log.Printf("Syncing prices for %s", currentDate.Format("January 2 2006"))

		// Save the prices in the database
		err = db.SavePrices(prices)
		if err != nil {
			panic(err)
		}
		currentDate = currentDate.AddDate(0, 0, 1)
	}
}
