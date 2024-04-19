package jwt

import (
	"context"
	"fmt"
	"testing"
	"time"
)

type Info struct {
	Name string
}

func TestJWT(t *testing.T) {
	core := Info{
		Name: "123",
	}
	token, err := Create[Info](core, time.Second*10)
	fmt.Println(token, err)
	info, err := Parse[Info](token)
	fmt.Println(info, err)

	fmt.Println(info.Name, err)
	c := context.Background()
	c2 := context.WithValue(c, "i", info)
	data := c2.Value("i")

	fmt.Println(data.(*Info).Name)
}
