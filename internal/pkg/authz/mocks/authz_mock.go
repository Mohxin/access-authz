// Code generated by MockGen. DO NOT EDIT.
// Source: authz.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	cachemanager "github.com/volvo-cars/connect-access-control/internal/pkg/gateway/cache-manager"
	plums "github.com/volvo-cars/connect-access-control/internal/pkg/gateway/plums"
	store "github.com/volvo-cars/connect-access-control/internal/pkg/store"
)

// MockcacheClient is a mock of cacheClient interface.
type MockcacheClient struct {
	ctrl     *gomock.Controller
	recorder *MockcacheClientMockRecorder
}

// MockcacheClientMockRecorder is the mock recorder for MockcacheClient.
type MockcacheClientMockRecorder struct {
	mock *MockcacheClient
}

// NewMockcacheClient creates a new mock instance.
func NewMockcacheClient(ctrl *gomock.Controller) *MockcacheClient {
	mock := &MockcacheClient{ctrl: ctrl}
	mock.recorder = &MockcacheClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockcacheClient) EXPECT() *MockcacheClientMockRecorder {
	return m.recorder
}

// GetPartnersByCodes mocks base method.
func (m *MockcacheClient) GetPartnersByCodes(ctx context.Context, partnerCodes []string, partnerType string) ([]*cachemanager.Partner, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPartnersByCodes", ctx, partnerCodes, partnerType)
	ret0, _ := ret[0].([]*cachemanager.Partner)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPartnersByCodes indicates an expected call of GetPartnersByCodes.
func (mr *MockcacheClientMockRecorder) GetPartnersByCodes(ctx, partnerCodes, partnerType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPartnersByCodes", reflect.TypeOf((*MockcacheClient)(nil).GetPartnersByCodes), ctx, partnerCodes, partnerType)
}

// MockplumsClient is a mock of plumsClient interface.
type MockplumsClient struct {
	ctrl     *gomock.Controller
	recorder *MockplumsClientMockRecorder
}

// MockplumsClientMockRecorder is the mock recorder for MockplumsClient.
type MockplumsClientMockRecorder struct {
	mock *MockplumsClient
}

// NewMockplumsClient creates a new mock instance.
func NewMockplumsClient(ctrl *gomock.Controller) *MockplumsClient {
	mock := &MockplumsClient{ctrl: ctrl}
	mock.recorder = &MockplumsClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockplumsClient) EXPECT() *MockplumsClientMockRecorder {
	return m.recorder
}

// GetUserByCDSID mocks base method.
func (m *MockplumsClient) GetUserByCDSID(ctx context.Context, cdsid string) (*plums.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByCDSID", ctx, cdsid)
	ret0, _ := ret[0].(*plums.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByCDSID indicates an expected call of GetUserByCDSID.
func (mr *MockplumsClientMockRecorder) GetUserByCDSID(ctx, cdsid interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByCDSID", reflect.TypeOf((*MockplumsClient)(nil).GetUserByCDSID), ctx, cdsid)
}

// MockauthzStore is a mock of authzStore interface.
type MockauthzStore struct {
	ctrl     *gomock.Controller
	recorder *MockauthzStoreMockRecorder
}

// MockauthzStoreMockRecorder is the mock recorder for MockauthzStore.
type MockauthzStoreMockRecorder struct {
	mock *MockauthzStore
}

// NewMockauthzStore creates a new mock instance.
func NewMockauthzStore(ctrl *gomock.Controller) *MockauthzStore {
	mock := &MockauthzStore{ctrl: ctrl}
	mock.recorder = &MockauthzStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockauthzStore) EXPECT() *MockauthzStoreMockRecorder {
	return m.recorder
}

// GetRoleMapping mocks base method.
func (m *MockauthzStore) GetRoleMapping(scopeID, roleID string) (store.RoleMapping, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoleMapping", scopeID, roleID)
	ret0, _ := ret[0].(store.RoleMapping)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoleMapping indicates an expected call of GetRoleMapping.
func (mr *MockauthzStoreMockRecorder) GetRoleMapping(scopeID, roleID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoleMapping", reflect.TypeOf((*MockauthzStore)(nil).GetRoleMapping), scopeID, roleID)
}

// GetRoleMappings mocks base method.
func (m *MockauthzStore) GetRoleMappings(scopeID string) ([]store.RoleMapping, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRoleMappings", scopeID)
	ret0, _ := ret[0].([]store.RoleMapping)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRoleMappings indicates an expected call of GetRoleMappings.
func (mr *MockauthzStoreMockRecorder) GetRoleMappings(scopeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRoleMappings", reflect.TypeOf((*MockauthzStore)(nil).GetRoleMappings), scopeID)
}
