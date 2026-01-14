// Package checkgrp maintains the group of handlers for health checking.
package checkgrp

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

// Handler manages the set of check endpoints.
type Handlers struct {
	Log   *logrus.Entry
	Build string
}

// will check if our application is ready.
func (h Handlers) Readiness(w http.ResponseWriter, r *http.Request) {

	status := "ok"
	statusCode := http.StatusOK
	data := struct {
		Build  string
		Status string `json:"status"`
	}{
		Build:  h.Build,
		Status: status,
	}
	if err := response(w, statusCode, data); err != nil {
		h.Log.WithFields(logrus.Fields{
			"ERROR": err,
		}).Error("readiness")
	}

}

// Liveness returns simple status info if the service is alive.
func (h Handlers) Liveness(w http.ResponseWriter, r *http.Request) {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := struct {
		Status string `json:"status,omit,empty"`
		Build  string `json:"build,omitempty"`
		Host   string `json:"host,omitempty"`
	}{
		Status: "up",
		Build:  h.Build,
		Host:   host,
	}
	statusCode := http.StatusOK
	if err := response(w, statusCode, data); err != nil {
		h.Log.WithFields(logrus.Fields{
			"ERROR": err,
		}).Error("liveness")
	}
}

// response will write the data to the ResponseWriter.
func response(w http.ResponseWriter, statusCode int, data interface{}) error {
	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}
	return nil
}
