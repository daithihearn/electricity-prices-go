package esios

import (
	"electricity-prices/pkg/date"
	"testing"
	"time"
)

func TestGetPrices(t *testing.T) {

	tests := []struct {
		name               string
		testDate           time.Time
		expectedResultSize int
		expectSynced       bool
		expectingError     bool
	}{
		{
			name:               "Date that is available",
			testDate:           time.Date(2022, 10, 11, 0, 0, 0, 0, date.Location),
			expectedResultSize: 24,
			expectSynced:       false,
			expectingError:     false,
		},
		{
			name:               "Date in the future that is not available",
			testDate:           time.Now().AddDate(1, 1, 1),
			expectedResultSize: 0,
			expectSynced:       true,
			expectingError:     false,
		},
		{
			name:               "Date in the past that is not available",
			testDate:           time.Date(2000, 10, 11, 0, 0, 0, 0, date.Location),
			expectedResultSize: 0,
			expectSynced:       false,
			expectingError:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			prices, synced, err := GetPrices(test.testDate)

			if test.expectingError {
				if err == nil {
					t.Errorf("Expected error but got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got %s", err)
				}
			}

			if synced != test.expectSynced {
				t.Errorf("Expected synced to be %t but got %t", test.expectSynced, synced)
			}

			if len(prices) != test.expectedResultSize {
				t.Errorf("Expected %d prices but got %d", test.expectedResultSize, len(prices))
			}

		})
	}

}
