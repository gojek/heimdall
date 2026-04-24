package internal_test

import (
	"errors"
	"testing"

	"github.com/gojek/heimdall/v7/internal"
	"github.com/stretchr/testify/assert"
)

var err1 = errors.New("err1")
var err2 = errors.New("err2")
var err3 = errors.New("err3")
var err4 = errors.New("err4")

func TestBuildMultiError(t *testing.T) {
	t.Parallel()

	assert.Nil(t, internal.BuildMultiError(nil))

	assert.Equal(t, err1, internal.BuildMultiError([]error{err1}))

	errs := internal.BuildMultiError([]error{err1, err2})
	assert.Equal(t, "err1, err2", errs.Error())
	assert.ErrorIs(t, errs, err1)
	assert.ErrorIs(t, errs, err2)
	assert.NotErrorIs(t, errs, err3)
	assert.NotErrorIs(t, errs, err4)

	errs = internal.BuildMultiError([]error{err1, err2, err3, err4})
	assert.Equal(t, "err1, err2, err3, err4", errs.Error())
	assert.ErrorIs(t, errs, err1)
	assert.ErrorIs(t, errs, err2)
	assert.ErrorIs(t, errs, err3)
	assert.ErrorIs(t, errs, err4)
}
