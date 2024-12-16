package ctl

import (
	"github.com/gin-gonic/gin"
	"github.com/wike2019/wike_go/plugin"
	"gorm.io/gorm"
	"time"
)

// 一致性避免重试不幂等
func RepeatCheck(context *gin.Context, tx *gorm.DB) error {
	traceId, _ := context.Get("trace_id")
	data := &plugin.RepeatCheck{
		TraceId: traceId.(string),
	}
	return tx.Create(data).Error
}

var ClearChan chan struct{}

func init() {
	ClearChan = make(chan struct{}, 10)
}

// todo
func Clear(db *gorm.DB) {
	for _ = range ClearChan {
		//定期清理老数据
		for {
			result := db.Where("time < ?", time.Now().AddDate(0, 0, -2)).Limit(5000).Delete(&plugin.RepeatCheck{})
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
func RepeatCheckStart(db *gorm.DB) {
	db.AutoMigrate(&plugin.RepeatCheck{})
	Clear(db)
}

// 关闭幂等服务
func RepeatCheckStop() {
	close(ClearChan)
}
