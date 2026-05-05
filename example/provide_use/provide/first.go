package provide

import "github.com/wike2019/wike_app/model"

func ProvideExampleFirst(dog *model.Dog) *Animal {
	return &Animal{dog: dog}
}
