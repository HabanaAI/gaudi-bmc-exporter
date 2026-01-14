package rasmonitoringapi

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/sirupsen/logrus"
)

var passwordRegexp = regexp.MustCompile(`password=(.*)"`)

const (
	tokenHeader = "X-CSRFTOKEN"
)

func (c *Client) RefreshToken(ctx context.Context) error {

	// Check if it's a valid token before creating a new one.
	if _, valid := c.ValidToken(ctx, c.log); valid {
		return nil
	}

	// We logout to revoke the current token, but it may fail due  to reasons other than invalid token, so we retry the logout function.
	//  However, we ignore the error in case the token is already valid.
	err := retryWithAttempts(5, 5, c.Logout)
	if err != nil {
		c.log.WithError(err).Warning("failed to logout, but we will try to re-authenticate")
	}

	// if the token is not valid we don't need to logout, only to authenticate.
	err = retryWithAttempts(5, 5, c.authenticate)
	if err != nil {
		return fmt.Errorf("failed to re-authenticate: %w", err)
	}

	return nil
}

// ValidToken checks if the token is valid.
func (c *Client) ValidToken(ctx context.Context, log *logrus.Entry) (string, bool) {

	url := fmt.Sprintf("https://%s/api/settings/services", c.Hostname)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "failed creating request", false
	}
	req.Header.Add(tokenHeader, c.Auth.CSRFToken)
	req.AddCookie(c.Auth.Cookie)

	resp, err := c.Do(req)
	if err != nil {
		return "failed executing request", false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "failed reading the body", false
	}

	// check if we got an error in the body of the response
	// if there is an error then the token is not valid
	err = CheckIfError(body)
	if err != nil {
		return err.Error(), false
	}
	return "", true
}

func (c *Client) Logout() error {

	url := fmt.Sprintf("https://%s/api/session", c.Hostname)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return err
	}
	req.Header.Add(tokenHeader, c.Auth.CSRFToken)
	req.AddCookie(c.Auth.Cookie)

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
		return fmt.Errorf("status code %d, %s", resp.StatusCode, string(body))
	}

	// If the token is valid then the logout operation has failed.
	if _, valid := c.ValidToken(context.Background(), c.log); valid {
		return fmt.Errorf("the token is still valid, logout operation failed")
	}

	return nil
}

func (c *Client) authenticate() error {
	url := fmt.Sprintf("https://%s/api/session?username=%s&password=%s", c.Hostname, c.Username, c.Password)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	resp, err := c.Do(req)
	if err != nil {
		if strings.Contains(err.Error(), "password") {
			labeledErr := passwordRegexp.ReplaceAllString(err.Error(), "*******")
			return errors.New(labeledErr)
		}
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("exit code: %d, %s", resp.StatusCode, string(body))
	}

	var auth authDetails
	err = json.Unmarshal(body, &auth)
	if err != nil {
		return err
	}

	for _, cookie := range resp.Cookies() {
		if cookie.Name == "QSESSIONID" {
			auth.Cookie = cookie
		}
	}

	c.Auth = auth

	err = c.DeleteAllTokens(context.TODO())
	if err != nil {
		c.log.WithError(err).Warning()
	}

	return nil
}

func (c *Client) DeleteAllTokens(ctx context.Context) error {

	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	// Get the web sessions.
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("https://%s/api/settings/service-sessions?service_id=1", c.Hostname), nil)
	if err != nil {
		return err
	}

	req.Header.Add(tokenHeader, c.Auth.CSRFToken)
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(c.Auth.Cookie)

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
		return fmt.Errorf("status code - %d, %s", resp.StatusCode, string(body))
	}

	var sessionData SessionData

	err = json.Unmarshal(body, &sessionData)
	if err != nil {
		return err
	}

	// Filter all the sessions the we created except for the current one.
	var sessionsToRevoke SessionData

	for _, session := range sessionData {
		// add all the sessions that we created except the current one.
		if session.SessionID != c.Auth.RacsessionID {
			sessionsToRevoke = append(sessionsToRevoke, session)
		}
	}

	// Delete the sessions.
	for _, session := range sessionsToRevoke {
		err := c.DeleteSession(ctx, session.ID)
		if err != nil {
			c.log.WithError(err).Warningf("failed to delete session %+v", session)
		}

	}

	return nil

}

// DeleteSession will revoke the session credentials.
func (c *Client) DeleteSession(ctx context.Context, sessionId int) error {

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete,
		fmt.Sprintf("https://%s/api/settings/service-sessions/%d", c.Hostname, sessionId), nil)
	if err != nil {
		return err
	}

	req.Header.Add(tokenHeader, c.Auth.CSRFToken)
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(c.Auth.Cookie)

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
		return fmt.Errorf("status code - %d, %s", resp.StatusCode, string(body))
	}

	return nil
}

// Methods.
const (
	methodDirect   = "direct"
	methodIndirect = "indirect"
)

// doRequest does the request and retry according to the error.
func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	req.Header.Add(tokenHeader, c.Auth.CSRFToken)
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(c.Auth.Cookie)
	var resp *http.Response
	var body []byte

	err := RetryWithAttempts(5, 5, func() (bool, error) {
		c.Mutex.Lock()
		defer c.Mutex.Unlock()
		var err error
		resp, err = c.Do(req)
		if err != nil {
			return false, err
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return false, fmt.Errorf("failed to read response body: %w", err)
		}

		err = CheckIfError(body)

		if err != nil {
			if errors.Is(err, ErrDeviceNotResponding) ||
				errors.Is(err, ErrTimeout) || errors.Is(err, ErrCRC32) ||
				errors.Is(err, ErrCRC8) {
				return false, err
			}

			// TODO: check invalid auth error?
		}

		if resp.StatusCode != http.StatusOK {
			return false, fmt.Errorf("status code %d, %s", resp.StatusCode, string(body))
		}

		return false, err
	})

	return body, err
}

// opcodeData returns the opcode data.
func (c *Client) opcodeData(ctx context.Context, op Opcode, oam int, method string) (OpcodeData, error) {

	u, err := url.Parse(fmt.Sprintf("https://%s/ext/ras/%s/read", c.Hostname, method))
	if err != nil {
		return OpcodeData{}, err
	}

	q := u.Query()
	q.Add("oam", fmt.Sprintf("%d", oam))
	q.Add("offset", fmt.Sprintf("%d", op.Offset))
	q.Add("length", fmt.Sprintf("%d", op.Length))
	if method == methodIndirect {
		q.Add("opcode", fmt.Sprintf("%d", op.OpcodeNumber))
	}

	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return OpcodeData{}, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return OpcodeData{}, err
	}

	// fmt.Printf("opcode: %d, oam: %d, offset: %d, length: %d, output: %+v\n", op.OpcodeNumber, oam, op.Offset, op.Length, string(body))

	var data OpcodeData
	err = json.Unmarshal(body, &data)
	if err != nil {
		return OpcodeData{}, err
	}

	return data, err
}

// decodeOpcode will decode the opcode according to the expected field of the opcode.
func (c *Client) decodeOpcode(ctx context.Context, op Opcode, oam int, fieldName, method string) (string, error) {

	// Add retry if CRC error
	var data OpcodeData

	data, err := c.opcodeData(ctx, op, oam, method)
	if err != nil {
		return "", err
	}

	if len(data.Data) == 0 {
		return "", fmt.Errorf("failed to get info of field %s, oam %d, len of data is 0", fieldName, oam)
	}

	val, err := convertData(data, op.ExpectedType)
	if err != nil {
		return "", err

	}

	// fmt.Printf("field: %s\n", fieldName)
	// fmt.Printf("data: %+v\n", data.Data)
	// fmt.Printf("val: %+v\n", val)

	return val, nil
}

type directRequest struct {
	Oam    string `json:"oam"`
	Opcode string `json:"opcode"`
	Cpsp   string `json:"cpsp"`
}

type directResponse struct {
	Cpsr int `json:"cpsr"`
}

// directCommand will execute direct commands
func (c *Client) directCommand(ctx context.Context, path string, oam, opcode int) (int, error) {

	url := fmt.Sprintf("https://%s/ext/ras/direct/%s", c.Hostname, path)
	reqData := directRequest{
		Oam:    fmt.Sprintf("%d", oam),
		Opcode: fmt.Sprintf("%d", opcode),
		Cpsp:   "0",
	}

	b, err := json.Marshal(reqData)
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(b))
	if err != nil {
		return 0, err
	}

	body, err := c.doRequest(req)
	if err != nil {
		return 0, err
	}

	var respBody directResponse
	err = json.Unmarshal(body, &respBody)
	if err != nil {
		return 0, err
	}
	return respBody.Cpsr, nil
}

// convertData will convert the data into the expected data type.
func convertData(data OpcodeData, expectedType ExpectedType) (string, error) {

	switch expectedType {
	case Int:
		data.Data = reverseSlice(data.Data)
		s := hex.EncodeToString(data.Data)
		decimal_num, err := strconv.ParseInt(s, 16, 64)
		if err != nil {
			return "", err
		}
		return fmt.Sprintf("%d", decimal_num), nil
	case ReverseString:
		data.Data = reverseSlice(data.Data)
		s := hex.EncodeToString(data.Data)
		return s, nil

	case String:
		s := hex.EncodeToString(data.Data)

		return s, nil
	case BinArray:
		// only address the first 24 ports.
		// TODO: this is only for gaudi2,need to be more generic.
		data.Data = data.Data[:3]

		data.Data = reverseSlice(data.Data)

		var s string

		for _, n := range data.Data {
			binNum := fmt.Sprintf("%.8b", n)
			s = fmt.Sprintf("%s%s", s, binNum)
		}

		return reverseString(s), nil

	case Bin:
		var s string
		for _, n := range data.Data {
			binNum := fmt.Sprintf("%.8b", n)
			s = fmt.Sprintf("%s%s", s, binNum)
		}

		return s, nil

	case AsciiString:
		var s string
		for _, n := range data.Data {

			// Only add letters/numbers.
			if unicode.IsLetter(rune(n)) || unicode.IsNumber(rune(n)) {
				s = fmt.Sprintf("%s%s", s, string(n))
			}
		}

		return s, nil

	default:
		return "", fmt.Errorf("unsupported type %s", expectedType)
	}
}

func reverseSlice[T any](slice []T) []T {
	reversed := make([]T, len(slice))
	for i, j := 0, len(slice)-1; i < len(slice); i, j = i+1, j-1 {
		reversed[i] = slice[j]
	}
	return reversed
}

func (c *Client) Name() string {
	return c.Hostname
}

func retryWithAttempts(attempts int, interval int, f func() error) error {
	var err error
	dur := time.Duration(interval) * time.Second
	for i := 0; i < attempts; i++ {
		if err = f(); err == nil {
			// IF NO error
			return nil
		}
		time.Sleep(dur)
	}
	return fmt.Errorf("failed after %d attempts: %v", attempts, err)
}

type writeRequest struct {
	Oam    string `json:"oam"`
	Opcode string `json:"opcode"`
	Offset string `json:"offset"`
	Length string `json:"length"`
	Data   []int  `json:"data"`
}

func (c *Client) write(ctx context.Context, oam int, op Opcode, index int) error {
	c.Mutex.Lock()
	defer c.Mutex.Unlock()
	url := fmt.Sprintf("https://%s/ext/ras/indirect/write", c.Hostname)

	body := writeRequest{
		Oam:    fmt.Sprintf("%d", oam),
		Opcode: fmt.Sprintf("%d", op.OpcodeNumber),
		Offset: fmt.Sprintf("%d", op.Offset),
		Length: fmt.Sprintf("%d", op.Length),
		Data:   []int{index},
	}

	// make sure the request will be the size of the length.
	bLen := len(body.Data)
	for i := 0; i < op.Length-bLen; i++ {
		body.Data = append(body.Data, 0)
	}

	b, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, bytes.NewReader(b))
	if err != nil {
		return err
	}

	req.Header.Add(tokenHeader, c.Auth.CSRFToken)
	req.Header.Add("Content-Type", "application/json")
	req.AddCookie(c.Auth.Cookie)

	resp, err := c.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(req.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to write, status code: %d, %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func reverseString(str string) string {
	byte_str := []rune(str)
	for i, j := 0, len(byte_str)-1; i < j; i, j = i+1, j-1 {
		byte_str[i], byte_str[j] = byte_str[j], byte_str[i]
	}
	return string(byte_str)
}

func RetryWithAttempts(attempts int, interval int, f func() (bool, error)) error {
	var err error
	dur := time.Duration(interval) * time.Second
	for i := 0; i < attempts; i++ {

		stop, err := f()

		// Stop and return the error
		if stop {
			return err
		}

		// IF NO error
		if err == nil {
			return nil
		}

		time.Sleep(dur)
	}
	return fmt.Errorf("failed after %d attempts: %w", attempts, err)
}
