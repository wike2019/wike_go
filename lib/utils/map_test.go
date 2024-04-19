package utils

import (
	"fmt"
	"testing"
)

func TestMap(t *testing.T) {
	mapstring := NewMap[string]()
	mapstring.Set("111", "2222")
	mapstring.Set("222", "3333")
	fmt.Println(mapstring.Keys())
	fmt.Println(mapstring.Values())
	fmt.Println(mapstring.Get("111"))

}
