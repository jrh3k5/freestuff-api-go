package http

import "net/http"

// Doer is an abstraction of the Go-provided net/http/Client struct to facilitate injection.
type Doer interface {
	// Do executes the given request.
	Do(request *http.Request) (*http.Response, error)
}
