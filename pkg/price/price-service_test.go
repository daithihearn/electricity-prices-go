package price

import (
	"context"
	"electricity-prices/pkg/date"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPrice(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	priceExample := Price{
		DateTime: now,
		Price:    1.0,
	}

	tests := []struct {
		name           string
		mockResult     *[]Price
		mockError      *[]error
		expectedResult Price
		expectingError bool
	}{
		{
			name:           "success",
			mockResult:     &[]Price{priceExample},
			mockError:      &[]error{},
			expectedResult: priceExample,
			expectingError: false,
		},
		{
			name:           "failure",
			mockResult:     &[]Price{},
			mockError:      &[]error{errors.New("not found")},
			expectedResult: Price{},
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{MockFindOneResult: tt.mockResult, MockFindOneErr: tt.mockError}
			service := &Receiver{Collection: mockCollection}

			result, err := service.GetPrice(ctx, now)

			if tt.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestGetPrices(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	priceExample := Price{
		DateTime: now,
		Price:    1.0,
	}

	tests := []struct {
		name           string
		mockResult     *[][]Price
		mockError      *[]error
		expectedResult []Price
		expectingError bool
	}{
		{
			name:           "success",
			mockResult:     &[][]Price{{priceExample}},
			mockError:      &[]error{},
			expectedResult: []Price{priceExample},
			expectingError: false,
		},
		{
			name:           "failure",
			mockResult:     &[][]Price{},
			mockError:      &[]error{errors.New("not found")},
			expectedResult: []Price{},
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{MockFindResult: tt.mockResult, MockFindErr: tt.mockError}
			service := &Receiver{Collection: mockCollection}

			result, err := service.GetPrices(ctx, now, now)

			if tt.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestGetDailyPrices(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	priceExample := Price{
		DateTime: now,
		Price:    1.0,
	}

	tests := []struct {
		name           string
		mockResult     *[][]Price
		mockError      *[]error
		expectedResult []Price
		expectingError bool
	}{
		{
			name:           "success",
			mockResult:     &[][]Price{{priceExample}},
			mockError:      &[]error{},
			expectedResult: []Price{priceExample},
			expectingError: false,
		},
		{
			name:           "failure",
			mockResult:     &[][]Price{{}},
			mockError:      &[]error{errors.New("not found")},
			expectedResult: []Price{},
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{MockFindResult: tt.mockResult, MockFindErr: tt.mockError}
			service := &Receiver{Collection: mockCollection}

			result, err := service.GetDailyPrices(ctx, now)

			if tt.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestGetDailyAverages(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	priceExample := Price{
		DateTime: now,
		Price:    1.0,
	}

	dailyAverageExample := DailyAverage{
		Date:    date.ParseToLocalDay(now),
		Average: 1.0,
	}

	tests := []struct {
		name           string
		mockResult     *[][]Price
		mockError      *[]error
		expectedResult []DailyAverage
		expectingError bool
	}{
		{
			name:           "success",
			mockResult:     &[][]Price{{priceExample}},
			mockError:      &[]error{},
			expectedResult: []DailyAverage{dailyAverageExample},
			expectingError: false,
		},
		{
			name:           "failure",
			mockResult:     &[][]Price{},
			mockError:      &[]error{errors.New("not found")},
			expectedResult: []DailyAverage{},
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{MockFindResult: tt.mockResult, MockFindErr: tt.mockError}
			service := &Receiver{Collection: mockCollection}

			result, err := service.GetDailyAverages(ctx, now, 3)

			if tt.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestGetDayRating(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	priceExample := Price{
		DateTime: now,
		Price:    1.0,
	}

	tests := []struct {
		name           string
		mockPrices     *[][]Price
		mockPricesErr  *[]error
		mockAvg        *[]float64
		mockAvgErr     *[]error
		expectedResult DayRating
		expectingError bool
	}{
		{
			name:           "success",
			mockPrices:     &[][]Price{{priceExample}},
			mockPricesErr:  &[]error{},
			mockAvg:        &[]float64{1.0},
			mockAvgErr:     &[]error{},
			expectedResult: "NORMAL",
			expectingError: false,
		},
		{
			name:           "failure",
			mockPrices:     &[][]Price{{}},
			mockPricesErr:  &[]error{errors.New("not found")},
			mockAvg:        &[]float64{1.0},
			mockAvgErr:     &[]error{},
			expectedResult: "NORMAL",
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{MockFindResult: tt.mockPrices, MockFindErr: tt.mockPricesErr, MockThirtyDayAvg: tt.mockAvg, MockThirtyDayErr: tt.mockAvgErr}
			service := &Receiver{Collection: mockCollection}

			result, err := service.GetDayRating(ctx, now)

			if tt.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestGetDayAverage(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	priceExample := Price{
		DateTime: now,
		Price:    1.0,
	}

	tests := []struct {
		name           string
		mockPrices     *[][]Price
		mockPricesErr  *[]error
		expectedResult float64
		expectingError bool
	}{
		{
			name:           "success",
			mockPrices:     &[][]Price{{priceExample}},
			mockPricesErr:  &[]error{},
			expectedResult: 1.0,
			expectingError: false,
		},
		{
			name:           "failure",
			mockPrices:     &[][]Price{{}},
			mockPricesErr:  &[]error{errors.New("not found")},
			expectedResult: 0.0,
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{MockFindResult: tt.mockPrices, MockFindErr: tt.mockPricesErr}
			service := &Receiver{Collection: mockCollection}

			result, err := service.GetDayAverage(ctx, now)

			if tt.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestGetCheapPeriods(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	priceExample := Price{
		DateTime: now,
		Price:    1.0,
	}

	tests := []struct {
		name           string
		mockPrices     *[][]Price
		mockPricesErr  *[]error
		mockAvg        *[]float64
		mockAvgErr     *[]error
		expectedResult [][]Price
		expectingError bool
	}{
		{
			name:           "success",
			mockPrices:     &[][]Price{{priceExample}},
			mockPricesErr:  &[]error{},
			mockAvg:        &[]float64{1.0},
			mockAvgErr:     &[]error{},
			expectedResult: [][]Price{{priceExample}},
			expectingError: false,
		},
		{
			name:           "failure",
			mockPrices:     &[][]Price{{}},
			mockPricesErr:  &[]error{errors.New("not found")},
			mockAvg:        &[]float64{1.0},
			mockAvgErr:     &[]error{},
			expectedResult: [][]Price{{}},
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{MockFindResult: tt.mockPrices, MockFindErr: tt.mockPricesErr, MockThirtyDayAvg: tt.mockAvg, MockThirtyDayErr: tt.mockAvgErr}
			service := &Receiver{Collection: mockCollection}

			result, err := service.GetCheapPeriods(ctx, now)

			if tt.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestGetExpensivePeriods(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	priceExample := Price{
		DateTime: now,
		Price:    1.0,
	}

	tests := []struct {
		name           string
		mockPrices     *[][]Price
		mockPricesErr  *[]error
		mockAvg        *[]float64
		mockAvgErr     *[]error
		expectedResult [][]Price
		expectingError bool
	}{
		{
			name:           "success",
			mockPrices:     &[][]Price{{priceExample}},
			mockPricesErr:  &[]error{},
			mockAvg:        &[]float64{1.0},
			mockAvgErr:     &[]error{},
			expectedResult: [][]Price{{priceExample}},
			expectingError: false,
		},
		{
			name:           "failure",
			mockPrices:     &[][]Price{{}},
			mockPricesErr:  &[]error{errors.New("not found")},
			mockAvg:        &[]float64{1.0},
			mockAvgErr:     &[]error{},
			expectedResult: [][]Price{{}},
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{MockFindResult: tt.mockPrices, MockFindErr: tt.mockPricesErr, MockThirtyDayAvg: tt.mockAvg, MockThirtyDayErr: tt.mockAvgErr}
			service := &Receiver{Collection: mockCollection}

			result, err := service.GetExpensivePeriods(ctx, now)

			if tt.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestGetThirtyDayAverage(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	tests := []struct {
		name           string
		mockAvg        *[]float64
		mockAvgErr     *[]error
		expectedResult float64
		expectingError bool
	}{
		{
			name:           "success",
			mockAvg:        &[]float64{1.0},
			mockAvgErr:     &[]error{},
			expectedResult: 1.0,
			expectingError: false,
		},
		{
			name:           "failure",
			mockAvg:        &[]float64{},
			mockAvgErr:     &[]error{errors.New("not found")},
			expectedResult: 0.0,
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{MockThirtyDayAvg: tt.mockAvg, MockThirtyDayErr: tt.mockAvgErr}
			service := &Receiver{Collection: mockCollection}

			result, err := service.GetThirtyDayAverage(ctx, now)

			if tt.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}
		})
	}
}

func TestGetLatestPrice(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	priceExample := Price{
		DateTime: now,
		Price:    1.0,
	}

	tests := []struct {
		name           string
		mockPrice      *[]Price
		mockPriceFound *[]bool
		mockPriceErr   *[]error
		expectedResult Price
		expectingError bool
		expectedFound  bool
	}{
		{
			name:           "success",
			mockPrice:      &[]Price{priceExample},
			mockPriceFound: &[]bool{true},
			mockPriceErr:   &[]error{},
			expectedResult: priceExample,
			expectingError: false,
			expectedFound:  true,
		},
		{
			name:           "failure",
			mockPrice:      &[]Price{},
			mockPriceFound: &[]bool{false},
			mockPriceErr:   &[]error{errors.New("not found")},
			expectedResult: Price{},
			expectingError: true,
			expectedFound:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{MockLatestPrice: tt.mockPrice, MockLatestPriceOk: tt.mockPriceFound, MockLatestPriceErr: tt.mockPriceErr}
			service := &Receiver{Collection: mockCollection}

			result, found, err := service.GetLatestPrice(ctx)

			if tt.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult, result)
			}

			assert.Equal(t, tt.expectedFound, found)
		})
	}
}

func TestGetDailyInfo(r *testing.T) {
	ctx := context.Background()
	now := time.Now()

	priceExample := Price{
		DateTime: now,
		Price:    1.0,
	}

	tests := []struct {
		name           string
		mockPrices     *[][]Price
		mockPricesErr  *[]error
		mockAvg        *[]float64
		mockAvgErr     *[]error
		expectedResult DailyPriceInfo
		expectingError bool
	}{
		{
			name:           "success",
			mockPrices:     &[][]Price{{priceExample}},
			mockPricesErr:  &[]error{},
			mockAvg:        &[]float64{1.0},
			mockAvgErr:     &[]error{},
			expectedResult: DailyPriceInfo{Prices: []Price{priceExample}, DayAverage: 1.0, DayRating: "NORMAL", ThirtyDayAverage: 1.0, CheapPeriods: [][]Price{{priceExample}}, ExpensivePeriods: [][]Price{{priceExample}}},
			expectingError: false,
		},
		{
			name:           "failure",
			mockPrices:     &[][]Price{{}},
			mockPricesErr:  &[]error{errors.New("not found")},
			mockAvg:        &[]float64{1.0},
			mockAvgErr:     &[]error{},
			expectedResult: DailyPriceInfo{Prices: []Price{}, DayAverage: 0.0, DayRating: "NORMAL", ThirtyDayAverage: 1.0, CheapPeriods: [][]Price{{}}, ExpensivePeriods: [][]Price{{}}},
			expectingError: true,
		},
	}

	for _, tt := range tests {
		r.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{MockFindResult: tt.mockPrices, MockFindErr: tt.mockPricesErr, MockThirtyDayAvg: tt.mockAvg, MockThirtyDayErr: tt.mockAvgErr}
			service := &Receiver{Collection: mockCollection}

			result, err := service.GetDailyInfo(ctx, now)

			if tt.expectingError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResult.Prices, result.Prices)
				assert.Equal(t, tt.expectedResult.DayAverage, result.DayAverage)
				assert.Equal(t, tt.expectedResult.DayRating, result.DayRating)
				assert.Equal(t, tt.expectedResult.ThirtyDayAverage, result.ThirtyDayAverage)
				assert.Equal(t, tt.expectedResult.CheapPeriods, result.CheapPeriods)
				assert.Equal(t, tt.expectedResult.ExpensivePeriods, result.ExpensivePeriods)

			}
		})
	}
}
