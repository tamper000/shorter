// Code generated by MockGen. DO NOT EDIT.
// Source: database.go
//
// Generated by this command:
//
//	mockgen -source=database.go -destination=../mocks/database_mock.go --package mocks
//

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"
	models "urlshort/internal/models"

	gomock "go.uber.org/mock/gomock"
)

// MockDatabase is a mock of Database interface.
type MockDatabase struct {
	ctrl     *gomock.Controller
	recorder *MockDatabaseMockRecorder
	isgomock struct{}
}

// MockDatabaseMockRecorder is the mock recorder for MockDatabase.
type MockDatabaseMockRecorder struct {
	mock *MockDatabase
}

// NewMockDatabase creates a new mock instance.
func NewMockDatabase(ctrl *gomock.Controller) *MockDatabase {
	mock := &MockDatabase{ctrl: ctrl}
	mock.recorder = &MockDatabaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDatabase) EXPECT() *MockDatabaseMockRecorder {
	return m.recorder
}

// AddClick mocks base method.
func (m *MockDatabase) AddClick(alias string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddClick", alias)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddClick indicates an expected call of AddClick.
func (mr *MockDatabaseMockRecorder) AddClick(alias any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddClick", reflect.TypeOf((*MockDatabase)(nil).AddClick), alias)
}

// AddUser mocks base method.
func (m *MockDatabase) AddUser(username, password string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", username, password)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUser indicates an expected call of AddUser.
func (mr *MockDatabaseMockRecorder) AddUser(username, password any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockDatabase)(nil).AddUser), username, password)
}

// CheckUserExists mocks base method.
func (m *MockDatabase) CheckUserExists(username string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckUserExists", username)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckUserExists indicates an expected call of CheckUserExists.
func (mr *MockDatabaseMockRecorder) CheckUserExists(username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckUserExists", reflect.TypeOf((*MockDatabase)(nil).CheckUserExists), username)
}

// CloseDatabase mocks base method.
func (m *MockDatabase) CloseDatabase() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "CloseDatabase")
}

// CloseDatabase indicates an expected call of CloseDatabase.
func (mr *MockDatabaseMockRecorder) CloseDatabase() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseDatabase", reflect.TypeOf((*MockDatabase)(nil).CloseDatabase))
}

// DeleteLink mocks base method.
func (m *MockDatabase) DeleteLink(alias, username string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteLink", alias, username)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteLink indicates an expected call of DeleteLink.
func (mr *MockDatabaseMockRecorder) DeleteLink(alias, username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteLink", reflect.TypeOf((*MockDatabase)(nil).DeleteLink), alias, username)
}

// GetByAlias mocks base method.
func (m *MockDatabase) GetByAlias(alias string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByAlias", alias)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByAlias indicates an expected call of GetByAlias.
func (mr *MockDatabaseMockRecorder) GetByAlias(alias any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByAlias", reflect.TypeOf((*MockDatabase)(nil).GetByAlias), alias)
}

// GetLinksByUser mocks base method.
func (m *MockDatabase) GetLinksByUser(username string) ([]models.Link, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLinksByUser", username)
	ret0, _ := ret[0].([]models.Link)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLinksByUser indicates an expected call of GetLinksByUser.
func (mr *MockDatabaseMockRecorder) GetLinksByUser(username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLinksByUser", reflect.TypeOf((*MockDatabase)(nil).GetLinksByUser), username)
}

// GetUser mocks base method.
func (m *MockDatabase) GetUser(username string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", username)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockDatabaseMockRecorder) GetUser(username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockDatabase)(nil).GetUser), username)
}

// InsertNew mocks base method.
func (m_2 *MockDatabase) InsertNew(m models.ShortLink, username string) (string, error) {
	m_2.ctrl.T.Helper()
	ret := m_2.ctrl.Call(m_2, "InsertNew", m, username)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertNew indicates an expected call of InsertNew.
func (mr *MockDatabaseMockRecorder) InsertNew(m, username any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertNew", reflect.TypeOf((*MockDatabase)(nil).InsertNew), m, username)
}
