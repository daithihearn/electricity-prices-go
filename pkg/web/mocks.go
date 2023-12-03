package web

import (
	"bytes"
	"io"
	"net/http"
)

type MockHTTPClient struct {
	MockResp *http.Response
	MockErr  error
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	return m.MockResp, m.MockErr
}

type MockReadCloser struct {
	io.Reader
}

func (m MockReadCloser) Close() error {
	return nil // or add logic to handle the closure of the Reader
}

func NewMockReadCloser(response string) io.ReadCloser {
	return MockReadCloser{Reader: bytes.NewBufferString(response)}
}
