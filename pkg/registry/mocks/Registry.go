// Copyright 2018 The ksonnet authors
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

// Code generated by mockery v1.0.0
package mocks

import app "github.com/ksonnet/ksonnet/metadata/app"
import mock "github.com/stretchr/testify/mock"
import parts "github.com/ksonnet/ksonnet/pkg/parts"
import registry "github.com/ksonnet/ksonnet/pkg/registry"

// Registry is an autogenerated mock type for the Registry type
type Registry struct {
	mock.Mock
}

// FetchRegistrySpec provides a mock function with given fields:
func (_m *Registry) FetchRegistrySpec() (*registry.Spec, error) {
	ret := _m.Called()

	var r0 *registry.Spec
	if rf, ok := ret.Get(0).(func() *registry.Spec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*registry.Spec)
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

// MakeRegistryRefSpec provides a mock function with given fields:
func (_m *Registry) MakeRegistryRefSpec() *app.RegistryRefSpec {
	ret := _m.Called()

	var r0 *app.RegistryRefSpec
	if rf, ok := ret.Get(0).(func() *app.RegistryRefSpec); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*app.RegistryRefSpec)
		}
	}

	return r0
}

// Name provides a mock function with given fields:
func (_m *Registry) Name() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// Protocol provides a mock function with given fields:
func (_m *Registry) Protocol() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// RegistrySpecDir provides a mock function with given fields:
func (_m *Registry) RegistrySpecDir() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// RegistrySpecFilePath provides a mock function with given fields:
func (_m *Registry) RegistrySpecFilePath() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ResolveLibrary provides a mock function with given fields: libID, libAlias, version, onFile, onDir
func (_m *Registry) ResolveLibrary(libID string, libAlias string, version string, onFile registry.ResolveFile, onDir registry.ResolveDirectory) (*parts.Spec, *app.LibraryRefSpec, error) {
	ret := _m.Called(libID, libAlias, version, onFile, onDir)

	var r0 *parts.Spec
	if rf, ok := ret.Get(0).(func(string, string, string, registry.ResolveFile, registry.ResolveDirectory) *parts.Spec); ok {
		r0 = rf(libID, libAlias, version, onFile, onDir)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*parts.Spec)
		}
	}

	var r1 *app.LibraryRefSpec
	if rf, ok := ret.Get(1).(func(string, string, string, registry.ResolveFile, registry.ResolveDirectory) *app.LibraryRefSpec); ok {
		r1 = rf(libID, libAlias, version, onFile, onDir)
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).(*app.LibraryRefSpec)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(string, string, string, registry.ResolveFile, registry.ResolveDirectory) error); ok {
		r2 = rf(libID, libAlias, version, onFile, onDir)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// ResolveLibrarySpec provides a mock function with given fields: libID, libRefSpec
func (_m *Registry) ResolveLibrarySpec(libID string, libRefSpec string) (*parts.Spec, error) {
	ret := _m.Called(libID, libRefSpec)

	var r0 *parts.Spec
	if rf, ok := ret.Get(0).(func(string, string) *parts.Spec); ok {
		r0 = rf(libID, libRefSpec)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*parts.Spec)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(libID, libRefSpec)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// URI provides a mock function with given fields:
func (_m *Registry) URI() string {
	ret := _m.Called()

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
