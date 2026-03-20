package internal

import (
	"bytes"
	"io"
)

func BuildReadSeekCloser(body io.Reader) (io.ReadSeekCloser, error) {
	if body == nil {
		return nil, nil
	}

	// Do not check for io.ReadSeekCloser as we don't know if Close is a no-op or not.
	// Instead, we check for io.ReadSeeker and wrap it with a no-op Close method.
	if seeker, ok := body.(io.ReadSeeker); ok {
		return bytesNopCloser{ReadSeeker: seeker}, nil
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}
	// return io.ReadSeekCloser so other wrappers can also reuse it if needed.
	return bytesNopCloser{ReadSeeker: bytes.NewReader(data)}, nil
}

type bytesNopCloser struct {
	io.ReadSeeker
}

func (b bytesNopCloser) Close() error { return nil }
