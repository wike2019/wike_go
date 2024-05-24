package jwt

import (
	"fmt"
	"testing"
	"time"
)

type Info struct {
	Name string
}

func TestJWT(t *testing.T) {
	core := Info{
		Name: "我是wike",
	}
	//生成一个token
	token, err := Create[Info](core, time.Second*10)
	fmt.Println(token, err)
	//根据token得到对象
	info, err := Parse[Info](token)
	fmt.Println(info, err)
}
