package bmcmonitoring

import (
	"errors"
	"fmt"
	"io"
)

const maxResponseBodySize int64 = 4 << 20

var ErrResponseBodyTooLarge = errors.New("response body exceeds size limit")

func ReadResponseBody(body io.Reader) ([]byte, error) {
	data, err := io.ReadAll(io.LimitReader(body, maxResponseBodySize+1))
	if err != nil {
		return nil, err
	}
	if int64(len(data)) > maxResponseBodySize {
		return nil, fmt.Errorf("%w: maximum is %d bytes", ErrResponseBodyTooLarge, maxResponseBodySize)
	}
	return data, nil
}
