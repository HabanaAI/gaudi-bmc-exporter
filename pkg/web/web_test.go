package web

import (
	"context"
	"habana_bmc_exporter/pkg/logger"
	"habana_bmc_exporter/pkg/web/checkgrp"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_debugStandardLibraryMux(t *testing.T) {

	mux := debugStandardLibraryMux()

	routes := []string{"/debug/pprof/", "/debug/pprof/cmdline", "/debug/pprof/symbol", "/debug/pprof/trace"}
	for _, route := range routes {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, route, nil)
		mux.ServeHTTP(w, r)

		// run the request.
		if w.Code != http.StatusOK {
			t.Fatalf("expected status ok, got: %s", http.StatusText(w.Code))
		}

	}
}

func TestServer(t *testing.T) {

	ctx := context.Background()
	mux := NewServer(ctx,
		ServerOpts{
			Port: "4000",
		}, checkgrp.Handlers{
			Build: "test",
			Log:   logger.New().WithField("service", "test"),
		},
		[]WebHandler{
			{Route: "/test", Handler: func(ctx context.Context, w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("ok"))
				w.WriteHeader(http.StatusOK)
			}}})
	defer mux.Close()

	t.Run("test adding custom handlers", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/test", nil)
		mux.Handler.ServeHTTP(w, r)
		// run the request.
		if w.Code != http.StatusOK {
			t.Fatalf("expected status ok, got: %s", http.StatusText(w.Code))
		}
	})
}
