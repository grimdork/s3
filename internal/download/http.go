package download

import (
	"io"
	"net/http"
)

// OpenWebReader reads chunks of unsecured HTTP data and writes it to a destination.
func OpenWebReader(url string) (io.ReadCloser, int, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, 0, err
	}

	return res.Body, int(res.ContentLength), nil
}
