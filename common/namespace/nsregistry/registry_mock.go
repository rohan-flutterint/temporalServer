// Code generated by MockGen. DO NOT EDIT.
// Source: registry.go
//
// Generated by this command:
//
//	mockgen -package nsregistry -source registry.go -destination registry_mock.go
//

// Package nsregistry is a generated GoMock package.
package nsregistry

import (
	context "context"
	reflect "reflect"
	time "time"

	persistence "go.temporal.io/server/common/persistence"
	gomock "go.uber.org/mock/gomock"
)

// MockClock is a mock of Clock interface.
type MockClock struct {
	ctrl     *gomock.Controller
	recorder *MockClockMockRecorder
	isgomock struct{}
}

// MockClockMockRecorder is the mock recorder for MockClock.
type MockClockMockRecorder struct {
	mock *MockClock
}

// NewMockClock creates a new mock instance.
func NewMockClock(ctrl *gomock.Controller) *MockClock {
	mock := &MockClock{ctrl: ctrl}
	mock.recorder = &MockClockMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClock) EXPECT() *MockClockMockRecorder {
	return m.recorder
}

// Now mocks base method.
func (m *MockClock) Now() time.Time {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Now")
	ret0, _ := ret[0].(time.Time)
	return ret0
}

// Now indicates an expected call of Now.
func (mr *MockClockMockRecorder) Now() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Now", reflect.TypeOf((*MockClock)(nil).Now))
}

// MockPersistence is a mock of Persistence interface.
type MockPersistence struct {
	ctrl     *gomock.Controller
	recorder *MockPersistenceMockRecorder
	isgomock struct{}
}

// MockPersistenceMockRecorder is the mock recorder for MockPersistence.
type MockPersistenceMockRecorder struct {
	mock *MockPersistence
}

// NewMockPersistence creates a new mock instance.
func NewMockPersistence(ctrl *gomock.Controller) *MockPersistence {
	mock := &MockPersistence{ctrl: ctrl}
	mock.recorder = &MockPersistenceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPersistence) EXPECT() *MockPersistenceMockRecorder {
	return m.recorder
}

// GetMetadata mocks base method.
func (m *MockPersistence) GetMetadata(arg0 context.Context) (*persistence.GetMetadataResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetMetadata", arg0)
	ret0, _ := ret[0].(*persistence.GetMetadataResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetMetadata indicates an expected call of GetMetadata.
func (mr *MockPersistenceMockRecorder) GetMetadata(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetMetadata", reflect.TypeOf((*MockPersistence)(nil).GetMetadata), arg0)
}

// GetNamespace mocks base method.
func (m *MockPersistence) GetNamespace(arg0 context.Context, arg1 *persistence.GetNamespaceRequest) (*persistence.GetNamespaceResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetNamespace", arg0, arg1)
	ret0, _ := ret[0].(*persistence.GetNamespaceResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetNamespace indicates an expected call of GetNamespace.
func (mr *MockPersistenceMockRecorder) GetNamespace(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetNamespace", reflect.TypeOf((*MockPersistence)(nil).GetNamespace), arg0, arg1)
}

// ListNamespaces mocks base method.
func (m *MockPersistence) ListNamespaces(arg0 context.Context, arg1 *persistence.ListNamespacesRequest) (*persistence.ListNamespacesResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListNamespaces", arg0, arg1)
	ret0, _ := ret[0].(*persistence.ListNamespacesResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListNamespaces indicates an expected call of ListNamespaces.
func (mr *MockPersistenceMockRecorder) ListNamespaces(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListNamespaces", reflect.TypeOf((*MockPersistence)(nil).ListNamespaces), arg0, arg1)
}
