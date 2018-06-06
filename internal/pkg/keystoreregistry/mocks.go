// Code generated by MockGen. DO NOT EDIT.
// Source: internal/pkg/keystoreregistry/keystore.go

// Package keystoreregistry is a generated GoMock package.
package keystoreregistry

import (
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockKeystore is a mock of Keystore interface
type MockKeystore struct {
	ctrl     *gomock.Controller
	recorder *MockKeystoreMockRecorder
}

// MockKeystoreMockRecorder is the mock recorder for MockKeystore
type MockKeystoreMockRecorder struct {
	mock *MockKeystore
}

// NewMockKeystore creates a new mock instance
func NewMockKeystore(ctrl *gomock.Controller) *MockKeystore {
	mock := &MockKeystore{ctrl: ctrl}
	mock.recorder = &MockKeystoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockKeystore) EXPECT() *MockKeystoreMockRecorder {
	return m.recorder
}

// Close mocks base method
func (m *MockKeystore) Close() error {
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close
func (mr *MockKeystoreMockRecorder) Close() *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockKeystore)(nil).Close))
}

// CleanPrefix mocks base method
func (m *MockKeystore) CleanPrefix(prefix string) error {
	ret := m.ctrl.Call(m, "CleanPrefix", prefix)
	ret0, _ := ret[0].(error)
	return ret0
}

// CleanPrefix indicates an expected call of CleanPrefix
func (mr *MockKeystoreMockRecorder) CleanPrefix(prefix interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CleanPrefix", reflect.TypeOf((*MockKeystore)(nil).CleanPrefix), prefix)
}

// Add mocks base method
func (m *MockKeystore) Add(keyValues []KeyValue) error {
	ret := m.ctrl.Call(m, "Add", keyValues)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add
func (mr *MockKeystoreMockRecorder) Add(keyValues interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockKeystore)(nil).Add), keyValues)
}

// Update mocks base method
func (m *MockKeystore) Update(keyValues []KeyValueVersion) error {
	ret := m.ctrl.Call(m, "Update", keyValues)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update
func (mr *MockKeystoreMockRecorder) Update(keyValues interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockKeystore)(nil).Update), keyValues)
}

// DeleteAll mocks base method
func (m *MockKeystore) DeleteAll(keyValues []KeyValueVersion) error {
	ret := m.ctrl.Call(m, "DeleteAll", keyValues)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteAll indicates an expected call of DeleteAll
func (mr *MockKeystoreMockRecorder) DeleteAll(keyValues interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteAll", reflect.TypeOf((*MockKeystore)(nil).DeleteAll), keyValues)
}

// GetAll mocks base method
func (m *MockKeystore) GetAll(prefix string) ([]KeyValueVersion, error) {
	ret := m.ctrl.Call(m, "GetAll", prefix)
	ret0, _ := ret[0].([]KeyValueVersion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll
func (mr *MockKeystoreMockRecorder) GetAll(prefix interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockKeystore)(nil).GetAll), prefix)
}

// Get mocks base method
func (m *MockKeystore) Get(key string) (KeyValueVersion, error) {
	ret := m.ctrl.Call(m, "Get", key)
	ret0, _ := ret[0].(KeyValueVersion)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockKeystoreMockRecorder) Get(key interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockKeystore)(nil).Get), key)
}

// WatchPrefix mocks base method
func (m *MockKeystore) WatchPrefix(prefix string, onUpdate func(*KeyValueVersion, *KeyValueVersion)) {
	m.ctrl.Call(m, "WatchPrefix", prefix, onUpdate)
}

// WatchPrefix indicates an expected call of WatchPrefix
func (mr *MockKeystoreMockRecorder) WatchPrefix(prefix, onUpdate interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchPrefix", reflect.TypeOf((*MockKeystore)(nil).WatchPrefix), prefix, onUpdate)
}

// AtomicAdd mocks base method
func (m *MockKeystore) AtomicAdd(key, value string) {
	m.ctrl.Call(m, "AtomicAdd", key, value)
}

// AtomicAdd indicates an expected call of AtomicAdd
func (mr *MockKeystoreMockRecorder) AtomicAdd(key, value interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AtomicAdd", reflect.TypeOf((*MockKeystore)(nil).AtomicAdd), key, value)
}

// WatchPutPrefix mocks base method
func (m *MockKeystore) WatchPutPrefix(prefix string, onPut func(string, string)) {
	m.ctrl.Call(m, "WatchPutPrefix", prefix, onPut)
}

// WatchPutPrefix indicates an expected call of WatchPutPrefix
func (mr *MockKeystoreMockRecorder) WatchPutPrefix(prefix, onPut interface{}) *gomock.Call {
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WatchPutPrefix", reflect.TypeOf((*MockKeystore)(nil).WatchPutPrefix), prefix, onPut)
}