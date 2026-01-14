package redfish

import (
	"net/http"
	"testing"
)

type transport struct {
	RoundTripFunc func(r *http.Request) (*http.Response, error)
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	return t.RoundTripFunc(r)
}

// checkAuth verifies that the auth is provided with the request
func checkAuth(t *testing.T, r *http.Request) {

}
