package price

import (
	"context"
	"electricity-prices/pkg/date"
	"errors"
	"testing"
	"time"

	"electricity-prices/pkg/db"
	"github.com/stretchr/testify/assert"
)

// MockCollection is a mock implementation of db.Collection[Price]
type MockCollection struct {
	db.Collection[Price]
	mockFindOneResult Price
	mockFindOneErr    error
	mockFindResult    []Price
	mockFindErr       error
	mockInsertManyErr error
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}) (Price, error) {
	return m.mockFindOneResult, m.mockFindOneErr
}

func (m *MockCollection) Find(ctx context.Context, filter interface{}) ([]Price, error) {
	return m.mockFindResult, m.mockFindErr
}

func (m *MockCollection) InsertMany(ctx context.Context, documents []Price) error {
	return m.mockInsertManyErr
}

func TestGetPrice(t *testing.T) {
	ctx := context.Background()
	now := time.Now()

	priceExample := Price{
		DateTime: now,
		Price:    1.0,
	}

	tests := []struct {
		name           string
		mockResult     Price
		mockError      error
		expectedResult Price
		expectingError bool
	}{
		{
			name:           "success",
			mockResult:     priceExample,
			mockError:      nil,
			expectedResult: priceExample,
			expectingError: false,
		},
		{
			name:           "failure",
			mockResult:     Price{},
			mockError:      errors.New("not found"),
			expectedResult: Price{},
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{mockFindOneResult: tt.mockResult, mockFindOneErr: tt.mockError}
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
		mockResult     []Price
		mockError      error
		expectedResult []Price
		expectingError bool
	}{
		{
			name:           "success",
			mockResult:     []Price{priceExample},
			mockError:      nil,
			expectedResult: []Price{priceExample},
			expectingError: false,
		},
		{
			name:           "failure",
			mockResult:     nil,
			mockError:      errors.New("not found"),
			expectedResult: []Price{},
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{mockFindResult: tt.mockResult, mockFindErr: tt.mockError}
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
		mockResult     []Price
		mockError      error
		expectedResult []Price
		expectingError bool
	}{
		{
			name:           "success",
			mockResult:     []Price{priceExample},
			mockError:      nil,
			expectedResult: []Price{priceExample},
			expectingError: false,
		},
		{
			name:           "failure",
			mockResult:     nil,
			mockError:      errors.New("not found"),
			expectedResult: []Price{},
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{mockFindResult: tt.mockResult, mockFindErr: tt.mockError}
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
		mockResult     []Price
		mockError      error
		expectedResult []DailyAverage
		expectingError bool
	}{
		{
			name:           "success",
			mockResult:     []Price{priceExample},
			mockError:      nil,
			expectedResult: []DailyAverage{dailyAverageExample},
			expectingError: false,
		},
		{
			name:           "failure",
			mockResult:     nil,
			mockError:      errors.New("not found"),
			expectedResult: []DailyAverage{},
			expectingError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockCollection := &MockCollection{mockFindResult: tt.mockResult, mockFindErr: tt.mockError}
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
