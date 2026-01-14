package g3redfish

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"habana_bmc_exporter/internal/g3redfish/config"
	bmcmonitoring "habana_bmc_exporter/pkg/bmc-monitoring"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type SessionData struct {
	SessionID string
	Token     string
}
type Client struct {
	log *logrus.Entry
	*bmcmonitoring.Client
	session SessionData
	config  config.Config
}

const (
	defaultNumCards = 8
	defaultNumPorts = 24
	defaultNumLanes = 16
	apiVersion      = "v1"
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
		err = client.loadConfig()
		if err != nil {
			ch <- res{err: fmt.Errorf("failed to load config: %w", err)}
			return
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

func (c *Client) DeleteAllTokens(context.Context) error {
	return nil
}

func (c *Client) Logout() error {
	requestUrl := c.createUrl(fmt.Sprintf("SessionService/Sessions/%s", c.session.SessionID))
	req, err := http.NewRequest(http.MethodDelete, requestUrl, nil)
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

func (c *Client) Name() string {
	return c.Hostname
}

// ValidToken checks if the token is valid by trying to get the session information
func (c *Client) ValidToken(ctx context.Context, log *logrus.Entry) (string, bool) {
	var res V1SessionService
	requestUrl := c.createUrl("SessionService")
	err := c.request(ctx, log, requestUrl, &res, nil)
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

	requestUrl := c.createUrl("SessionService/Sessions")
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, requestUrl, bytes.NewReader(reqBody))
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

func (c *Client) loadConfig() error {
	yamlFile, err := os.ReadFile("/etc/g3redfish.yaml")
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &c.config)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) createUrl(apiPath string) string {
	// check if the apiPath already contains the apiVersion
	if apiPath[0] == '/' {
		return fmt.Sprintf("https://%s%s", c.Hostname, apiPath)
	}
	return fmt.Sprintf("https://%s/redfish/%s/%s", c.Hostname, apiVersion, apiPath)
}

// request creates and execute a request and unmarshal the response to the expected variable.
func (c *Client) request(ctx context.Context, log *logrus.Entry, url string, obj any, reqBody io.Reader) error {
	now := time.Now()

	defer func() {
		log.WithField("url", url).Debugf("took %s to complete the request", time.Since(now).String())
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

	err = json.Unmarshal(body, obj)
	if err != nil {
		return fmt.Errorf("%s :%w", string(body), err)
	}

	return nil
}
