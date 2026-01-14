package redfish

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
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

func (r *RedFishErr) Error() string {
	errMsg := fmt.Sprintf("code: %s, message: %s.", r.Err.Code, r.Err.Message)

	for _, extended := range r.Err.ExtendedInfo {
		errMsg = fmt.Sprintf("%s MessageId: %s, Severity: %s, Resolution: %s,Message: %s",
			errMsg, extended.MessageID, extended.Severity, extended.Resolution, extended.Message)

		for _, msgArg := range extended.MessageArgs {
			errMsg = fmt.Sprintf("%s message arg: %s", errMsg, msgArg)
		}

		for _, relatedTopic := range extended.RelatedProperties {
			errMsg = fmt.Sprintf("%s related topic: %s", errMsg, relatedTopic)
		}
	}
	return errMsg
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		var generalErr RedFishErr

		err = json.Unmarshal(body, &generalErr)
		if err != nil {
			return fmt.Errorf("failed to decode error, %s: %w", string(body), err)
		}

		return &generalErr
	}

	err = json.Unmarshal(body, expected)
	if err != nil {
		return fmt.Errorf("%s :%w", string(body), err)
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
