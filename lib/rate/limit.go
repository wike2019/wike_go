package rate

import (
	"github.com/spf13/viper"
	"github.com/wike2019/wike_go/lib/rateLimiter"
)

func Limit(cfg *viper.Viper) *rateLimiter.RateLimiterCache {
	limit := cfg.GetInt("LRULimit")
	return rateLimiter.NewRateLimiterCache(limit)
}
