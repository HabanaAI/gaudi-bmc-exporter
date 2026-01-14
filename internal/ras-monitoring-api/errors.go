package rasmonitoringapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

var (
	// Ras errors

	ErrUninitialized           = fmt.Errorf("not initialized")
	ErrInvalidArg              = fmt.Errorf("invalid argument")
	ErrNotSupported            = fmt.Errorf("not supported")
	ErrResourceUnavailable     = fmt.Errorf("resource not available")
	ErrInsufficientSize        = fmt.Errorf("insufficient size")
	ErrI2cDev                  = fmt.Errorf("i2c dev")
	ErrDeviceNotResponding     = fmt.Errorf("device not responding")
	ErrNoResponseOfDeviceOnBus = fmt.Errorf("no response of device on bus")
	ErrTimeout                 = fmt.Errorf("timeout")
	ErrCRC8                    = fmt.Errorf("crc 8")
	ErrCRC32                   = fmt.Errorf("crc 32")
	ErrMemory                  = fmt.Errorf("memory")
	ErrNoData                  = fmt.Errorf("no data")
	ErrUnknown                 = fmt.Errorf("unknown")

	ErrInvalidAuthentication = fmt.Errorf("invalid authentication")
)

type RasError struct {
	Err  string `json:"error"`
	Code int    `json:"code"`
}

type AuthError struct {
	Error string `json:"error"`
	CC    int    `json:"cc"`
}

func CheckIfError(body []byte) error {

	d := json.NewDecoder(bytes.NewReader(body))

	d.DisallowUnknownFields()

	// opErr represents the way we get errors from the api.
	var opErr RasError
	err := d.Decode(&opErr)

	if err == nil {

		switch {
		case strings.Contains(opErr.Err, "ret = 1"):
			return ErrUninitialized
		case strings.Contains(opErr.Err, "ret = 2"):
			return ErrInvalidArg
		case strings.Contains(opErr.Err, "ret = 3"):
			return ErrNotSupported
		case strings.Contains(opErr.Err, "ret = 4"):
			return ErrResourceUnavailable
		case strings.Contains(opErr.Err, "ret = 5"):
			return ErrInsufficientSize
		case strings.Contains(opErr.Err, "ret = 6"):
			return ErrI2cDev
		case strings.Contains(opErr.Err, "ret = 7"):
			return ErrDeviceNotResponding
		case strings.Contains(opErr.Err, "ret = 8"):
			return ErrTimeout
		case strings.Contains(opErr.Err, "ret = 9"):
			return ErrCRC8
		case strings.Contains(opErr.Err, "ret = 10"):
			return ErrCRC32
		case strings.Contains(opErr.Err, "ret = 11"):
			return ErrMemory
		case strings.Contains(opErr.Err, "ret = 12"):
			return ErrNoData
		case strings.Contains(opErr.Err, "ret = 13"):
			return ErrUnknown
		default:
			return fmt.Errorf("error: %s, code: %d", opErr.Err, opErr.Code)

		}

	}

	d = json.NewDecoder(bytes.NewReader(body))
	d.DisallowUnknownFields()

	var authErr AuthError

	err = d.Decode(&authErr)

	// if the Unmarshal don't succeed it's not an error.
	if err == nil {

		if strings.Contains(authErr.Error, "Invalid Authentication") {
			return ErrInvalidAuthentication
		}

		return fmt.Errorf("error: %s, cc: %d", authErr.Error, authErr.CC)
	}

	return nil
}
