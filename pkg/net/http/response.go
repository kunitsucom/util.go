package httpz

import (
	"bytes"
	"fmt"
	"net/http"
)

func NewBufferFromResponseBody(response *http.Response) (body *bytes.Buffer, err error) {
	buf := bytes.NewBuffer(nil)

	if _, err = buf.ReadFrom(response.Body); err != nil {
		return nil, fmt.Errorf("ReadFrom: %w", err)
	}

	return buf, nil
}
