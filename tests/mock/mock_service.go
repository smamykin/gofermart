// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/service/contracts.go

// Package mock_contracts is a generated GoMock package.
package mock_contracts

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	entity "github.com/smamykin/gofermart/internal/entity"
	service "github.com/smamykin/gofermart/internal/service"
)

// MockUserRepositoryInterface is a mock of UserRepositoryInterface interface.
type MockUserRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepositoryInterfaceMockRecorder
}

// MockUserRepositoryInterfaceMockRecorder is the mock recorder for MockUserRepositoryInterface.
type MockUserRepositoryInterfaceMockRecorder struct {
	mock *MockUserRepositoryInterface
}

// NewMockUserRepositoryInterface creates a new mock instance.
func NewMockUserRepositoryInterface(ctrl *gomock.Controller) *MockUserRepositoryInterface {
	mock := &MockUserRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockUserRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepositoryInterface) EXPECT() *MockUserRepositoryInterfaceMockRecorder {
	return m.recorder
}

// GetUserByLogin mocks base method.
func (m *MockUserRepositoryInterface) GetUserByLogin(login string) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", login)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockUserRepositoryInterfaceMockRecorder) GetUserByLogin(login interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockUserRepositoryInterface)(nil).GetUserByLogin), login)
}

// UpsertUser mocks base method.
func (m *MockUserRepositoryInterface) UpsertUser(login, pwd string) (entity.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpsertUser", login, pwd)
	ret0, _ := ret[0].(entity.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpsertUser indicates an expected call of UpsertUser.
func (mr *MockUserRepositoryInterfaceMockRecorder) UpsertUser(login, pwd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpsertUser", reflect.TypeOf((*MockUserRepositoryInterface)(nil).UpsertUser), login, pwd)
}

// MockOrderRepositoryInterface is a mock of OrderRepositoryInterface interface.
type MockOrderRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockOrderRepositoryInterfaceMockRecorder
}

// MockOrderRepositoryInterfaceMockRecorder is the mock recorder for MockOrderRepositoryInterface.
type MockOrderRepositoryInterfaceMockRecorder struct {
	mock *MockOrderRepositoryInterface
}

// NewMockOrderRepositoryInterface creates a new mock instance.
func NewMockOrderRepositoryInterface(ctrl *gomock.Controller) *MockOrderRepositoryInterface {
	mock := &MockOrderRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockOrderRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockOrderRepositoryInterface) EXPECT() *MockOrderRepositoryInterfaceMockRecorder {
	return m.recorder
}

// AddOrder mocks base method.
func (m *MockOrderRepositoryInterface) AddOrder(o entity.Order) (entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddOrder", o)
	ret0, _ := ret[0].(entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddOrder indicates an expected call of AddOrder.
func (mr *MockOrderRepositoryInterfaceMockRecorder) AddOrder(o interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddOrder", reflect.TypeOf((*MockOrderRepositoryInterface)(nil).AddOrder), o)
}

// GetAccrualSumByUserId mocks base method.
func (m *MockOrderRepositoryInterface) GetAccrualSumByUserId(userID int) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccrualSumByUserId", userID)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAccrualSumByUserId indicates an expected call of GetAccrualSumByUserId.
func (mr *MockOrderRepositoryInterfaceMockRecorder) GetAccrualSumByUserId(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccrualSumByUserId", reflect.TypeOf((*MockOrderRepositoryInterface)(nil).GetAccrualSumByUserId), userID)
}

// GetAllByUserID mocks base method.
func (m *MockOrderRepositoryInterface) GetAllByUserID(userID int) ([]entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllByUserID", userID)
	ret0, _ := ret[0].([]entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllByUserID indicates an expected call of GetAllByUserID.
func (mr *MockOrderRepositoryInterfaceMockRecorder) GetAllByUserID(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllByUserID", reflect.TypeOf((*MockOrderRepositoryInterface)(nil).GetAllByUserID), userID)
}

// GetOrder mocks base method.
func (m *MockOrderRepositoryInterface) GetOrder(ID int) (entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrder", ID)
	ret0, _ := ret[0].(entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrder indicates an expected call of GetOrder.
func (mr *MockOrderRepositoryInterfaceMockRecorder) GetOrder(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrder", reflect.TypeOf((*MockOrderRepositoryInterface)(nil).GetOrder), ID)
}

// GetOrderByOrderNumber mocks base method.
func (m *MockOrderRepositoryInterface) GetOrderByOrderNumber(orderNumber string) (entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByOrderNumber", orderNumber)
	ret0, _ := ret[0].(entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByOrderNumber indicates an expected call of GetOrderByOrderNumber.
func (mr *MockOrderRepositoryInterfaceMockRecorder) GetOrderByOrderNumber(orderNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByOrderNumber", reflect.TypeOf((*MockOrderRepositoryInterface)(nil).GetOrderByOrderNumber), orderNumber)
}

// GetOrdersWithUnfinishedStatus mocks base method.
func (m *MockOrderRepositoryInterface) GetOrdersWithUnfinishedStatus() ([]entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersWithUnfinishedStatus")
	ret0, _ := ret[0].([]entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersWithUnfinishedStatus indicates an expected call of GetOrdersWithUnfinishedStatus.
func (mr *MockOrderRepositoryInterfaceMockRecorder) GetOrdersWithUnfinishedStatus() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersWithUnfinishedStatus", reflect.TypeOf((*MockOrderRepositoryInterface)(nil).GetOrdersWithUnfinishedStatus))
}

// UpdateOrder mocks base method.
func (m *MockOrderRepositoryInterface) UpdateOrder(order entity.Order) (entity.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateOrder", order)
	ret0, _ := ret[0].(entity.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateOrder indicates an expected call of UpdateOrder.
func (mr *MockOrderRepositoryInterfaceMockRecorder) UpdateOrder(order interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateOrder", reflect.TypeOf((*MockOrderRepositoryInterface)(nil).UpdateOrder), order)
}

// MockWithdrawalRepositoryInterface is a mock of WithdrawalRepositoryInterface interface.
type MockWithdrawalRepositoryInterface struct {
	ctrl     *gomock.Controller
	recorder *MockWithdrawalRepositoryInterfaceMockRecorder
}

// MockWithdrawalRepositoryInterfaceMockRecorder is the mock recorder for MockWithdrawalRepositoryInterface.
type MockWithdrawalRepositoryInterfaceMockRecorder struct {
	mock *MockWithdrawalRepositoryInterface
}

// NewMockWithdrawalRepositoryInterface creates a new mock instance.
func NewMockWithdrawalRepositoryInterface(ctrl *gomock.Controller) *MockWithdrawalRepositoryInterface {
	mock := &MockWithdrawalRepositoryInterface{ctrl: ctrl}
	mock.recorder = &MockWithdrawalRepositoryInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWithdrawalRepositoryInterface) EXPECT() *MockWithdrawalRepositoryInterfaceMockRecorder {
	return m.recorder
}

// AddWithdrawal mocks base method.
func (m *MockWithdrawalRepositoryInterface) AddWithdrawal(withdrawal entity.Withdrawal) (entity.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddWithdrawal", withdrawal)
	ret0, _ := ret[0].(entity.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddWithdrawal indicates an expected call of AddWithdrawal.
func (mr *MockWithdrawalRepositoryInterfaceMockRecorder) AddWithdrawal(withdrawal interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddWithdrawal", reflect.TypeOf((*MockWithdrawalRepositoryInterface)(nil).AddWithdrawal), withdrawal)
}

// GetAmountSumByUserId mocks base method.
func (m *MockWithdrawalRepositoryInterface) GetAmountSumByUserId(userID int) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAmountSumByUserId", userID)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAmountSumByUserId indicates an expected call of GetAmountSumByUserId.
func (mr *MockWithdrawalRepositoryInterfaceMockRecorder) GetAmountSumByUserId(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAmountSumByUserId", reflect.TypeOf((*MockWithdrawalRepositoryInterface)(nil).GetAmountSumByUserId), userID)
}

// GetWithdrawal mocks base method.
func (m *MockWithdrawalRepositoryInterface) GetWithdrawal(ID int) (entity.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawal", ID)
	ret0, _ := ret[0].(entity.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawal indicates an expected call of GetWithdrawal.
func (mr *MockWithdrawalRepositoryInterfaceMockRecorder) GetWithdrawal(ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawal", reflect.TypeOf((*MockWithdrawalRepositoryInterface)(nil).GetWithdrawal), ID)
}

// MockAccrualClientInterface is a mock of AccrualClientInterface interface.
type MockAccrualClientInterface struct {
	ctrl     *gomock.Controller
	recorder *MockAccrualClientInterfaceMockRecorder
}

// MockAccrualClientInterfaceMockRecorder is the mock recorder for MockAccrualClientInterface.
type MockAccrualClientInterfaceMockRecorder struct {
	mock *MockAccrualClientInterface
}

// NewMockAccrualClientInterface creates a new mock instance.
func NewMockAccrualClientInterface(ctrl *gomock.Controller) *MockAccrualClientInterface {
	mock := &MockAccrualClientInterface{ctrl: ctrl}
	mock.recorder = &MockAccrualClientInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAccrualClientInterface) EXPECT() *MockAccrualClientInterfaceMockRecorder {
	return m.recorder
}

// GetOrder mocks base method.
func (m *MockAccrualClientInterface) GetOrder(orderNumber string) (service.AccrualOrder, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrder", orderNumber)
	ret0, _ := ret[0].(service.AccrualOrder)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrder indicates an expected call of GetOrder.
func (mr *MockAccrualClientInterfaceMockRecorder) GetOrder(orderNumber interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrder", reflect.TypeOf((*MockAccrualClientInterface)(nil).GetOrder), orderNumber)
}

// MockHashGeneratorInterface is a mock of HashGeneratorInterface interface.
type MockHashGeneratorInterface struct {
	ctrl     *gomock.Controller
	recorder *MockHashGeneratorInterfaceMockRecorder
}

// MockHashGeneratorInterfaceMockRecorder is the mock recorder for MockHashGeneratorInterface.
type MockHashGeneratorInterfaceMockRecorder struct {
	mock *MockHashGeneratorInterface
}

// NewMockHashGeneratorInterface creates a new mock instance.
func NewMockHashGeneratorInterface(ctrl *gomock.Controller) *MockHashGeneratorInterface {
	mock := &MockHashGeneratorInterface{ctrl: ctrl}
	mock.recorder = &MockHashGeneratorInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockHashGeneratorInterface) EXPECT() *MockHashGeneratorInterfaceMockRecorder {
	return m.recorder
}

// Generate mocks base method.
func (m *MockHashGeneratorInterface) Generate(stringToHash string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Generate", stringToHash)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Generate indicates an expected call of Generate.
func (mr *MockHashGeneratorInterfaceMockRecorder) Generate(stringToHash interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Generate", reflect.TypeOf((*MockHashGeneratorInterface)(nil).Generate), stringToHash)
}

// IsEqual mocks base method.
func (m *MockHashGeneratorInterface) IsEqual(hashedPassword, plainTxtPwd string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsEqual", hashedPassword, plainTxtPwd)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// IsEqual indicates an expected call of IsEqual.
func (mr *MockHashGeneratorInterfaceMockRecorder) IsEqual(hashedPassword, plainTxtPwd interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsEqual", reflect.TypeOf((*MockHashGeneratorInterface)(nil).IsEqual), hashedPassword, plainTxtPwd)
}
