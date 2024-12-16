package ctl

type PageDocList[T any] struct {
	List T `json:"list"`
	Page Page
}
type HttpReturn struct {
	Code    errorCode   `json:"code"`
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"trace_id"`
}

type HttpDoc[T any] struct {
	Code    errorCode `json:"code"`
	Msg     string    `json:"msg"`
	Data    T         `json:"data"`
	TraceId string    `json:"trace_id" desc:"追踪id"`
}

func PageDoc[T any]() HttpDoc[PageDocList[T]] {
	return HttpDoc[PageDocList[T]]{}
}
