package limiter

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/Wammero/Balancer/internal/models"
)

type Limiter interface {
	Check(ctx context.Context, bucket *models.TokenBucket) error
}

type tokenBucketLimiter struct {
	mu sync.Mutex
}

func New() Limiter {
	return &tokenBucketLimiter{}
}

func (l *tokenBucketLimiter) Check(ctx context.Context, bucket *models.TokenBucket) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	elapsed := now.Sub(bucket.LastRefill).Seconds()

	newTokens := int(elapsed * float64(bucket.TokensPerSecond))
	if newTokens > 0 {
		bucket.Tokens = min(bucket.Tokens+newTokens, bucket.Capacity)
		bucket.LastRefill = now
	}

	if bucket.Tokens <= 0 {
		return fmt.Errorf("rate limit exceeded for client %s", bucket.Key)
	}

	bucket.Tokens--
	bucket.UpdatedAt = time.Now()

	return nil
}
