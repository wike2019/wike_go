package rateLimiter

import (
	"github.com/gin-gonic/gin"
	"github.com/golang/groupcache/lru"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
	"net/http"
	"sync"
)

type RateLimiterCache struct {
	cache *lru.Cache
	mu    sync.Mutex
}

func NewRateLimiterCache(maxEntries int) *RateLimiterCache {
	return &RateLimiterCache{
		cache: lru.New(maxEntries),
	}
}

func (rlc *RateLimiterCache) GetLimiter(key string, r rate.Limit, b int) *rate.Limiter {
	rlc.mu.Lock()
	defer rlc.mu.Unlock()

	if lim, ok := rlc.cache.Get(key); ok {
		return lim.(*rate.Limiter)
	}

	lim := rate.NewLimiter(r, b)
	rlc.cache.Add(key, lim)
	return lim
}
func RateLimiter(rateLimiterCache *RateLimiterCache, cfg *viper.Viper) gin.HandlerFunc {
	r := cfg.GetFloat64("LimitRate")
	b := cfg.GetInt("LimitBucket")
	return func(c *gin.Context) {
		limiter := rateLimiterCache.GetLimiter(c.ClientIP(), rate.Limit(r), b)
		if limiter.Allow() == false {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":     http.StatusTooManyRequests,
				"error":    "请求太频繁了",
				"trace_id": c.GetString("trace_id"),
			})
			return
		}
		c.Next()
	}
}
