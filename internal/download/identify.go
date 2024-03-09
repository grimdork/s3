package download

import (
	"errors"
	"io"
	"os"
	"strings"
)

const (
	// URIUnknown is an error or not implemented scheme.
	URIUnknown = iota
	// URILocal is a lcoal file.
	URILocal
	// URIWeb is an unsecured web address.
	URIWeb
	// URIS3 is a Simple Storage Service address.
	URIS3
	// URISSH is a secure shell address.
	URISSH
)

// ErrUnknownURI is returned when an unknown source or destination URI/scheme is provided.
var (
	ErrUnknownURI   = errors.New("unknown URI")
	ErrCantWriteWeb = errors.New("can't write to web addresses")
	ErrNotAFile     = errors.New("not a file")
)

// IdentifyURI to know which read/write function to use.
func IdentifyURI(url string) uint8 {
	url = strings.ToLower(url)
	url = strings.TrimSpace(url)
	if url == "" {
		return URIUnknown
	}

	if strings.HasPrefix(url, "http://") {
		return URIWeb
	}

	if strings.HasPrefix(url, "https://") {
		return URIWeb
	}

	if strings.HasPrefix(url, "s3://") {
		return URIS3
	}

	if strings.HasPrefix(url, "ssh://") {
		return URISSH
	}

	if strings.Contains(url, "://") {
		return URIUnknown
	}

	return URILocal
}

// OpenSource for reading.
func OpenSource(url string) (io.ReadCloser, int, error) {
	switch IdentifyURI(url) {
	case URIWeb:
		return OpenWebReader(url)
	case URIS3:
	case URISSH:
	case URILocal:
		f, err := os.Open(url)
		if err != nil {
			return nil, 0, err
		}

		st, err := f.Stat()
		if err != nil {
			return nil, 0, err
		}

		if st.IsDir() {
			return nil, 0, ErrNotAFile
		}

		return f, int(st.Size()), nil
	}

	return nil, 0, ErrUnknownURI
}

// OpenDestination for writing.
func OpenDestination(url string) (io.WriteCloser, error) {
	switch IdentifyURI(url) {
	case URIWeb:
		return nil, ErrCantWriteWeb
	case URIS3:
	case URISSH:
	case URILocal:
		return os.Create(url)
	}

	return nil, ErrUnknownURI
}
