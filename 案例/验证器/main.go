package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/wike2019/wike_go/src/util/Validate"
)

func main()  {
	Validate.New().AddValiDate("CheckName", CheckName)
	Validate.New().AddValiDate("NewEmail", func(fl validator.FieldLevel) bool {
		return  true
	})
	user:=&User{Name:"不合法的名字"}
	err:=Validate.New().Validate.Struct(user)

	fmt.Println(Validate.New().Msg(user,err))

}
