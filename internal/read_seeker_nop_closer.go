package internal

import (
	"bytes"
	"io"
	"sync"
)

func BuildReadSeekCloser(body io.Reader) (io.ReadSeekCloser, error) {
	if body == nil {
		return nil, nil
	}

	// Do not check for io.ReadSeekCloser as we don't know if Close is a no-op or not.
	// Instead, we check for io.ReadSeeker and wrap it with a no-op Close method.
	if seeker, ok := body.(io.ReadSeeker); ok {
		return newBytesNopCloser(seeker), nil
	}

	data, err := io.ReadAll(body)
	if err != nil {
		return nil, err
	}

	// return io.ReadSeekCloser so other wrappers can also reuse it if needed.
	return newBytesNopCloser(bytes.NewReader(data)), nil
}

type readSeekerNopCloser struct {
	d io.ReadSeeker
	m sync.Mutex
}

func (b *readSeekerNopCloser) Close() error { return nil }

func (b *readSeekerNopCloser) Read(p []byte) (n int, err error) {
	b.m.Lock()
	defer b.m.Unlock()

	return b.d.Read(p)
}

func (b *readSeekerNopCloser) Seek(offset int64, whence int) (int64, error) {
	b.m.Lock()
	defer b.m.Unlock()

	return b.d.Seek(offset, whence)
}

func newBytesNopCloser(data io.ReadSeeker) *readSeekerNopCloser {
	return &readSeekerNopCloser{d: data}
}
