package ctl

type PageDocList[T any] struct {
	List T `json:"list"`
	Page Page
}
type HttpReturn struct {
	Msg     string      `json:"msg"`
	Data    interface{} `json:"data"`
	TraceId string      `json:"trace_id"`
	Code    errorCode   `json:"code"`
}

type HttpDoc[T any] struct {
	Msg     string    `json:"msg"`
	TraceId string    `json:"trace_id" desc:"追踪id"`
	Code    errorCode `json:"code"`
	Data    T         `json:"data"`
}

func PageDoc[T any]() HttpDoc[PageDocList[T]] {
	return HttpDoc[PageDocList[T]]{}
}

type HttpDocList[T any] struct {
	Msg     string `json:"msg"`
	TraceId string `json:"trace_id" desc:"追踪id"`
	Code    errorCode
	Data    []T `json:"data"`
}
