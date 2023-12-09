package testdata

import (
	"context"
	"electricity-prices/pkg/price"
	"time"
)

// Define a struct with a counter field
type CallCounter struct {
	Count int
}

// A mock implementation of Service

type MockPriceService struct {
	MockGetLatestPriceResult      *[]price.Price
	MockGetLatestPriceNoResult    *[]bool
	MockGetLatestPriceError       *[]error
	MockGetPriceResult            *[]price.Price
	MockGetPriceError             *[]error
	MockGetPricesResult           *[][]price.Price
	MockGetPricesError            *[]error
	MockSavePricesCount           *CallCounter
	MockSavePricesError           *[]error
	MockGetDailyPricesResult      *[][]price.Price
	MockGetDailyPricesError       *[]error
	MockGetDailyAveragesResult    *[][]price.DailyAverage
	MockGetDailyAveragesError     *[]error
	MockGetDailyInfoResult        *[]price.DailyPriceInfo
	MockGetDailyInfoError         *[]error
	MockGetDayRatingResult        *[]price.DayRating
	MockGetDayRatingError         *[]error
	MockGetDayAverageResult       *[]float64
	MockGetDayAverageError        *[]error
	MockGetCheapPeriodsResult     *[][][]price.Price
	MockGetCheapPeriodsError      *[]error
	MockGetExpensivePeriodsResult *[][][]price.Price
	MockGetExpensivePeriodsError  *[]error
	MockGetThirtyDayAverageResult *[]float64
	MockGetThirtyDayAverageError  *[]error
}

func (m *MockPriceService) GetLatestPrice(ctx context.Context) (price.Price, bool, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result price.Price
	if len(*m.MockGetLatestPriceResult) > 0 {
		result = (*m.MockGetLatestPriceResult)[0]
		*m.MockGetLatestPriceResult = (*m.MockGetLatestPriceResult)[1:]
	} else {
		result = price.Price{}
	}

	var noResult bool
	if len(*m.MockGetLatestPriceNoResult) > 0 {
		noResult = (*m.MockGetLatestPriceNoResult)[0]
		*m.MockGetLatestPriceNoResult = (*m.MockGetLatestPriceNoResult)[1:]
	} else {
		noResult = false
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockGetLatestPriceError) > 0 {
		err = (*m.MockGetLatestPriceError)[0]
		*m.MockGetLatestPriceError = (*m.MockGetLatestPriceError)[1:]
	} else {
		err = nil
	}

	return result, noResult, err
}

func (m *MockPriceService) GetPrice(ctx context.Context, t time.Time) (price.Price, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result price.Price
	if len(*m.MockGetPriceResult) > 0 {
		result = (*m.MockGetPriceResult)[0]
		*m.MockGetPriceResult = (*m.MockGetPriceResult)[1:]
	} else {
		result = price.Price{}
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockGetPriceError) > 0 {
		err = (*m.MockGetPriceError)[0]
		*m.MockGetPriceError = (*m.MockGetPriceError)[1:]
	} else {
		err = nil
	}

	return result, err
}

func (m *MockPriceService) GetPrices(ctx context.Context, start time.Time, end time.Time) ([]price.Price, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result []price.Price
	if len(*m.MockGetPricesResult) > 0 {
		result = (*m.MockGetPricesResult)[0]
		*m.MockGetPricesResult = (*m.MockGetPricesResult)[1:]
	} else {
		result = nil
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockGetPricesError) > 0 {
		err = (*m.MockGetPricesError)[0]
		*m.MockGetPricesError = (*m.MockGetPricesError)[1:]
	} else {
		err = nil
	}

	return result, err
}

func (m *MockPriceService) SavePrices(ctx context.Context, prices []price.Price) error {
	// Decrement the counter
	m.MockSavePricesCount.Count--

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockSavePricesError) > 0 {
		err = (*m.MockSavePricesError)[0]
		*m.MockSavePricesError = (*m.MockSavePricesError)[1:]
	} else {
		err = nil
	}

	return err
}

func (m *MockPriceService) GetDailyPrices(ctx context.Context, t time.Time) ([]price.Price, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result []price.Price
	if len(*m.MockGetDailyPricesResult) > 0 {
		result = (*m.MockGetDailyPricesResult)[0]
		*m.MockGetDailyPricesResult = (*m.MockGetDailyPricesResult)[1:]
	} else {
		result = nil
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockGetDailyPricesError) > 0 {
		err = (*m.MockGetDailyPricesError)[0]
		*m.MockGetDailyPricesError = (*m.MockGetDailyPricesError)[1:]
	} else {
		err = nil
	}

	return result, err
}

func (m *MockPriceService) GetDailyAverages(ctx context.Context, t time.Time, numberOfDays int) ([]price.DailyAverage, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result []price.DailyAverage
	if len(*m.MockGetDailyAveragesResult) > 0 {
		result = (*m.MockGetDailyAveragesResult)[0]
		*m.MockGetDailyAveragesResult = (*m.MockGetDailyAveragesResult)[1:]
	} else {
		result = nil
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockGetDailyAveragesError) > 0 {
		err = (*m.MockGetDailyAveragesError)[0]
		*m.MockGetDailyAveragesError = (*m.MockGetDailyAveragesError)[1:]
	} else {
		err = nil
	}

	return result, err
}

func (m *MockPriceService) GetDailyInfo(ctx context.Context, t time.Time) (price.DailyPriceInfo, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result price.DailyPriceInfo
	if len(*m.MockGetDailyInfoResult) > 0 {
		result = (*m.MockGetDailyInfoResult)[0]
		*m.MockGetDailyInfoResult = (*m.MockGetDailyInfoResult)[1:]
	} else {
		result = price.DailyPriceInfo{}
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockGetDailyInfoError) > 0 {
		err = (*m.MockGetDailyInfoError)[0]
		*m.MockGetDailyInfoError = (*m.MockGetDailyInfoError)[1:]
	} else {
		err = nil
	}

	return result, err
}

func (m *MockPriceService) GetDayRating(ctx context.Context, t time.Time) (price.DayRating, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result price.DayRating
	if len(*m.MockGetDayRatingResult) > 0 {
		result = (*m.MockGetDayRatingResult)[0]
		*m.MockGetDayRatingResult = (*m.MockGetDayRatingResult)[1:]
	} else {
		result = price.Nil
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockGetDayRatingError) > 0 {
		err = (*m.MockGetDayRatingError)[0]
		*m.MockGetDayRatingError = (*m.MockGetDayRatingError)[1:]
	} else {
		err = nil
	}

	return result, err
}

func (m *MockPriceService) GetDayAverage(ctx context.Context, t time.Time) (float64, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result float64
	if len(*m.MockGetDayAverageResult) > 0 {
		result = (*m.MockGetDayAverageResult)[0]
		*m.MockGetDayAverageResult = (*m.MockGetDayAverageResult)[1:]
	} else {
		result = 0.0
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockGetDayAverageError) > 0 {
		err = (*m.MockGetDayAverageError)[0]
		*m.MockGetDayAverageError = (*m.MockGetDayAverageError)[1:]
	} else {
		err = nil
	}

	return result, err
}

func (m *MockPriceService) GetCheapPeriods(ctx context.Context, t time.Time) ([][]price.Price, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result [][]price.Price
	if len(*m.MockGetCheapPeriodsResult) > 0 {
		result = (*m.MockGetCheapPeriodsResult)[0]
		*m.MockGetCheapPeriodsResult = (*m.MockGetCheapPeriodsResult)[1:]
	} else {
		result = nil
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockGetCheapPeriodsError) > 0 {
		err = (*m.MockGetCheapPeriodsError)[0]
		*m.MockGetCheapPeriodsError = (*m.MockGetCheapPeriodsError)[1:]
	} else {
		err = nil
	}

	return result, err
}

func (m *MockPriceService) GetExpensivePeriods(ctx context.Context, t time.Time) ([][]price.Price, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result [][]price.Price
	if len(*m.MockGetExpensivePeriodsResult) > 0 {
		result = (*m.MockGetExpensivePeriodsResult)[0]
		*m.MockGetExpensivePeriodsResult = (*m.MockGetExpensivePeriodsResult)[1:]
	} else {
		result = nil
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockGetExpensivePeriodsError) > 0 {
		err = (*m.MockGetExpensivePeriodsError)[0]
		*m.MockGetExpensivePeriodsError = (*m.MockGetExpensivePeriodsError)[1:]
	} else {
		err = nil
	}

	return result, err
}

func (m *MockPriceService) GetThirtyDayAverage(ctx context.Context, t time.Time) (float64, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result float64
	if len(*m.MockGetThirtyDayAverageResult) > 0 {
		result = (*m.MockGetThirtyDayAverageResult)[0]
		*m.MockGetThirtyDayAverageResult = (*m.MockGetThirtyDayAverageResult)[1:]
	} else {
		result = 0.0
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockGetThirtyDayAverageError) > 0 {
		err = (*m.MockGetThirtyDayAverageError)[0]
		*m.MockGetThirtyDayAverageError = (*m.MockGetThirtyDayAverageError)[1:]
	} else {
		err = nil
	}

	return result, err
}

// A mock implementation of Client

type MockPriceClient struct {
	MockGetPricesResult *[][]price.Price
	MockGetPricesSynced *[]bool
	MockGetPricesError  *[]error
}

func (m *MockPriceClient) GetPrices(t time.Time) ([]price.Price, bool, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result []price.Price
	if len(*m.MockGetPricesResult) > 0 {
		result = (*m.MockGetPricesResult)[0]
		*m.MockGetPricesResult = (*m.MockGetPricesResult)[1:]
	} else {
		result = nil
	}

	// Get the first element of the synced array and remove it from the array, return nil if the array is empty
	var synced bool
	if len(*m.MockGetPricesSynced) > 0 {
		synced = (*m.MockGetPricesSynced)[0]
		*m.MockGetPricesSynced = (*m.MockGetPricesSynced)[1:]
	} else {
		synced = false
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockGetPricesError) > 0 {
		err = (*m.MockGetPricesError)[0]
		*m.MockGetPricesError = (*m.MockGetPricesError)[1:]
	} else {
		err = nil
	}

	return result, synced, err
}
