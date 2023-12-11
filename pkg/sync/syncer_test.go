package sync

import (
	"context"
	"electricity-prices/pkg/date"
	"electricity-prices/pkg/price"
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
		savePricesCount          *price.CallCounter
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
			primaryGetPricesSynced:   &[]bool{false, true},
			primaryGetPricesErr:      &[]error{nil},
			secondaryGetPricesResp:   &[][]price.Price{},
			secondaryGetPricesSynced: &[]bool{true},
			secondaryGetPricesErr:    &[]error{},
			savePricesCount: &price.CallCounter{
				Count: 1,
			},
			mockSavePricesErr: &[]error{nil},
			expectError:       false,
			expectSynced:      true,
		},
		{
			name:    "Save failed",
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
			primaryGetPricesSynced:   &[]bool{false, true},
			primaryGetPricesErr:      &[]error{nil},
			secondaryGetPricesResp:   &[][]price.Price{},
			secondaryGetPricesSynced: &[]bool{true},
			secondaryGetPricesErr:    &[]error{},
			savePricesCount: &price.CallCounter{
				Count: 1,
			},
			mockSavePricesErr: &[]error{fmt.Errorf("error")},
			expectError:       true,
			expectSynced:      false,
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
			primaryGetPricesErr:    &[]error{fmt.Errorf("error"), fmt.Errorf("error")},
			secondaryGetPricesResp: &[][]price.Price{
				{
					{
						DateTime: time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local),
						Price:    1.0,
					},
				},
			},
			secondaryGetPricesSynced: &[]bool{false, true},
			secondaryGetPricesErr:    &[]error{nil},
			savePricesCount: &price.CallCounter{
				Count: 1,
			},
			mockSavePricesErr: &[]error{nil},
			expectError:       false,
			expectSynced:      true,
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
			secondaryGetPricesSynced: &[]bool{false, true},
			secondaryGetPricesErr:    &[]error{nil},
			savePricesCount: &price.CallCounter{
				Count: 1,
			},
			mockSavePricesErr: &[]error{nil},
			expectError:       false,
			expectSynced:      true,
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
			savePricesCount: &price.CallCounter{
				Count: 0,
			},
			mockSavePricesErr: &[]error{nil},
			expectError:       true,
			expectSynced:      false,
		},
		{
			name:    "Primary Client empty, Secondary Client empty",
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
			primaryGetPricesSynced:   &[]bool{false},
			primaryGetPricesErr:      &[]error{},
			secondaryGetPricesResp:   &[][]price.Price{},
			secondaryGetPricesSynced: &[]bool{false},
			secondaryGetPricesErr:    &[]error{},
			savePricesCount: &price.CallCounter{
				Count: 0,
			},
			mockSavePricesErr: &[]error{nil},
			expectError:       true,
			expectSynced:      false,
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
			savePricesCount: &price.CallCounter{
				Count: 0,
			},
			mockSavePricesErr: &[]error{nil},
			expectError:       true,
			expectSynced:      false,
		},
		{
			name:    "endDate is day after last synced date",
			endDate: time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local),
			getLatestPriceResp: &[]price.Price{
				{
					DateTime: time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local),
					Price:    1.0,
				},
			},
			getLatestPriceNoResult:   &[]bool{false},
			getLatestPriceErr:        &[]error{nil},
			primaryGetPricesResp:     &[][]price.Price{},
			primaryGetPricesSynced:   &[]bool{},
			primaryGetPricesErr:      &[]error{},
			secondaryGetPricesResp:   &[][]price.Price{},
			secondaryGetPricesSynced: &[]bool{},
			secondaryGetPricesErr:    &[]error{},
			savePricesCount: &price.CallCounter{
				Count: 0,
			},
			mockSavePricesErr: &[]error{nil},
			expectError:       false,
			expectSynced:      true,
		},
		{
			name:                   "Last price no result",
			endDate:                time.Date(2021, 6, 1, 0, 0, 0, 0, time.Local),
			getLatestPriceResp:     &[]price.Price{},
			getLatestPriceNoResult: &[]bool{true},
			getLatestPriceErr:      &[]error{nil},
			primaryGetPricesResp: &[][]price.Price{
				{
					{
						DateTime: date.StartOfDay(time.Date(2021, 5, 31, 0, 0, 0, 0, time.Local)),
						Price:    1.0,
					},
				},
			},
			primaryGetPricesSynced:   &[]bool{false, true},
			primaryGetPricesErr:      &[]error{nil},
			secondaryGetPricesResp:   &[][]price.Price{},
			secondaryGetPricesSynced: &[]bool{true},
			secondaryGetPricesErr:    &[]error{},
			savePricesCount: &price.CallCounter{
				Count: 1,
			},
			mockSavePricesErr: &[]error{nil},
			expectError:       false,
			expectSynced:      true,
		},
	}
	for _, test := range tests {
		// Create a mock Service
		mockPriceService := &price.MockPriceService{
			MockGetLatestPriceResult:   test.getLatestPriceResp,
			MockGetLatestPriceNoResult: test.getLatestPriceNoResult,
			MockGetLatestPriceError:    test.getLatestPriceErr,
			MockSavePricesCount:        test.savePricesCount,
			MockSavePricesError:        test.mockSavePricesErr,
		}

		// Mock primary and secondary clients
		mockPrimaryClient := &price.MockPriceClient{
			MockGetPricesResult: test.primaryGetPricesResp,
			MockGetPricesSynced: test.primaryGetPricesSynced,
			MockGetPricesError:  test.primaryGetPricesErr,
		}
		mockSecondaryClient := &price.MockPriceClient{
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
			if test.savePricesCount.Count != 0 {
				t.Errorf("Expected savePricesCount to be 0 but got %v", test.savePricesCount.Count)
			}
		})
	}
}
