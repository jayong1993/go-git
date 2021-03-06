package ioutil

import (
	"bufio"
	"errors"
	"io"
)

type readPeeker interface {
	io.Reader
	Peek(int) ([]byte, error)
}

var (
	ErrEmptyReader = errors.New("reader is empty")
)

// NonEmptyReader takes a reader and returns it if it is not empty, or
// `ErrEmptyReader` if it is empty. If there is an error when reading the first
// byte of the given reader, it will be propagated.
func NonEmptyReader(r io.Reader) (io.Reader, error) {
	pr, ok := r.(readPeeker)
	if !ok {
		pr = bufio.NewReader(r)
	}

	_, err := pr.Peek(1)
	if err == io.EOF {
		return nil, ErrEmptyReader
	}

	if err != nil {
		return nil, err
	}

	return pr, nil
}

type readCloser struct {
	io.Reader
	closer io.Closer
}

func (r *readCloser) Close() error {
	return r.closer.Close()
}

// NewReadCloser creates an `io.ReadCloser` with the given `io.Reader` and
// `io.Closer`.
func NewReadCloser(r io.Reader, c io.Closer) io.ReadCloser {
	return &readCloser{Reader: r, closer: c}
}
