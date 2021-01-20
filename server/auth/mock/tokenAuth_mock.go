// Code generated by MockGen. DO NOT EDIT.
// Source: tokenAuth.go

// Package mock_auth is a generated GoMock package.
package mock_auth

import (
	jwt "github.com/dgrijalva/jwt-go"
	gomock "github.com/golang/mock/gomock"
	auth "go-app/server/auth"
	reflect "reflect"
)

// MockTokenAuth is a mock of TokenAuth interface
type MockTokenAuth struct {
	ctrl     *gomock.Controller
	recorder *MockTokenAuthMockRecorder
}

// MockTokenAuthMockRecorder is the mock recorder for MockTokenAuth
type MockTokenAuthMockRecorder struct {
	mock *MockTokenAuth
}

// NewMockTokenAuth creates a new mock instance
func NewMockTokenAuth(ctrl *gomock.Controller) *MockTokenAuth {
	mock := &MockTokenAuth{ctrl: ctrl}
	mock.recorder = &MockTokenAuthMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockTokenAuth) EXPECT() *MockTokenAuthMockRecorder {
	return m.recorder
}

// SignToken mocks base method
func (m *MockTokenAuth) SignToken() (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SignToken")
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SignToken indicates an expected call of SignToken
func (mr *MockTokenAuthMockRecorder) SignToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SignToken", reflect.TypeOf((*MockTokenAuth)(nil).SignToken))
}

// VerifyToken mocks base method
func (m *MockTokenAuth) VerifyToken(arg0 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "VerifyToken", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// VerifyToken indicates an expected call of VerifyToken
func (mr *MockTokenAuthMockRecorder) VerifyToken(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "VerifyToken", reflect.TypeOf((*MockTokenAuth)(nil).VerifyToken), arg0)
}

// GetClaim mocks base method
func (m *MockTokenAuth) GetClaim() auth.Claim {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClaim")
	ret0, _ := ret[0].(auth.Claim)
	return ret0
}

// GetClaim indicates an expected call of GetClaim
func (mr *MockTokenAuthMockRecorder) GetClaim() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClaim", reflect.TypeOf((*MockTokenAuth)(nil).GetClaim))
}

// SetClaim mocks base method
func (m *MockTokenAuth) SetClaim(arg0 auth.Claim) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetClaim", arg0)
}

// SetClaim indicates an expected call of SetClaim
func (mr *MockTokenAuthMockRecorder) SetClaim(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetClaim", reflect.TypeOf((*MockTokenAuth)(nil).SetClaim), arg0)
}

// MockClaim is a mock of Claim interface
type MockClaim struct {
	ctrl     *gomock.Controller
	recorder *MockClaimMockRecorder
}

// MockClaimMockRecorder is the mock recorder for MockClaim
type MockClaimMockRecorder struct {
	mock *MockClaim
}

// NewMockClaim creates a new mock instance
func NewMockClaim(ctrl *gomock.Controller) *MockClaim {
	mock := &MockClaim{ctrl: ctrl}
	mock.recorder = &MockClaimMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockClaim) EXPECT() *MockClaimMockRecorder {
	return m.recorder
}

// ToJSON mocks base method
func (m *MockClaim) ToJSON() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ToJSON")
	ret0, _ := ret[0].(string)
	return ret0
}

// ToJSON indicates an expected call of ToJSON
func (mr *MockClaimMockRecorder) ToJSON() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ToJSON", reflect.TypeOf((*MockClaim)(nil).ToJSON))
}

// GetJWTToken mocks base method
func (m *MockClaim) GetJWTToken() *jwt.Token {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetJWTToken")
	ret0, _ := ret[0].(*jwt.Token)
	return ret0
}

// GetJWTToken indicates an expected call of GetJWTToken
func (mr *MockClaimMockRecorder) GetJWTToken() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetJWTToken", reflect.TypeOf((*MockClaim)(nil).GetJWTToken))
}

// IsAdmin mocks base method
func (m *MockClaim) IsAdmin() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsAdmin")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsAdmin indicates an expected call of IsAdmin
func (mr *MockClaimMockRecorder) IsAdmin() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsAdmin", reflect.TypeOf((*MockClaim)(nil).IsAdmin))
}
