package internal_test

import (
	"context"
	"testing"
	"time"

	"github.com/gojek/heimdall/v7/internal"
	"github.com/stretchr/testify/assert"
)

func TestIsCtxDone(t *testing.T) {
	assert.False(t, internal.IsCtxDone(context.Background()))

	ctx, cancel := context.WithCancel(context.Background())
	assert.False(t, internal.IsCtxDone(ctx))
	cancel()
	assert.True(t, internal.IsCtxDone(ctx))

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	assert.False(t, internal.IsCtxDone(ctx))
	time.Sleep(12 * time.Millisecond)
	assert.True(t, internal.IsCtxDone(ctx))
}
