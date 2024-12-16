package ctl

import (
	"math"
	"sync"
)

type Page struct {
	SumPage  int `json:"sumPage"  desc:"总页数"`                              // 总页数
	SumCount int `json:"sumCount" desc:"总条数"`                              // 总条数
	CurPage  int `json:"page" form:"page" desc:"当前页"  required:"true"  `   // 当前页
	Offset   int `json:"-"`                                                // 起始量
	Count    int `json:"count" form:"count" desc:"每页返回数量" required:"true"` // 每页返回数量
}

type Pagination struct {
	List interface{} `json:"list"`
	Page
}

func (this *Pagination) Reset() {
	this.List = nil
	this.SumPage = 0
	this.SumCount = 0
	this.CurPage = 0
	this.Offset = 0
	this.Count = 0
}

// PageTime 适用 列表-分页-时间范围
type PageTime struct {
	Page
	StartTime int64 `json:"startTime" form:"startTime"` // 开始时间
	EndTime   int64 `json:"endTime" form:"endTime"`     // 结束时间
}

var PageResult *sync.Pool

func init() {
	PageResult = &sync.Pool{
		New: func() interface{} {
			return &Pagination{}
		},
	}
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
