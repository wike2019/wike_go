package provide

import (
	"fmt"

	"github.com/wike2019/wike_app/model"
)

type Animal struct {
	dog *model.Dog
}

func (this *Animal) Show() string {
	fmt.Println(this.dog.Name)
	return this.dog.Name + "这个是聚合方法"
}
func ProvideExampleSecond(dog *model.Dog) *Animal {
	return &Animal{dog: dog}
}
