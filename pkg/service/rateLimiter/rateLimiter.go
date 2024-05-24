package rateLimiter

import (
	"github.com/golang/groupcache/lru"
	"golang.org/x/time/rate"
	"sync"
)

type LimiterCache struct {
	cache *lru.Cache
	mu    sync.Mutex
}

func NewRateLimiterCache(maxEntries int) GetLimit {
	return &LimiterCache{
		cache: lru.New(maxEntries),
	}
}

func (rlc *LimiterCache) GetLimiter(key string, r rate.Limit, b int) *rate.Limiter {
	rlc.mu.Lock()
	defer rlc.mu.Unlock()

	if lim, ok := rlc.cache.Get(key); ok {
		return lim.(*rate.Limiter)
	}

	lim := rate.NewLimiter(r, b)
	rlc.cache.Add(key, lim)
	return lim
}
func RateLimit() GetLimit {
	return NewRateLimiterCache(40960)
}

type GetLimit interface {
	GetLimiter(key string, r rate.Limit, b int) *rate.Limiter
}
