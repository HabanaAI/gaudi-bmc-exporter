package redfish

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type SessionData struct {
	SessionID string
	Token     string
}
type Client struct {
	log *logrus.Entry
	*bmcmonitoring.Client
	session SessionData
}

const (
	defaultNumCards = 8
	defaultNumPorts = 24
	defaultNumLanes = 16
)

func NewClient(opts bmcmonitoring.ClientOpts, log *logrus.Entry) (*Client, error) {
	if opts.CreationTimeout == 0 {
		opts.CreationTimeout = 10 * time.Second
	}

	if opts.Oams == 0 {
		opts.Oams = defaultNumCards
	}
	if opts.Ports == 0 {
		opts.Ports = defaultNumPorts
	}

	if opts.Lanes == 0 {
		opts.Lanes = defaultNumLanes
	}

	ctx, cancel := context.WithTimeout(context.Background(), opts.CreationTimeout)
	defer cancel()

	type res struct {
		client *Client
		err    error
	}

	ch := make(chan res)

	go func() {
		c, err := bmcmonitoring.NewClient(opts)
		if err != nil {
			ch <- res{err: err}
			return
		}

		c.Timeout = time.Minute * 3

		ll := log.WithField("hostname", opts.Hostname)

		client := &Client{
			Client: c,
			log:    ll,
		}
		err = client.authenticate(ctx)
		if err != nil {
			ch <- res{err: fmt.Errorf("failed to authenticate: %w", err)}
			return
		}

		ch <- res{client: client}
	}()

	select {
	case <-ctx.Done():
		return nil, fmt.Errorf("timeout exceeded for creating client %s", opts.Hostname)
	case resp := <-ch:
		if resp.err != nil {
			return nil, resp.err
		}
		return resp.client, nil
	}
}

func (c *Client) authenticate(ctx context.Context) error {
	type creds struct {
		Username string `json:"UserName"`
		Password string `json:"Password"`
	}

	reqBody, err := json.Marshal(creds{
		Username: c.Username,
		Password: c.Password,
	})

	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, fmt.Sprintf("https://%s/redfish/v1/SessionService/Sessions", c.Hostname), bytes.NewReader(reqBody))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("exit code: %d, %s", resp.StatusCode, string(body))
	}

	tokens := resp.Header["X-Auth-Token"]

	var auth struct {
		ID string `json:"Id"`
	}

	err = json.Unmarshal(body, &auth)
	if err != nil {
		return err
	}

	if len(tokens) != 1 {
		return fmt.Errorf("failed to get token")
	}

	c.session = SessionData{
		Token:     tokens[0],
		SessionID: auth.ID,
	}
	return nil
}

func (c *Client) Logout() error {
	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf("https://%s/redfish/v1/SessionService/Sessions/%s", c.Hostname, c.session.SessionID), nil)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
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

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("exit code: %d, %s", resp.StatusCode, string(body))
	}
	return nil

}
