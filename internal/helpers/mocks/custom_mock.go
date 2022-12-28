package mocks

import "net/http"

type MockHelpers struct {
}

func NewMockHelpers() *MockHelpers {
	return &MockHelpers{}
}

/*CheckURL(longURL string) bool
GetIP(r *http.Request) string
RandString() string*/

func (m *MockHelpers) CheckURL(longURL string) bool {
	return false
}

func (m *MockHelpers) GetIP(r *http.Request) string {
	return "testIpInfo"
}

func (m *MockHelpers) RandString() string {
	return "testString"
}