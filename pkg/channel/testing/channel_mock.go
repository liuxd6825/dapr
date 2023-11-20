// Code generated by mockery v2.14.0.

package testing

import (
	"context"
	"sync"

	mock "github.com/stretchr/testify/mock"

	"github.com/liuxd6825/dapr/pkg/apphealth"
	"github.com/liuxd6825/dapr/pkg/config"
	invokev1 "github.com/liuxd6825/dapr/pkg/messaging/v1"
)

// MockAppChannel is an autogenerated mock type for the AppChannel type
type MockAppChannel struct {
	mock.Mock
	requestsReceived map[string][]byte
	mutex            sync.Mutex
}

// GetAppConfig provides a mock function with given fields:
func (_m *MockAppChannel) GetAppConfig(_ context.Context, _ string) (*config.ApplicationConfig, error) {
	ret := _m.Called()

	var r0 *config.ApplicationConfig
	if rf, ok := ret.Get(0).(func() *config.ApplicationConfig); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*config.ApplicationConfig)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// HealthProbe provides a mock function with given fields: ctx
func (_m *MockAppChannel) HealthProbe(ctx context.Context) (bool, error) {
	ret := _m.Called(ctx)

	var r0 bool
	if rf, ok := ret.Get(0).(func(context.Context) bool); ok {
		r0 = rf(ctx)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockAppChannel) Init() {
	_m.mutex.Lock()
	_m.requestsReceived = make(map[string][]byte)
	_m.mutex.Unlock()
}

// InvokeMethod provides a mock function with given fields: ctx, req, appID
func (_m *MockAppChannel) InvokeMethod(ctx context.Context, req *invokev1.InvokeMethodRequest, appID string) (*invokev1.InvokeMethodResponse, error) {
	_m.mutex.Lock()
	if _m.requestsReceived != nil {
		req.WithReplay(true)
		pd, err := req.ProtoWithData()
		if err != nil {
			return nil, err
		}
		var data []byte
		if pd != nil && pd.Message != nil && pd.Message.Data != nil {
			data = pd.Message.Data.Value
		}
		_m.requestsReceived[req.Message().Method] = data
	}
	_m.mutex.Unlock()
	ret := _m.Called(ctx, req)

	var r0 *invokev1.InvokeMethodResponse
	if rf, ok := ret.Get(0).(func(context.Context, *invokev1.InvokeMethodRequest, string) *invokev1.InvokeMethodResponse); ok {
		r0 = rf(ctx, req, appID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*invokev1.InvokeMethodResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *invokev1.InvokeMethodRequest) error); ok {
		r1 = rf(ctx, req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *MockAppChannel) GetInvokedRequest() map[string][]byte {
	_m.mutex.Lock()
	defer _m.mutex.Unlock()
	return _m.requestsReceived
}

// SetAppHealth provides a mock function with given fields: ah
func (_m *MockAppChannel) SetAppHealth(ah *apphealth.AppHealth) {
	_m.Called(ah)
}

type mockConstructorTestingTNewMockAppChannel interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockAppChannel creates a new instance of MockAppChannel. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockAppChannel(t mockConstructorTestingTNewMockAppChannel) *MockAppChannel {
	mock := &MockAppChannel{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
