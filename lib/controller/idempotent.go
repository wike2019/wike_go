package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike-go/MysqlInit"
	"github.com/wike2019/wike-go/lib/bloom"
	"gorm.io/gorm"
	"time"
)

// 一致性避免重试不幂等
func IdempotentCheck(context *gin.Context, tx *gorm.DB) error {
	traceId, _ := context.Get("trace_id")
	data := &MysqlInit.IdempotentCheck{
		TraceId: traceId.(string),
	}
	return tx.Create(data).Error
}

func Clear(db *gorm.DB) {
	for _ = range bloom.Clear {
		//定期清理老数据
		for {
			result := db.Where("time < ?", time.Now().AddDate(0, 0, -2)).Limit(5000).Delete(&MysqlInit.IdempotentCheck{})
			if result.Error != nil {
				continue
			}
			if result.RowsAffected < 5000 {
				break
			}
		}
	}
}

// 开启幂等服务
func IdempotentCheckStart(db *gorm.DB) {
	db.AutoMigrate(&MysqlInit.IdempotentCheck{})
	Clear(db)
}
