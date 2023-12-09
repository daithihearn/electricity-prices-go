package price

import (
	"context"
	"time"
)

type MockCollection struct {
	Collection
	MockFindOneResult  *[]Price
	MockFindOneErr     *[]error
	MockFindResult     *[][]Price
	MockFindErr        *[]error
	MockInsertManyErr  *[]error
	MockThirtyDayAvg   *[]float64
	MockThirtyDayErr   *[]error
	MockLatestPrice    *[]Price
	MockLatestPriceOk  *[]bool
	MockLatestPriceErr *[]error
}

func (m *MockCollection) FindOne(ctx context.Context, filter interface{}) (Price, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result Price
	if len(*m.MockFindOneResult) > 0 {
		result = (*m.MockFindOneResult)[0]
		*m.MockFindOneResult = (*m.MockFindOneResult)[1:]
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockFindOneErr) > 0 {
		err = (*m.MockFindOneErr)[0]
		*m.MockFindOneErr = (*m.MockFindOneErr)[1:]
	} else {
		err = nil
	}

	return result, err
}

func (m *MockCollection) Find(ctx context.Context, filter interface{}) ([]Price, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result []Price
	if len(*m.MockFindResult) > 0 {
		result = (*m.MockFindResult)[0]
		*m.MockFindResult = (*m.MockFindResult)[1:]
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockFindErr) > 0 {
		err = (*m.MockFindErr)[0]
		*m.MockFindErr = (*m.MockFindErr)[1:]
	} else {
		err = nil
	}

	return result, err
}

func (m *MockCollection) InsertMany(ctx context.Context, documents []Price) error {
	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockInsertManyErr) > 0 {
		err = (*m.MockInsertManyErr)[0]
		*m.MockInsertManyErr = (*m.MockInsertManyErr)[1:]
	} else {
		err = nil
	}

	return err
}

func (m *MockCollection) GetThirtyDayAverage(ctx context.Context, t time.Time) (float64, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result float64
	if len(*m.MockThirtyDayAvg) > 0 {
		result = (*m.MockThirtyDayAvg)[0]
		*m.MockThirtyDayAvg = (*m.MockThirtyDayAvg)[1:]
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockThirtyDayErr) > 0 {
		err = (*m.MockThirtyDayErr)[0]
		*m.MockThirtyDayErr = (*m.MockThirtyDayErr)[1:]
	} else {
		err = nil
	}

	return result, err
}

func (m *MockCollection) GetLatestPrice(ctx context.Context) (Price, bool, error) {
	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var result Price
	if len(*m.MockLatestPrice) > 0 {
		result = (*m.MockLatestPrice)[0]
		*m.MockLatestPrice = (*m.MockLatestPrice)[1:]
	}

	// Get the first element of the result array and remove it from the array, return nil if the array is empty
	var ok bool
	if len(*m.MockLatestPriceOk) > 0 {
		ok = (*m.MockLatestPriceOk)[0]
		*m.MockLatestPriceOk = (*m.MockLatestPriceOk)[1:]
	}

	// Get the first element of the error array and remove it from the array, return nil if the array is empty
	var err error
	if len(*m.MockLatestPriceErr) > 0 {
		err = (*m.MockLatestPriceErr)[0]
		*m.MockLatestPriceErr = (*m.MockLatestPriceErr)[1:]
	} else {
		err = nil
	}

	return result, ok, err
}
