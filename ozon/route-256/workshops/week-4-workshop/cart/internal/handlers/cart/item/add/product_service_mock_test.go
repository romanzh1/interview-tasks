package add

// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

//go:generate minimock -i week-4-workshop/cart/internal/handlers/cart/item/add.productService -o ./cart/internal/handlers/cart/item/add/product_service_mock_test.go -n ProductServiceMock

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"
	"week-4-workshop/cart/internal/domain"

	"github.com/gojuno/minimock/v3"
)

// ProductServiceMock implements productService
type ProductServiceMock struct {
	t minimock.Tester

	funcGetProductInfo          func(ctx context.Context, sku uint32) (pp1 *domain.Product, err error)
	inspectFuncGetProductInfo   func(ctx context.Context, sku uint32)
	afterGetProductInfoCounter  uint64
	beforeGetProductInfoCounter uint64
	GetProductInfoMock          mProductServiceMockGetProductInfo
}

// NewProductServiceMock returns a mock for productService
func NewProductServiceMock(t minimock.Tester) *ProductServiceMock {
	m := &ProductServiceMock{t: t}
	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetProductInfoMock = mProductServiceMockGetProductInfo{mock: m}
	m.GetProductInfoMock.callArgs = []*ProductServiceMockGetProductInfoParams{}

	return m
}

type mProductServiceMockGetProductInfo struct {
	mock               *ProductServiceMock
	defaultExpectation *ProductServiceMockGetProductInfoExpectation
	expectations       []*ProductServiceMockGetProductInfoExpectation

	callArgs []*ProductServiceMockGetProductInfoParams
	mutex    sync.RWMutex
}

// ProductServiceMockGetProductInfoExpectation specifies expectation struct of the productService.GetProductInfo
type ProductServiceMockGetProductInfoExpectation struct {
	mock    *ProductServiceMock
	params  *ProductServiceMockGetProductInfoParams
	results *ProductServiceMockGetProductInfoResults
	Counter uint64
}

// ProductServiceMockGetProductInfoParams contains parameters of the productService.GetProductInfo
type ProductServiceMockGetProductInfoParams struct {
	ctx context.Context
	sku uint32
}

// ProductServiceMockGetProductInfoResults contains results of the productService.GetProductInfo
type ProductServiceMockGetProductInfoResults struct {
	pp1 *domain.Product
	err error
}

// Expect sets up expected params for productService.GetProductInfo
func (mmGetProductInfo *mProductServiceMockGetProductInfo) Expect(ctx context.Context, sku uint32) *mProductServiceMockGetProductInfo {
	if mmGetProductInfo.mock.funcGetProductInfo != nil {
		mmGetProductInfo.mock.t.Fatalf("ProductServiceMock.GetProductInfo mock is already set by Set")
	}

	if mmGetProductInfo.defaultExpectation == nil {
		mmGetProductInfo.defaultExpectation = &ProductServiceMockGetProductInfoExpectation{}
	}

	mmGetProductInfo.defaultExpectation.params = &ProductServiceMockGetProductInfoParams{ctx, sku}
	for _, e := range mmGetProductInfo.expectations {
		if minimock.Equal(e.params, mmGetProductInfo.defaultExpectation.params) {
			mmGetProductInfo.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetProductInfo.defaultExpectation.params)
		}
	}

	return mmGetProductInfo
}

// Inspect accepts an inspector function that has same arguments as the productService.GetProductInfo
func (mmGetProductInfo *mProductServiceMockGetProductInfo) Inspect(f func(ctx context.Context, sku uint32)) *mProductServiceMockGetProductInfo {
	if mmGetProductInfo.mock.inspectFuncGetProductInfo != nil {
		mmGetProductInfo.mock.t.Fatalf("Inspect function is already set for ProductServiceMock.GetProductInfo")
	}

	mmGetProductInfo.mock.inspectFuncGetProductInfo = f

	return mmGetProductInfo
}

// Return sets up results that will be returned by productService.GetProductInfo
func (mmGetProductInfo *mProductServiceMockGetProductInfo) Return(pp1 *domain.Product, err error) *ProductServiceMock {
	if mmGetProductInfo.mock.funcGetProductInfo != nil {
		mmGetProductInfo.mock.t.Fatalf("ProductServiceMock.GetProductInfo mock is already set by Set")
	}

	if mmGetProductInfo.defaultExpectation == nil {
		mmGetProductInfo.defaultExpectation = &ProductServiceMockGetProductInfoExpectation{mock: mmGetProductInfo.mock}
	}
	mmGetProductInfo.defaultExpectation.results = &ProductServiceMockGetProductInfoResults{pp1, err}
	return mmGetProductInfo.mock
}

// Set uses given function f to mock the productService.GetProductInfo method
func (mmGetProductInfo *mProductServiceMockGetProductInfo) Set(f func(ctx context.Context, sku uint32) (pp1 *domain.Product, err error)) *ProductServiceMock {
	if mmGetProductInfo.defaultExpectation != nil {
		mmGetProductInfo.mock.t.Fatalf("Default expectation is already set for the productService.GetProductInfo method")
	}

	if len(mmGetProductInfo.expectations) > 0 {
		mmGetProductInfo.mock.t.Fatalf("Some expectations are already set for the productService.GetProductInfo method")
	}

	mmGetProductInfo.mock.funcGetProductInfo = f
	return mmGetProductInfo.mock
}

// When sets expectation for the productService.GetProductInfo which will trigger the result defined by the following
// Then helper
func (mmGetProductInfo *mProductServiceMockGetProductInfo) When(ctx context.Context, sku uint32) *ProductServiceMockGetProductInfoExpectation {
	if mmGetProductInfo.mock.funcGetProductInfo != nil {
		mmGetProductInfo.mock.t.Fatalf("ProductServiceMock.GetProductInfo mock is already set by Set")
	}

	expectation := &ProductServiceMockGetProductInfoExpectation{
		mock:   mmGetProductInfo.mock,
		params: &ProductServiceMockGetProductInfoParams{ctx, sku},
	}
	mmGetProductInfo.expectations = append(mmGetProductInfo.expectations, expectation)
	return expectation
}

// Then sets up productService.GetProductInfo return parameters for the expectation previously defined by the When method
func (e *ProductServiceMockGetProductInfoExpectation) Then(pp1 *domain.Product, err error) *ProductServiceMock {
	e.results = &ProductServiceMockGetProductInfoResults{pp1, err}
	return e.mock
}

// GetProductInfo implements productService
func (mmGetProductInfo *ProductServiceMock) GetProductInfo(ctx context.Context, sku uint32) (pp1 *domain.Product, err error) {
	mm_atomic.AddUint64(&mmGetProductInfo.beforeGetProductInfoCounter, 1)
	defer mm_atomic.AddUint64(&mmGetProductInfo.afterGetProductInfoCounter, 1)

	if mmGetProductInfo.inspectFuncGetProductInfo != nil {
		mmGetProductInfo.inspectFuncGetProductInfo(ctx, sku)
	}

	mm_params := &ProductServiceMockGetProductInfoParams{ctx, sku}

	// Record call args
	mmGetProductInfo.GetProductInfoMock.mutex.Lock()
	mmGetProductInfo.GetProductInfoMock.callArgs = append(mmGetProductInfo.GetProductInfoMock.callArgs, mm_params)
	mmGetProductInfo.GetProductInfoMock.mutex.Unlock()

	for _, e := range mmGetProductInfo.GetProductInfoMock.expectations {
		if minimock.Equal(e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.pp1, e.results.err
		}
	}

	if mmGetProductInfo.GetProductInfoMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetProductInfo.GetProductInfoMock.defaultExpectation.Counter, 1)
		mm_want := mmGetProductInfo.GetProductInfoMock.defaultExpectation.params
		mm_got := ProductServiceMockGetProductInfoParams{ctx, sku}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetProductInfo.t.Errorf("ProductServiceMock.GetProductInfo got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetProductInfo.GetProductInfoMock.defaultExpectation.results
		if mm_results == nil {
			mmGetProductInfo.t.Fatal("No results are set for the ProductServiceMock.GetProductInfo")
		}
		return (*mm_results).pp1, (*mm_results).err
	}
	if mmGetProductInfo.funcGetProductInfo != nil {
		return mmGetProductInfo.funcGetProductInfo(ctx, sku)
	}
	mmGetProductInfo.t.Fatalf("Unexpected call to ProductServiceMock.GetProductInfo. %v %v", ctx, sku)
	return
}

// GetProductInfoAfterCounter returns a count of finished ProductServiceMock.GetProductInfo invocations
func (mmGetProductInfo *ProductServiceMock) GetProductInfoAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetProductInfo.afterGetProductInfoCounter)
}

// GetProductInfoBeforeCounter returns a count of ProductServiceMock.GetProductInfo invocations
func (mmGetProductInfo *ProductServiceMock) GetProductInfoBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetProductInfo.beforeGetProductInfoCounter)
}

// Calls returns a list of arguments used in each call to ProductServiceMock.GetProductInfo.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetProductInfo *mProductServiceMockGetProductInfo) Calls() []*ProductServiceMockGetProductInfoParams {
	mmGetProductInfo.mutex.RLock()

	argCopy := make([]*ProductServiceMockGetProductInfoParams, len(mmGetProductInfo.callArgs))
	copy(argCopy, mmGetProductInfo.callArgs)

	mmGetProductInfo.mutex.RUnlock()

	return argCopy
}

// MinimockGetProductInfoDone returns true if the count of the GetProductInfo invocations corresponds
// the number of defined expectations
func (m *ProductServiceMock) MinimockGetProductInfoDone() bool {
	for _, e := range m.GetProductInfoMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetProductInfoMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetProductInfoCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetProductInfo != nil && mm_atomic.LoadUint64(&m.afterGetProductInfoCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetProductInfoInspect logs each unmet expectation
func (m *ProductServiceMock) MinimockGetProductInfoInspect() {
	for _, e := range m.GetProductInfoMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ProductServiceMock.GetProductInfo with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetProductInfoMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetProductInfoCounter) < 1 {
		if m.GetProductInfoMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ProductServiceMock.GetProductInfo")
		} else {
			m.t.Errorf("Expected call to ProductServiceMock.GetProductInfo with params: %#v", *m.GetProductInfoMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetProductInfo != nil && mm_atomic.LoadUint64(&m.afterGetProductInfoCounter) < 1 {
		m.t.Error("Expected call to ProductServiceMock.GetProductInfo")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ProductServiceMock) MinimockFinish() {
	if !m.minimockDone() {
		m.MinimockGetProductInfoInspect()
		m.t.FailNow()
	}
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ProductServiceMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *ProductServiceMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetProductInfoDone()
}
