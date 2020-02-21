package httpclient

import (
	"net/http"

	"github.com/stretchr/testify/mock"
)

// MockPlugin provides a mock plugin for heimdall
type MockPlugin struct {
	mock.Mock
}

// OnRequestStart is called when the request starts
func (m *MockPlugin) OnRequestStart(req *http.Request) {
	m.Called(req)
}

// OnRequestEnd is called when the request ends
func (m *MockPlugin) OnRequestEnd(req *http.Request, res *http.Response) {
	m.Called(req, res)
}

// OnError is called when the request errors out
func (m *MockPlugin) OnError(req *http.Request, err error) {
	m.Called(req, err)
}
