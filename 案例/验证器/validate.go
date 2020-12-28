package main

import (
	"github.com/go-playground/validator/v10"
	"unicode/utf8"
)

func CheckName(f validator.FieldLevel) bool { // FieldLevel contains all the information and helper functions to validate a field
	count := utf8.RuneCountInString(f.Field().String()) //通过utf8编码，获取字符串长度
	if count >= 2 && count <= 5 {
		return true
	}
	return false
}
