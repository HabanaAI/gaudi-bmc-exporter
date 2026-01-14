package web

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/sirupsen/logrus"
)

type Middleware func(Handler, *logrus.Entry) Handler

// wrapMiddleware creates a new handler by wrapping middleware around a final
// handler. The middleware Handlers will be executed by requests in the order
// they are provided.
func wrapMiddleware(mw []Middleware, handler Handler, log *logrus.Entry) Handler {

	// Loop backwards through the middleware invoking each one. Replace the
	// Handler with the new wrapped handler. Lopping backwards ensures that the
	// first middleware of the slice is the first to be executed by requests.
	for i := len(mw) - 1; i >= 0; i-- {
		h := mw[i]
		if h != nil {
			handler = h(handler, log)
		}
	}
	return handler
}

// TimeSinceMid is middleware that register the time it took to complete the call to the url.
func TimeSinceMid(handler Handler, log *logrus.Entry) Handler {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		log.Infof("got a call to handler %s", r.URL)
		handler(ctx, w, r) // original function call
		log.Infof("took %+v to complete %s", time.Since(now).String(), r.URL)
	}

}

// PanicMid is a middleware that will deal with panic.
func PanicMid(handler Handler, log *logrus.Entry) Handler {

	return func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {

				trace := debug.Stack()
				// Stack trace will we provided.
				log.WithError(fmt.Errorf("PANIC [%v] TRACE [%s]", rec, string(trace))).Error()
			}
		}()

		handler(ctx, w, r) // original function call
	}

}
