// Licensed to Elasticsearch B.V. under one or more contributor
// license agreements. See the NOTICE file distributed with
// this work for additional information regarding copyright
// ownership. Elasticsearch B.V. licenses this file to you under
// the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied.  See the License for the
// specific language governing permissions and limitations
// under the License.

// Code generated by mockery v2.24.0. DO NOT EDIT.

package identity

import (
	config "github.com/elastic/cloudbeat/config"
	cloud "github.com/elastic/cloudbeat/dataprovider/providers/cloud"

	context "context"

	mock "github.com/stretchr/testify/mock"
)

// MockIdentityProviderGetter is an autogenerated mock type for the ProviderGetter type
type MockIdentityProviderGetter struct {
	mock.Mock
}

type MockIdentityProviderGetter_Expecter struct {
	mock *mock.Mock
}

func (_m *MockIdentityProviderGetter) EXPECT() *MockIdentityProviderGetter_Expecter {
	return &MockIdentityProviderGetter_Expecter{mock: &_m.Mock}
}

// GetIdentity provides a mock function with given fields: ctx, cfg
func (_m *MockIdentityProviderGetter) GetIdentity(ctx context.Context, cfg config.GcpConfig) (*cloud.Identity, error) {
	ret := _m.Called(ctx, cfg)

	var r0 *cloud.Identity
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, config.GcpConfig) (*cloud.Identity, error)); ok {
		return rf(ctx, cfg)
	}
	if rf, ok := ret.Get(0).(func(context.Context, config.GcpConfig) *cloud.Identity); ok {
		r0 = rf(ctx, cfg)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cloud.Identity)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, config.GcpConfig) error); ok {
		r1 = rf(ctx, cfg)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// MockIdentityProviderGetter_GetIdentity_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetIdentity'
type MockIdentityProviderGetter_GetIdentity_Call struct {
	*mock.Call
}

// GetIdentity is a helper method to define mock.On call
//   - ctx context.Context
//   - cfg config.GcpConfig
func (_e *MockIdentityProviderGetter_Expecter) GetIdentity(ctx interface{}, cfg interface{}) *MockIdentityProviderGetter_GetIdentity_Call {
	return &MockIdentityProviderGetter_GetIdentity_Call{Call: _e.mock.On("GetIdentity", ctx, cfg)}
}

func (_c *MockIdentityProviderGetter_GetIdentity_Call) Run(run func(ctx context.Context, cfg config.GcpConfig)) *MockIdentityProviderGetter_GetIdentity_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(config.GcpConfig))
	})
	return _c
}

func (_c *MockIdentityProviderGetter_GetIdentity_Call) Return(_a0 *cloud.Identity, _a1 error) *MockIdentityProviderGetter_GetIdentity_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *MockIdentityProviderGetter_GetIdentity_Call) RunAndReturn(run func(context.Context, config.GcpConfig) (*cloud.Identity, error)) *MockIdentityProviderGetter_GetIdentity_Call {
	_c.Call.Return(run)
	return _c
}

type mockConstructorTestingTNewMockIdentityProviderGetter interface {
	mock.TestingT
	Cleanup(func())
}

// NewMockIdentityProviderGetter creates a new instance of MockIdentityProviderGetter. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewMockIdentityProviderGetter(t mockConstructorTestingTNewMockIdentityProviderGetter) *MockIdentityProviderGetter {
	mock := &MockIdentityProviderGetter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}