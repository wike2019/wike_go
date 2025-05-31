package common

// 根据id 搜索
type IDSearch struct {
	ID uint `json:"id" form:"id"`
}

// 根据关键字搜索
type KeyWordsSearch struct {
	KeyWords string `json:"keywords" form:"keywords"`
}

// 空结构,一般用于参数为空
type Empty struct {
}

func (this IDSearch) GetID() uint {
	return this.ID
}

type IDSearcher interface {
	GetID() uint
}
