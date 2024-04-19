package utils

import (
	"encoding/json"
)

// set 集合
type Set struct {
	mapData *MapSync[Empty]
}
type Empty struct{}

func NewSet() *Set {
	return &Set{
		mapData: NewMap[Empty](),
	}
}
func (this *Set) Add(value interface{}) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	this.mapData.Set(string(b), Empty{})
	return err
}
func (this *Set) List(dist interface{}) error {
	keys := this.mapData.Keys()
	res := make([]interface{}, len(keys))
	for i, item := range keys {
		var obj interface{}
		err := json.Unmarshal([]byte(item), &obj)
		if err != nil {
			return err
		}
		res[i] = obj
	}
	b, err := json.Marshal(res)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, dist)
	if err != nil {
		return err
	}
	return nil
}
