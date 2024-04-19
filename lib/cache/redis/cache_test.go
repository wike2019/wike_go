package redis

import (
	"fmt"
	"testing"
)

type Test struct {
	Name string
}

func TestNewCache(t *testing.T) {

	c := NewCache("127.0.0.1:6379")
	// Set the value of the key "foo" to "bar", with the default expiration time

	a := &Test{
		Name: "caasc",
	}
	c.Set("a", a, NoExpiration)
	c.Set("a1", "1111", NoExpiration)
	b := &Test{}
	c.Get("a", b)
	td := ""
	c.Get("a1", &td)
	fmt.Println(b.Name)
	fmt.Println(td)
}
