package Validate

import (
	"github.com/gin-gonic/gin/binding"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strings"
	"sync"
)

type Validate struct {
	Validate *validator.Validate
	Trans ut.Translator
}


var once sync.Once
var CheckVaid *Validate
func New() *Validate {
	once.Do(func() {
		v:=validator.New()
		v.SetTagName("binding")
		CheckVaid= &Validate{Validate:v}
	})
	return CheckVaid
}

func (this * Validate)AddValiDate(name string,checkfunc validator.Func)  {
	this.Validate.RegisterValidation(name,checkfunc)
	binding.Validator.Engine().(*validator.Validate).RegisterValidation(name,checkfunc)
}

func (this * Validate)Msg(obj interface{},errs error)string  {

	err:=errs.(validator.ValidationErrors)[0]
	field:=err.Field()
	name:=err.Tag()
	return this.GetValidMsg(obj,field,name,err.Error())

}


func (this * Validate) GetValidMsg(obj interface{},field string,tagName string,err string) string {

	getObj := reflect.TypeOf(obj)
	if f,exist:=getObj.Elem().FieldByName(field);exist{
		msg:=f.Tag.Get("vmsg");
		msginfo:=strings.Split(msg,",")
		for _,value:=range msginfo  {
			errMsg:=strings.Split(value,"=")
			if errMsg[0]==tagName{
				return  errMsg[1]
			}
		}
	}
	return err
}