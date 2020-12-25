package Validate

import (
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh2 "github.com/go-playground/validator/v10/translations/zh"
	"reflect"
	"sync"
)

type Validate struct {
	Validate *validator.Validate
	msg map[string]string
	Trans ut.Translator
}
var once sync.Once
var CheckVaid *Validate
func New() *Validate {
	once.Do(func() {
		uni := ut.New(zh.New(), zh.New())
		v:=validator.New()
		v.SetTagName("checking")
		trans, _ := uni.GetTranslator("zh")
		zh2.RegisterDefaultTranslations(v, trans)
		CheckVaid= &Validate{Validate:v,msg: map[string]string{},Trans:trans}
	})
	return CheckVaid
}

func (this * Validate)AddValiDate(name string,checkfunc validator.Func,msg string)  {
	this.msg[name]=msg
	this.Validate.RegisterValidation(name,checkfunc)
}

func (this * Validate)Msg(obj interface{},errs error)string  {

	err:=errs.(validator.ValidationErrors)[0]

	field:=err.Field()
	name:=err.Tag()
	errMsg:=err.Translate(CheckVaid.Trans)
	msg,ok:=this.msg[name]
	if !ok {
		return this.GetValidMsg(obj,field,errMsg)
	}
	return  msg
}


func (this * Validate) GetValidMsg(obj interface{},field string,err string) string {

	getObj := reflect.TypeOf(obj)
	if f,exist:=getObj.Elem().FieldByName(field);exist{
		msg:=f.Tag.Get("vmsg");
		if  msg!=""{
			return  msg
		}
		return err
	}
	return err
}