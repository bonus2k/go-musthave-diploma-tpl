// Code generated by MockGen. DO NOT EDIT.
// Source: \go-musthave-diploma-tpl\internal\repositories\storage.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	reflect "reflect"

	internal "github.com/bonus2k/go-musthave-diploma-tpl/internal"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// AddOrder mocks base method.
func (m *MockStore) AddOrder(ctx context.Context, order *internal.Order) (*internal.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrder", ctx, order)
	ret0, _ := ret[0].(*internal.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddOrder indicates an expected call of AddOrder.
func (mr *MockStoreMockRecorder) AddOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrder", reflect.TypeOf((*MockStore)(nil).AddOrder), ctx, order)
}

// AddUser mocks base method.
func (m *MockStore) AddUser(ctx context.Context, user *internal.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", ctx, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUser indicates an expected call of AddUser.
func (mr *MockStoreMockRecorder) AddUser(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockStore)(nil).AddUser), ctx, user)
}

// CheckConnection mocks base method.
func (m *MockStore) CheckConnection() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckConnection")
	ret0, _ := ret[0].(error)
	return ret0
}

// CheckConnection indicates an expected call of CheckConnection.
func (mr *MockStoreMockRecorder) CheckConnection() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckConnection", reflect.TypeOf((*MockStore)(nil).CheckConnection))
}

// FindUserByLogin mocks base method.
func (m *MockStore) FindUserByLogin(ctx context.Context, login string) (*internal.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FindUserByLogin", ctx, login)
	ret0, _ := ret[0].(*internal.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// FindUserByLogin indicates an expected call of FindUserByLogin.
func (mr *MockStoreMockRecorder) FindUserByLogin(ctx, login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FindUserByLogin", reflect.TypeOf((*MockStore)(nil).FindUserByLogin), ctx, login)
}

// GetOrders mocks base method.
func (m *MockStore) GetOrders(ctx context.Context, userID uuid.UUID) (*[]internal.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrders", ctx, userID)
	ret0, _ := ret[0].(*[]internal.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrders indicates an expected call of GetOrders.
func (mr *MockStoreMockRecorder) GetOrders(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrders", reflect.TypeOf((*MockStore)(nil).GetOrders), ctx, userID)
}

// GetOrdersNotProcessed mocks base method.
func (m *MockStore) GetOrdersNotProcessed(ctx context.Context) (*[]internal.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersNotProcessed", ctx)
	ret0, _ := ret[0].(*[]internal.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersNotProcessed indicates an expected call of GetOrdersNotProcessed.
func (mr *MockStoreMockRecorder) GetOrdersNotProcessed(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersNotProcessed", reflect.TypeOf((*MockStore)(nil).GetOrdersNotProcessed), ctx)
}

// GetUser mocks base method.
func (m *MockStore) GetUser(ctx context.Context, id uuid.UUID) (*internal.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, id)
	ret0, _ := ret[0].(*internal.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockStoreMockRecorder) GetUser(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockStore)(nil).GetUser), ctx, id)
}

// GetWithdrawals mocks base method.
func (m *MockStore) GetWithdrawals(ctx context.Context, userID uuid.UUID) (*[]internal.Withdraw, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawals", ctx, userID)
	ret0, _ := ret[0].(*[]internal.Withdraw)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawals indicates an expected call of GetWithdrawals.
func (mr *MockStoreMockRecorder) GetWithdrawals(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawals", reflect.TypeOf((*MockStore)(nil).GetWithdrawals), ctx, userID)
}

// SaveWithdrawal mocks base method.
func (m *MockStore) SaveWithdrawal(ctx context.Context, withdrawal *internal.Withdraw) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveWithdrawal", ctx, withdrawal)
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveWithdrawal indicates an expected call of SaveWithdrawal.
func (mr *MockStoreMockRecorder) SaveWithdrawal(ctx, withdrawal interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveWithdrawal", reflect.TypeOf((*MockStore)(nil).SaveWithdrawal), ctx, withdrawal)
}

// UpdateOrder mocks base method.
func (m *MockStore) UpdateOrder(ctx context.Context, order *internal.Order) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", ctx, order)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockStoreMockRecorder) UpdateOrder(ctx, order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockStore)(nil).UpdateOrder), ctx, order)
}