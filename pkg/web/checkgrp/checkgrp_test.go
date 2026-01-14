// Package checkgrp maintains the group of handlers for health checking.
package checkgrp

import (
	"habana_bmc_exporter/pkg/logger"
	"net/http"
	"net/http/httptest"
	"testing"
)

type CgpTests struct {
	app http.Handler
}

// TestGrp is the entry point for testing check group handlers.
func TestGrp(t *testing.T) {

	cgh := Handlers{
		Build: "test",
		Log:   logger.New().WithField("service", "test"),
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/readiness", cgh.Readiness)
	mux.HandleFunc("/debug/liveness", cgh.Liveness)

	tests := CgpTests{
		app: mux,
	}
	t.Run("Readiness", tests.TestHandlers_Readiness)
	t.Run("Liveness", tests.TestHandlers_Liveness)
}

func (c *CgpTests) TestHandlers_Readiness(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/debug/readiness", nil)

	// run the request.
	c.app.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status ok, got: %s", http.StatusText(w.Code))
	}
}

func (c *CgpTests) TestHandlers_Liveness(t *testing.T) {

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/debug/liveness", nil)

	// run the request.
	c.app.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status ok, got: %s", http.StatusText(w.Code))
	}

}
