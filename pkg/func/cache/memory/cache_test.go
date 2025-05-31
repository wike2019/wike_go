package memory

import (
	"fmt"
	"testing"
)

type Test struct {
	Name string
}

func TestNewCache(t *testing.T) {
	c := NewCache()
	// Set the value of the key "foo" to "bar", with the default expiration time

	a := &Test{
		Name: "caasc",
	}
	c.Set("a", a, NoExpiration)
	c.Set("a1", "11131", NoExpiration)
	b := &Test{}
	c.Get("a", b)
	td := ""
	c.Get("a1", &td)
	fmt.Println(b.Name)
	fmt.Println(td)
}
