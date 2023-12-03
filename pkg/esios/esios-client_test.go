package esios

import (
	"electricity-prices/pkg/date"
	"electricity-prices/pkg/testutils"
	"electricity-prices/pkg/web"
	"errors"
	"net/http"
	"testing"
	"time"
)

func TestGetPrices_Integration(t *testing.T) {

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

	client := Client{Http: &http.Client{Timeout: time.Second * 30}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			prices, synced, err := client.GetPrices(test.testDate)

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

func TestGetPrices(t *testing.T) {
	tests := []struct {
		name               string
		testDate           time.Time
		mockResponse       *http.Response
		mockError          error
		expectedResultSize int
		expectSynced       bool
		expectingError     bool
	}{
		{
			name:     "Valid response",
			testDate: time.Date(2023, 11, 29, 0, 0, 0, 0, date.Location),
			mockResponse: &http.Response{StatusCode: 200, Body: web.NewMockReadCloser(
				testutils.ReadJsonStringFromFile("testdata/valid-2023-11-29.json"))},
			mockError:          nil,
			expectedResultSize: 24,
			expectSynced:       false,
			expectingError:     false,
		},
		{
			name:     "No values for specified archive - in future",
			testDate: time.Now().AddDate(0, 0, 1),
			mockResponse: &http.Response{StatusCode: 200, Body: web.NewMockReadCloser(
				testutils.ReadJsonStringFromFile("testdata/no-values-for-archive.json"))},
			mockError:          nil,
			expectedResultSize: 0,
			expectSynced:       true,
			expectingError:     false,
		},
		{
			name:     "No values for specified archive - in past",
			testDate: time.Date(2000, 10, 11, 0, 0, 0, 0, date.Location),
			mockResponse: &http.Response{StatusCode: 200, Body: web.NewMockReadCloser(
				testutils.ReadJsonStringFromFile("testdata/no-values-for-archive.json"))},
			mockError:          nil,
			expectedResultSize: 0,
			expectSynced:       false,
			expectingError:     true,
		},
		{
			name:               "Invalid data returned",
			testDate:           time.Date(2022, 10, 11, 0, 0, 0, 0, date.Location),
			mockResponse:       &http.Response{StatusCode: 200, Body: web.NewMockReadCloser(`{"data": "invalid"}`)},
			mockError:          nil,
			expectedResultSize: 0,
			expectSynced:       false,
			expectingError:     true,
		},
		{
			name:           "500 error",
			testDate:       time.Date(2022, 10, 11, 0, 0, 0, 0, date.Location),
			mockResponse:   &http.Response{StatusCode: 500, Body: web.NewMockReadCloser("")},
			mockError:      nil,
			expectSynced:   false,
			expectingError: true,
		},
		{
			name:           "Error calling to API",
			testDate:       time.Date(2022, 10, 11, 0, 0, 0, 0, date.Location),
			mockResponse:   nil,
			mockError:      errors.New("mock error"),
			expectSynced:   false,
			expectingError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			client := Client{Http: &web.MockHTTPClient{
				MockResp: test.mockResponse,
				MockErr:  test.mockError,
			}}

			prices, synced, err := client.GetPrices(test.testDate)

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
