package testutils

import (
	"net/http"
	"net/http/httptest"
)

func MakeTestHTTPServer(handler http.Handler) (*httptest.Server, *http.Client) {
	server := httptest.NewServer(handler)

	return server, server.Client()
}
