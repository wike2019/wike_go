package utils

import (
	"math/rand"
	"time"
)

var letters = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// RandomString 生成随机长度的字符串
func RandomString(length int64) string {
	b := make([]rune, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := range b {
		b[i] = letters[r.Intn(62)]
	}
	return string(b)
}

// GetRandomNum 生成随机数字(包括边界)
func GetRandomNum(min, max int) int64 {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return int64(rand.Intn(max-min+1) + min)
}
