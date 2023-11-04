package client

import (
	"electricity-prices/pkg/model"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const urlTemplate = "https://apidatos.ree.es/en/datos/mercados/precios-mercados-tiempo-real?time_trunc=hour&start_date=%sT00:00&end_date=%sT23:59"

// GetPricesFromRee returns the prices for the given date from the REE API
func GetPricesFromRee(date time.Time) ([]model.Price, bool, error) {
	// Parse date to day string
	day := date.Format("2006-01-02")

	// Call to endpoint
	resp, err := http.Get(fmt.Sprintf(urlTemplate, day, day))
	if err != nil {
		log.Fatalf("Error occurred while sending request to the server: %s", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("Error occurred while reading response body: %s", err)
	}

	// Check if the status code indicates success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Initialize an instance of PVPC
		var res model.ReeResponse

		// Parse the JSON response body into the PVPC struct
		err := json.Unmarshal(body, &res)
		if err != nil {
			log.Fatalf("Error occurred while unmarshaling the response body: %s", err)
		}
		if len(res.Errors) > 0 {
			if res.Errors[0].Detail == "There are no data for the selected filters." {
				return nil, true, nil
			}
			return nil, false, fmt.Errorf("error returned from API: %v", res.Errors[0])
		}

		var included model.ReeIncluded
		for _, inc := range res.Included {
			if inc.ID == "600" {
				included = inc
				continue
			}
		}
		if len(included.Attributes.Values) == 0 {
			return nil, false, fmt.Errorf("no prices returned from API")
		}

		prices := make([]model.Price, len(included.Attributes.Values))

		for i, p := range included.Attributes.Values {
			prices[i] = model.Price{
				DateTime: p.DateTime,
				Price:    p.Price / 1000,
			}
		}

		return prices, false, nil
	}

	return nil, false, fmt.Errorf("server responded with a non-successful status code: %d", resp.StatusCode)

}
