package Redis

import (
	"fmt"
	"github.com/wike2019/wike_go/src/result"
	"time"
)

type empty struct {}
const (
	ATTR_EXPIRE="expr"   //过期时间
	ATTR_NX ="nx"  // setnx
	ATTR_XX ="xx"  // setxx
)
//属性
type OperationAttr struct {
	Name string
	Value interface{}
}
type OperationAttrs []*OperationAttr
func(this OperationAttrs) Find(name string) *Result.ErrorResult { //查找指定属性
	 for _,attr:=range this {
	 	if attr.Name==name{
	 		return Result.Result(attr.Value,nil)
		}
	 }
	 return Result.Result(nil,fmt.Errorf("OperationAttrs found error:%s",name))
}
func WithExpire(t time.Duration) *OperationAttr  {  //设置过期时间
		return &OperationAttr{Name:ATTR_EXPIRE,Value:t}
}
func WithNX() *OperationAttr  {  //设置set 的NX
	return &OperationAttr{Name:ATTR_NX,Value:empty{}}
}
func WithXX() *OperationAttr  {  //设置set 的XX
	return &OperationAttr{Name:ATTR_XX,Value:empty{}}
}
