package sync

import (
	"context"
	"electricity-prices/pkg/price"
	"electricity-prices/pkg/price/testdata"
	"fmt"
	"testing"
	"time"
)

func TestSync(t *testing.T) {
	tests := []struct {
		name                     string
		endDate                  time.Time
		getLatestPriceResp       *[]price.Price
		getLatestPriceNoResult   *[]bool
		getLatestPriceErr        *[]error
		primaryGetPricesResp     *[][]price.Price
		primaryGetPricesSynced   *[]bool
		primaryGetPricesErr      *[]error
		secondaryGetPricesResp   *[][]price.Price
		secondaryGetPricesSynced *[]bool
		secondaryGetPricesErr    *[]error
		mockSavePricesErr        *[]error
		expectError              bool
		expectSynced             bool
	}{
		{
			name:    "Primary Client successful",
			endDate: time.Date(2023, 6, 1, 0, 0, 0, 0, time.Local),
			getLatestPriceResp: &[]price.Price{
				{
					DateTime: time.Date(2021, 5, 31, 0, 0, 0, 0, time.Local),
					Price:    1.0,
				},
			},
			getLatestPriceNoResult: &[]bool{false},
			getLatestPriceErr:      &[]error{nil},
			primaryGetPricesResp: &[][]price.Price{
				{
					{
						DateTime: time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local),
						Price:    1.0,
					},
				},
			},
			primaryGetPricesSynced:   &[]bool{true},
			primaryGetPricesErr:      &[]error{nil},
			secondaryGetPricesResp:   &[][]price.Price{},
			secondaryGetPricesSynced: &[]bool{},
			secondaryGetPricesErr:    &[]error{},
			mockSavePricesErr:        &[]error{nil},
			expectError:              false,
			expectSynced:             true,
		},
		{
			name:    "Primary Client error, Secondary Client successful",
			endDate: time.Date(2023, 6, 1, 0, 0, 0, 0, time.Local),
			getLatestPriceResp: &[]price.Price{
				{
					DateTime: time.Date(2021, 5, 31, 0, 0, 0, 0, time.Local),
				},
			},
			getLatestPriceNoResult: &[]bool{false},
			getLatestPriceErr:      &[]error{nil},
			primaryGetPricesResp:   &[][]price.Price{},
			primaryGetPricesSynced: &[]bool{},
			primaryGetPricesErr:    &[]error{fmt.Errorf("error")},
			secondaryGetPricesResp: &[][]price.Price{
				{
					{
						DateTime: time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local),
						Price:    1.0,
					},
				},
			},
			secondaryGetPricesSynced: &[]bool{true},
			secondaryGetPricesErr:    &[]error{nil},
			mockSavePricesErr:        &[]error{nil},
			expectError:              false,
			expectSynced:             true,
		},
		{
			name:    "Primary Client unsuccessful, Secondary Client successful",
			endDate: time.Date(2023, 6, 1, 0, 0, 0, 0, time.Local),
			getLatestPriceResp: &[]price.Price{
				{
					DateTime: time.Date(2021, 5, 31, 0, 0, 0, 0, time.Local),
					Price:    1.0,
				},
			},
			getLatestPriceNoResult: &[]bool{false},
			getLatestPriceErr:      &[]error{nil},
			primaryGetPricesResp:   &[][]price.Price{},
			primaryGetPricesSynced: &[]bool{},
			primaryGetPricesErr:    &[]error{},
			secondaryGetPricesResp: &[][]price.Price{
				{
					{
						DateTime: time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local),
						Price:    1.0,
					},
				},
			},
			secondaryGetPricesSynced: &[]bool{true},
			secondaryGetPricesErr:    &[]error{nil},
			mockSavePricesErr:        &[]error{nil},
			expectError:              false,
			expectSynced:             true,
		},
		{
			name:    "Primary Client unsuccessful, Secondary Client unsuccessful",
			endDate: time.Date(2023, 6, 1, 0, 0, 0, 0, time.Local),
			getLatestPriceResp: &[]price.Price{
				{
					DateTime: time.Date(2021, 5, 31, 0, 0, 0, 0, time.Local),
					Price:    1.0,
				},
			},
			getLatestPriceNoResult:   &[]bool{false},
			getLatestPriceErr:        &[]error{nil},
			primaryGetPricesResp:     &[][]price.Price{},
			primaryGetPricesSynced:   &[]bool{},
			primaryGetPricesErr:      &[]error{fmt.Errorf("error")},
			secondaryGetPricesResp:   &[][]price.Price{},
			secondaryGetPricesSynced: &[]bool{},
			secondaryGetPricesErr:    &[]error{fmt.Errorf("error")},
			mockSavePricesErr:        &[]error{nil},
			expectError:              true,
			expectSynced:             false,
		},
		{
			name:                     "Error getting latest price",
			endDate:                  time.Date(2023, 6, 1, 0, 0, 0, 0, time.Local),
			getLatestPriceResp:       &[]price.Price{},
			getLatestPriceNoResult:   &[]bool{false},
			getLatestPriceErr:        &[]error{fmt.Errorf("error")},
			primaryGetPricesResp:     &[][]price.Price{},
			primaryGetPricesSynced:   &[]bool{},
			primaryGetPricesErr:      &[]error{},
			secondaryGetPricesResp:   &[][]price.Price{},
			secondaryGetPricesSynced: &[]bool{},
			secondaryGetPricesErr:    &[]error{},
			mockSavePricesErr:        &[]error{nil},
			expectError:              true,
			expectSynced:             false,
		},
	}
	for _, test := range tests {
		// Create a mock PriceService
		mockPriceService := &testdata.MockPriceService{
			MockGetLatestPriceResult:   test.getLatestPriceResp,
			MockGetLatestPriceNoResult: test.getLatestPriceNoResult,
			MockGetLatestPriceError:    test.getLatestPriceErr,
			MockSavePricesError:        test.mockSavePricesErr,
		}

		// Mock primary and secondary clients
		mockPrimaryClient := &testdata.MockPriceClient{
			MockGetPricesResult: test.primaryGetPricesResp,
			MockGetPricesSynced: test.primaryGetPricesSynced,
			MockGetPricesError:  test.primaryGetPricesErr,
		}
		mockSecondaryClient := &testdata.MockPriceClient{
			MockGetPricesResult: test.secondaryGetPricesResp,
			MockGetPricesSynced: test.secondaryGetPricesSynced,
			MockGetPricesError:  test.secondaryGetPricesErr,
		}

		syncer := Syncer{
			PriceService:    mockPriceService,
			PrimaryClient:   mockPrimaryClient,
			SecondaryClient: mockSecondaryClient,
		}

		t.Run(test.name, func(t *testing.T) {
			synced, err := syncer.Sync(context.Background(), test.endDate)
			if test.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error but got %v", err)
				}
			}
			if synced != test.expectSynced {
				t.Errorf("Expected expectSynced to be %v but got %v", test.expectSynced, synced)
			}
		})
	}
}
