package ree

import (
	"electricity-prices/pkg/price"
	"electricity-prices/pkg/web"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"
)

const urlTemplate = "https://apidatos.ree.es/en/datos/mercados/precios-mercados-tiempo-real?time_trunc=hour&start_date=%sT00:00&end_date=%sT23:59"

type ReeClient struct {
	Http web.HTTPClient
}

// GetPrices returns the prices for the given date from the REE API
func (c *ReeClient) GetPrices(t time.Time) ([]price.Price, bool, error) {
	// Parse date to day string
	day := t.Format("2006-01-02")

	// Call to endpoint
	resp, err := c.Http.Get(fmt.Sprintf(urlTemplate, day, day))
	if err != nil {
		return nil, false, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Fatalf("Error occurred while closing response body: %s", err)
		}
	}(resp.Body)

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, false, err
	}

	// Check if the status code indicates success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {

		// Initialize the response object
		var res ReeResponse

		// Parse the JSON response body into the response struct
		err = json.Unmarshal(body, &res)
		if err != nil {
			return nil, false, err
		}

		var included ReeIncluded
		for _, inc := range res.Included {
			if inc.ID == "1001" {
				included = inc
				continue
			}
		}

		if len(included.Attributes.Values) == 0 {
			return nil, false, fmt.Errorf("failed to parse response for date %s", day)
		}

		prices := make([]price.Price, len(included.Attributes.Values))

		for i, p := range included.Attributes.Values {
			prices[i] = price.Price{
				DateTime: p.DateTime,
				Price:    p.Price / 1000,
			}
		}

		return prices, false, nil

	} else if resp.StatusCode == 404 {
		// If the date is in the future, return synced as true
		if t.After(time.Now()) {
			return nil, true, nil
		}
	} else if resp.StatusCode == 502 {
		// Initialize the error response object
		var res ReeErrorResponse

		// Parse the JSON response body into the response struct
		err = json.Unmarshal(body, &res)
		if err != nil {
			log.Fatalf("Error occurred while unmarshaling the error response body: %s", err)
		}

		if len(res.Errors) > 0 && res.Errors[0].Detail == "There are no data for the selected filters." {
			// If the date is in the future, return synced as true
			if t.After(time.Now()) {
				return nil, true, nil
			}
		}
	}

	return nil, false, fmt.Errorf("server responded with a non-successful status code: %d", resp.StatusCode)

}
