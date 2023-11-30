package esios

import (
	"electricity-prices/pkg/price"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const urlTemplate = "https://api.esios.ree.es/archives/70/download_json?date=%s"

// GetPrices returns the prices for the given date from the ERIOS API
func GetPrices(date time.Time) ([]price.Price, bool, error) {
	// Parse date to day string
	day := date.Format("2006-01-02")

	// Call to endpoint
	resp, err := http.Get(fmt.Sprintf(urlTemplate, day))
	if err != nil {
		log.Fatalf("Error occurred while sending request to the server: %s", err)
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
		log.Fatalf("Error occurred while reading response body: %s", err)
	}

	// Check if the status code indicates success
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		// Initialize the response object
		var res EsiosResponse

		// Parse the JSON response body into the response struct
		err := json.Unmarshal(body, &res)
		if err != nil {
			log.Fatalf("Error occurred while unmarshaling the response body: %s", err)
		}

		if len(res.PVPC) == 0 {
			log.Printf("No prices for %s", day)
			return nil, true, nil
		}

		prices := make([]price.Price, len(res.PVPC))

		for i, p := range res.PVPC {
			convertedP, err := convertStringToFloat(p.PCB)
			if err != nil {
				return nil, false, fmt.Errorf("error converting price: %v", err)
			}
			convetedDate, err := convertToTime(p.Day, p.Hour)
			if err != nil {
				return nil, false, fmt.Errorf("error converting date: %v", err)
			}
			prices[i] = price.Price{
				DateTime: convetedDate,
				Price:    convertedP / 1000,
			}
		}

		return prices, false, nil

	}
	return nil, false, fmt.Errorf("server responded with a non-successful status code: %d", resp.StatusCode)
}

func convertStringToFloat(s string) (float64, error) {
	// Replace comma with period
	s = strings.Replace(s, ",", ".", -1)

	// Convert string to float
	return strconv.ParseFloat(s, 64)
}

func convertToTime(dateStr string, hourRange string) (time.Time, error) {
	// Convert hour range to integer
	hour, err := convertHourRangeToIn(hourRange)
	if err != nil {
		return time.Time{}, err
	}

	// Layout of the input date string (this must match the format of dateStr)
	layout := "02/01/2006"

	// Parse the date string
	date, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, err
	}

	// Create a new time with the specified hour, minute, and second
	newTime := time.Date(date.Year(), date.Month(), date.Day(), hour, 0, 0, 0, date.Location())

	return newTime, nil
}

func convertHourRangeToIn(hourRange string) (int, error) {
	// Check if the string is at least 2 characters long
	if len(hourRange) < 2 {
		return 0, fmt.Errorf("string is too short")
	}

	// Extract the first two characters
	firstTwo := hourRange[:2]

	// Convert to integer
	return strconv.Atoi(firstTwo)
}