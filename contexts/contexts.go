package contexts

import (
	"context"
	"time"
)

var (
	DefaultContextTimeout = time.Second * 100
)

func NewTimeoutContext(ctx context.Context, timeouts ...time.Duration) (context.Context, context.CancelFunc) {
	timeout := DefaultContextTimeout
	if len(timeouts) > 0 {
		timeout = timeouts[0]
	}
	return context.WithTimeout(ctx, timeout)
}
