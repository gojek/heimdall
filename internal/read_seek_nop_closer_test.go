package internal

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildReadSeekCloser_WithNilBody_ReturnsNilReadSeekCloser(t *testing.T) {
	result, err := BuildReadSeekCloser(nil)
	assert.Nil(t, result)
	assert.NoError(t, err)
}

func TestBuildReadSeekCloser_WithReadSeeker_WrapsWithNopCloser(t *testing.T) {
	data := []byte("test data")
	seeker := &readSeekCounter{ReadSeeker: bytes.NewReader(data)}

	result, err := BuildReadSeekCloser(seeker)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Zero(t, seeker.readCount)
	err = result.Close()
	assert.NoError(t, err)

	// can read after no/op close
	buf := make([]byte, len(data))
	n, err := result.Read(buf)
	assert.NoError(t, err)
	assert.Equal(t, len(data), n)
	assert.Equal(t, data, buf)
	assert.Equal(t, 1, seeker.readCount)
	err = result.Close() // close again to ensure no-op
	assert.NoError(t, err)

}

func TestBuildReadSeekCloser_WithPlainReader_ReadsAllDataAndWrapsWithNopCloser(t *testing.T) {
	data := []byte("test data content")
	reader := &readSeekCounter{ReadSeeker: bytes.NewReader(data)}
	result, err := BuildReadSeekCloser(io.NopCloser(reader))
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 2, reader.readCount) // should read all data once

	buf := make([]byte, len(data))
	n, err := result.Read(buf)
	assert.True(t, err == nil || err == io.EOF)
	assert.Equal(t, len(data), n)
	assert.Equal(t, data, buf)
	assert.Equal(t, 2, reader.readCount) // no further reads
}

func TestBuildReadSeekCloser_WithPlainReader_SupportsSeek(t *testing.T) {
	data := []byte("seek test data")
	reader := bytes.NewReader(data)

	result, _ := BuildReadSeekCloser(reader)

	result.Seek(5, io.SeekStart)
	buf := make([]byte, 4)
	result.Read(buf)
	assert.Equal(t, []byte("test"), buf)
}

func TestBuildReadSeekCloser_WithEmptyReader_ReturnsValidReadSeekCloser(t *testing.T) {
	reader := bytes.NewReader([]byte{})

	result, err := BuildReadSeekCloser(reader)
	assert.NoError(t, err)
	assert.NotNil(t, result)

	buf := make([]byte, 1)
	n, err := result.Read(buf)
	assert.Equal(t, 0, n)
	assert.Equal(t, io.EOF, err)
}

func TestBuildReadSeekCloser_WithReadSeeker_CanSeekAndRead(t *testing.T) {
	data := []byte("position test")
	seeker := bytes.NewReader(data)

	result, _ := BuildReadSeekCloser(seeker)

	result.Seek(9, io.SeekStart)
	buf := make([]byte, 4)
	result.Read(buf)
	assert.Equal(t, []byte("test"), buf)
}

func TestBuildReadSeekCloser_WithReadAll_ErrorPropagates(t *testing.T) {
	reader := &errorReader{}

	result, err := BuildReadSeekCloser(reader)
	assert.Nil(t, result)
	assert.Error(t, err)
}

type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

type readSeekCounter struct {
	readCount int
	io.ReadSeeker
}

func (r *readSeekCounter) Read(p []byte) (n int, err error) {
	r.readCount++
	return r.ReadSeeker.Read(p)
}
