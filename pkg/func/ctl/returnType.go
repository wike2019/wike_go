package ctl

type PageList[T any] struct {
	Data    []T       `json:"data"`
	Page    Page      `json:"page"`
	Msg     string    `json:"msg"`
	TraceId string    `json:"trace_id" desc:"追踪id"`
	Code    errorCode `json:"code"`
}

type DataList[T any] struct {
	Msg     string    `json:"msg"`
	TraceId string    `json:"trace_id" desc:"追踪id"`
	Code    errorCode `json:"code"`
	Data    T         `json:"data"`
}
