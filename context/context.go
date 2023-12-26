package context

import (
	"context"
	"time"
)

var (
	DefaultContextTimeout = time.Second * 100
)

func NewTimeoutContext(ctx context.Context, timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(ctx, timeout)
}

func NewDefaultTimeoutContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), DefaultContextTimeout)
}
