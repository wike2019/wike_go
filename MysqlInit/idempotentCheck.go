package MysqlInit

import (
	"gorm.io/gorm"
	"time"
)

type IdempotentCheck struct {
	ID      uint64    `gorm:"primaryKey;column:id;comment:主键"`
	TraceId string    `gorm:"column:trace_id;type:varchar(255);uniqueIndex:idx_id;comment:追踪id"`
	Time    time.Time `gorm:"column:time;type:date;uniqueIndex:idx_id;comment:插入日期"`
}

func (this *IdempotentCheck) TableName() string {
	return "IdempotentCheck"
}

// BeforeCreate Gorm钩子，在创建记录之前调用
func (m *IdempotentCheck) BeforeCreate(tx *gorm.DB) (err error) {
	m.Time = time.Now()
	return
}

// BeforeUpdate Gorm钩子，在更新记录之前调用
func (m *IdempotentCheck) BeforeUpdate(tx *gorm.DB) (err error) {
	m.Time = time.Now()
	return
}
