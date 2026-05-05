package model

import "errors"

type ModelSum struct {
	Num1 float64
	Num2 float64
}

func (this *ModelSum) Check() error {
	if this.Num1 < 0 {
		return errors.New("第一个数不能小于零")
	}
	if this.Num2 < 0 {
		return errors.New("第二个数不能小于零")
	}
	return nil
}
