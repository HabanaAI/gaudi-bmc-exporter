package rasmonitoringapi

import "net/http"

type transport struct {
	RoundTripFunc func(r *http.Request) (*http.Response, error)
}

func (t *transport) RoundTrip(r *http.Request) (*http.Response, error) {
	return t.RoundTripFunc(r)
}
