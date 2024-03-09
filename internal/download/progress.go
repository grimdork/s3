package download

import (
	"fmt"
	"io"
)

// Progress reporter.
type Progress struct {
	read    int
	written int
	total   int
	reader  io.ReadCloser
	writer  io.WriteCloser
}

// NewProgress creates a new progress reporter.
func NewProgress(r io.ReadCloser, w io.WriteCloser, t int) *Progress {
	return &Progress{
		reader: r,
		writer: w,
		total:  t,
	}
}

// report progress.
func (p *Progress) report() {
	if p.total == -1 {
		fmt.Printf("\033[2K\rRead %d bytes, wrote %d bytes (unknown total)", p.read, p.written)
	} else {
		fmt.Printf("\033[2K\rRead %d/%d bytes, wrote %d/%d bytes", p.read, p.total, p.written, p.total)
	}
}

// Copy data from reader to writer.
func (p *Progress) Copy() error {
	buf := make([]byte, 1024)
	for {
		n, err := p.reader.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
		p.read += n
		p.written += n
		p.report()
		if _, err := p.writer.Write(buf[:n]); err != nil {
			return err
		}
	}

	return nil
}
