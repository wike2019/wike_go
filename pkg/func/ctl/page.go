package ctl

import (
	"math"
)

type Page struct {
	SumPage  int `json:"sumPage"  desc:"总页数"`                              // 总页数
	SumCount int `json:"sumCount" desc:"总条数"`                              // 总条数
	CurPage  int `json:"page" form:"page" desc:"当前页"  required:"true"  `   // 当前页
	Offset   int `json:"-"`                                                // 起始量
	Count    int `json:"count" form:"count" desc:"每页返回数量" required:"true"` // 每页返回数量
}

// PageTime 适用 列表-分页-时间范围
type Time struct {
	StartTime int64 `json:"startTime" form:"startTime"` // 开始时间
	EndTime   int64 `json:"endTime" form:"endTime"`     // 结束时间
}

// FormatPage 格式化分页数据-起始量

// FormatTotal 格式化分页数据-总页数
func (p *Page) FormatTotal(total int64) {
	p.SumCount = int(total)
	if p.Count > 0 {
		p.SumPage = int(math.Ceil(float64(p.SumCount) / float64(p.Count)))
	} else {
		p.SumPage = 1
	}
}
