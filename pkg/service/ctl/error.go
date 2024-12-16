package ctl

type StatusError struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func Error(err error, code int) {
	if err != nil {
		panic(StatusError{Msg: err.Error(), Code: code})
	}
}
