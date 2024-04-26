package rate

import (
	"github.com/spf13/viper"
	"github.com/wike2019/wike_go/lib/rateLimiter"
	"golang.org/x/time/rate"
)

func Limit(cfg *viper.Viper) *rateLimiter.RateLimiterCache {
	limit := cfg.GetInt("LRULimit")
	return rateLimiter.NewRateLimiterCache(limit)
}

func CreateLimiter(r rate.Limit, b int) *rate.Limiter {
	return rate.NewLimiter(r, b)
}
