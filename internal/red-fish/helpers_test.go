package redfish

import (
	"bytes"
	"context"
	"encoding/json"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
)

func TestRedFishErrError(t *testing.T) {
	err := &RedFishErr{Err: GeneralErr{
		Code:    "Base.1.0.GeneralError",
		Message: "A general error has occurred.",
		ExtendedInfo: []ExtendedErrInfo{{
			MessageID:         "Base.1.0.PropertyValueError",
			Severity:          "Warning",
			Resolution:        "Correct the value and resubmit.",
			Message:           "The value is invalid.",
			MessageArgs:       []string{"value", "property"},
			RelatedProperties: []string{"#/Property"},
		}},
	}}

	require.Equal(t, "code: Base.1.0.GeneralError, message: A general error has occurred.."+
		" MessageId: Base.1.0.PropertyValueError, Severity: Warning,"+
		" Resolution: Correct the value and resubmit.,Message: The value is invalid."+
		" message arg: value message arg: property related topic: #/Property", err.Error())
	require.NoError(t, err.validate())
}

func TestRedFishErrErrorBoundsUntrustedArrays(t *testing.T) {
	err := &RedFishErr{Err: GeneralErr{
		ExtendedInfo: []ExtendedErrInfo{{
			MessageArgs: make([]string, 8192),
		}},
	}}
	for index := range err.Err.ExtendedInfo[0].MessageArgs {
		err.Err.ExtendedInfo[0].MessageArgs[index] = "x"
	}

	message := err.Error()

	require.LessOrEqual(t, len(message), maxRedFishErrorMessageLength)
	require.Equal(t, maxRedFishErrorArrayEntries, strings.Count(message, " message arg: "))
	require.True(t, strings.HasSuffix(message, redFishErrorTruncated))
	require.ErrorIs(t, err.validate(), ErrRedFishErrorTooLarge)
}

func TestRequestRejectsRedFishErrorWithTooManyArrayEntries(t *testing.T) {
	responseBody, err := json.Marshal(RedFishErr{Err: GeneralErr{
		ExtendedInfo: []ExtendedErrInfo{{
			MessageArgs: make([]string, 8192),
		}},
	}})
	require.NoError(t, err)

	client := &Client{Client: &bmcmonitoring.Client{Client: http.Client{
		Transport: &transport{RoundTripFunc: func(*http.Request) (*http.Response, error) {
			return &http.Response{
				StatusCode: http.StatusInternalServerError,
				Body:       io.NopCloser(bytes.NewReader(responseBody)),
			}, nil
		}},
	}}}

	err = client.request(context.Background(), logrus.New().WithField("test", t.Name()),
		"https://bmc.example/redfish/v1/test", &struct{}{}, nil)

	require.ErrorIs(t, err, ErrRedFishErrorTooLarge)
}
