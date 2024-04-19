package utils

import (
	"fmt"
	"testing"
)

type Test struct {
	Name string
}

func TestSet(t *testing.T) {
	set := NewSet()
	set.Add("111")
	set.Add("222")
	set.Add("333")
	set.Add("222")
	var str []string
	set.List(&str)
	for i, item := range str {
		fmt.Println(i, item)
	}
	set2 := NewSet()
	obj1 := &Test{
		Name: "a1",
	}
	obj2 := &Test{
		Name: "a1",
	}
	obj3 := &Test{
		Name: "a2",
	}
	set2.Add(obj1)
	set2.Add(obj2)
	set2.Add(obj3)
	var obj []*Test
	set2.List(&obj)
	for i, item := range obj {
		fmt.Println(i, item)
	}
}
