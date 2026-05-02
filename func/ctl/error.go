package ctl

import (
	"errors"
	"gorm.io/gorm"
)

type StatusError struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
}

func Error(err error, code int) {
	if err != nil {
		panic(StatusError{Msg: err.Error(), Code: code})
	}
}
func (r *Controller) IsExist(err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
func IsExist(err error) bool {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}
	return true
}
