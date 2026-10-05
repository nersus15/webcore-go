package port

import (
	"context"
	"time"
)

type RateLimitResult struct {
	Allowed   bool
	Remaining int64
	ResetAt   time.Time
}

// sliding window counter; store bersama (Redis) membuat limit berlaku lintas pod
type IRateLimitStore interface {
	RateLimit(ctx context.Context, key string, limit int64, window time.Duration) (RateLimitResult, error)
}
