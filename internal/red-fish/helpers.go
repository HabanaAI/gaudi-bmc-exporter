package redfish

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

func convertAlertTime(alertTimeEpoch int64) string {
	return time.Unix(alertTimeEpoch, 0).UTC().String()
}

type RedFishErr struct {
	Err GeneralErr `json:"error"`
}
type GeneralErr struct {
	Code         string            `json:"code"`
	Message      string            `json:"message"`
	ExtendedInfo []ExtendedErrInfo `json:"@Message.ExtendedInfo"`
}
type ExtendedErrInfo struct {
	MessageID         string   `json:"MessageId"`
	Severity          string   `json:"Severity"`
	Resolution        string   `json:"Resolution"`
	Message           string   `json:"Message"`
	MessageArgs       []string `json:"MessageArgs"`
	RelatedProperties []string `json:"RelatedProperties"`
}

const (
	maxRedFishExtendedInfoEntries = 64
	maxRedFishErrorArrayEntries   = 128
	maxRedFishErrorFieldLength    = 4 << 10
	maxRedFishErrorMessageLength  = 64 << 10
)

const redFishErrorTruncated = "... (truncated)"

var ErrRedFishErrorTooLarge = errors.New("redfish error response exceeds limits")

type limitedErrorBuilder struct {
	builder   strings.Builder
	remaining int
	truncated bool
}

func newLimitedErrorBuilder(limit int) *limitedErrorBuilder {
	return &limitedErrorBuilder{remaining: limit - len(redFishErrorTruncated)}
}

func (b *limitedErrorBuilder) write(values ...string) {
	for _, value := range values {
		if len(value) > maxRedFishErrorFieldLength {
			value = value[:maxRedFishErrorFieldLength]
			b.truncated = true
		}
		if len(value) > b.remaining {
			value = value[:b.remaining]
			b.truncated = true
		}
		b.builder.WriteString(value)
		b.remaining -= len(value)
		if b.remaining == 0 {
			return
		}
	}
}

func (b *limitedErrorBuilder) String() string {
	if b.truncated {
		b.builder.WriteString(redFishErrorTruncated)
	}
	return b.builder.String()
}

func (r *RedFishErr) render() (string, bool) {
	builder := newLimitedErrorBuilder(maxRedFishErrorMessageLength)
	builder.write("code: ", r.Err.Code, ", message: ", r.Err.Message, ".")

	for extendedIndex, extended := range r.Err.ExtendedInfo {
		if extendedIndex == maxRedFishExtendedInfoEntries {
			builder.truncated = true
			break
		}
		builder.write(" MessageId: ", extended.MessageID, ", Severity: ", extended.Severity,
			", Resolution: ", extended.Resolution, ",Message: ", extended.Message)

		for messageArgIndex, messageArg := range extended.MessageArgs {
			if messageArgIndex == maxRedFishErrorArrayEntries {
				builder.truncated = true
				break
			}
			builder.write(" message arg: ", messageArg)
		}

		for relatedPropertyIndex, relatedProperty := range extended.RelatedProperties {
			if relatedPropertyIndex == maxRedFishErrorArrayEntries {
				builder.truncated = true
				break
			}
			builder.write(" related topic: ", relatedProperty)
		}
	}
	return builder.String(), builder.truncated
}

func (r *RedFishErr) Error() string {
	message, _ := r.render()
	return message
}

func (r *RedFishErr) validate() error {
	if len(r.Err.ExtendedInfo) > maxRedFishExtendedInfoEntries {
		return fmt.Errorf("%w: extended info has %d entries, maximum is %d",
			ErrRedFishErrorTooLarge, len(r.Err.ExtendedInfo), maxRedFishExtendedInfoEntries)
	}

	fields := []string{r.Err.Code, r.Err.Message}
	for _, extended := range r.Err.ExtendedInfo {
		if len(extended.MessageArgs) > maxRedFishErrorArrayEntries {
			return fmt.Errorf("%w: message arguments have %d entries, maximum is %d",
				ErrRedFishErrorTooLarge, len(extended.MessageArgs), maxRedFishErrorArrayEntries)
		}
		if len(extended.RelatedProperties) > maxRedFishErrorArrayEntries {
			return fmt.Errorf("%w: related properties have %d entries, maximum is %d",
				ErrRedFishErrorTooLarge, len(extended.RelatedProperties), maxRedFishErrorArrayEntries)
		}

		fields = append(fields, extended.MessageID, extended.Severity, extended.Resolution, extended.Message)
		fields = append(fields, extended.MessageArgs...)
		fields = append(fields, extended.RelatedProperties...)
	}

	for _, field := range fields {
		if len(field) > maxRedFishErrorFieldLength {
			return fmt.Errorf("%w: field has %d bytes, maximum is %d",
				ErrRedFishErrorTooLarge, len(field), maxRedFishErrorFieldLength)
		}
	}

	if _, truncated := r.render(); truncated {
		return fmt.Errorf("%w: rendered message exceeds %d bytes",
			ErrRedFishErrorTooLarge, maxRedFishErrorMessageLength)
	}

	return nil
}

// request creates and execute a request and unmarshal the response to the expected variable.
func (c *Client) request(ctx context.Context, log *logrus.Entry, url string, expected any, reqBody io.Reader) error {
	now := time.Now()

	defer func() {
		log.WithField("url", url).Infof("took %s to complete the request", time.Since(now).String())
	}()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, reqBody)
	if err != nil {
		return err
	}

	req.Header.Set("X-Auth-Token", c.session.Token)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := bmcmonitoring.ReadResponseBody(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var generalErr RedFishErr

		err = json.Unmarshal(body, &generalErr)
		if err != nil {
			return fmt.Errorf("failed to decode redfish error response: %w", err)
		}
		if err = generalErr.validate(); err != nil {
			return fmt.Errorf("failed to decode redfish error response: %w", err)
		}

		return &generalErr
	}

	err = json.Unmarshal(body, expected)
	if err != nil {
		return fmt.Errorf("failed to decode redfish response: %w", err)
	}

	return nil
}

func (c *Client) DeleteAllTokens(context.Context) error {
	return nil
}

func (c *Client) Name() string {
	return c.Hostname
}

// ValidToken checks if the token is valid by trying to get the session information
func (c *Client) ValidToken(ctx context.Context, log *logrus.Entry) (string, bool) {
	var res SessionInformation
	url := fmt.Sprintf("https://%s/redfish/v1/SessionService", c.Hostname)

	err := c.request(ctx, log, url, &res, nil)
	if err != nil {
		return err.Error(), false
	}
	return "", true
}

func (c *Client) RefreshToken(ctx context.Context) error {
	// if token is valid then we don't need to do anything
	if _, valid := c.ValidToken(ctx, c.log); valid {
		return nil
	}

	return c.authenticate(ctx)
}
