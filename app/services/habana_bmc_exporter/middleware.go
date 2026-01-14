package main

import (
	"context"
	"habana_bmc_exporter/pkg/web"
	"net/http"

	"github.com/sirupsen/logrus"
)

// combineContextMiddleware combines the main context of the application into the context of the
// web request, so any cancellation signal in the application level, will be propagated into
// the request and cancel all executing requested.
func (app application) combineContextMiddleware(handler web.Handler, log *logrus.Entry) web.Handler {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		r = r.WithContext(ctx)
		handler(ctx, w, r)

	}
}
