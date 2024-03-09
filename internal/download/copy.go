package download

import "io"

// Copy from source ReadCloser to dest WriteCloser and close them afterwards,
// showing progress during the process. Source and destination will both be
// closed no matter how this returns.
func Copy(src io.ReadCloser, dest io.WriteCloser, size int) error {
	defer src.Close()
	defer dest.Close()

	p := NewProgress(src, dest, size)
	return p.Copy()
}
