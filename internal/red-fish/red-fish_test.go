package redfish

import (
	"compress/gzip"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestNewClientRejectsOversizedGzipAuthenticationResponse(t *testing.T) {
	server := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if encoding := r.Header.Get("Accept-Encoding"); encoding != "gzip" {
			t.Errorf("Accept-Encoding = %q, want gzip", encoding)
		}
		w.Header().Set("Content-Encoding", "gzip")
		w.WriteHeader(http.StatusUnauthorized)

		writer := gzip.NewWriter(w)
		_, _ = io.CopyN(writer, strings.NewReader(strings.Repeat("x", 4<<20+1)), 4<<20+1)
		_ = writer.Close()
	}))
	defer server.Close()

	_, err := NewClient(bmcmonitoring.ClientOpts{
		Hostname: strings.TrimPrefix(server.URL, "https://"),
	}, logrus.New().WithField("test", t.Name()))

	require.ErrorIs(t, err, bmcmonitoring.ErrResponseBodyTooLarge)
}
