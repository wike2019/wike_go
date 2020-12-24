package web

import (
	"encoding/json"
	"log"
)
// 模型
type Model interface {
	String() string
}

type Models string

func MakeModels(v interface{}) Models {
	b, err := json.Marshal(v)
	if err != nil {
		log.Println(err)
	}
	return Models(b)
}
