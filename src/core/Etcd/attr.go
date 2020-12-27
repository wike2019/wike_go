package Etcd

import (
	"fmt"
	"github.com/wike2019/wike_go/src/Result"
	"go.etcd.io/etcd/clientv3"
)
type empty struct {}
//可添加的属性
const (
	ATTR_WithPrevKV="WithPrevKV"   //过期时间
	ATTR_Lease ="Lease"  // setnx
	ATTR_WithTime ="WithTime"  // setxx
)
//属性
type OperationAttr struct {
	Name string
	Value interface{}
}
type OperationAttrs []*OperationAttr
//查找对应属性限制传入的参数
func(this OperationAttrs) Find(name string) *Result.ErrorResult { //查找指定属性
	for _,attr:=range this {
		if attr.Name==name{
			return Result.Result(attr.Value,nil)
		}
	}
	return Result.Result(nil,fmt.Errorf("OperationAttrs found error:%s",name))
}
func WithPrevKV() *OperationAttr  {  //设置前缀
	return &OperationAttr{Name:ATTR_WithPrevKV,Value:empty{}}
}
func WithLease(Lease clientv3.LeaseID) *OperationAttr  {  //设置租约
	return &OperationAttr{Name:ATTR_Lease,Value:Lease}
}
func WithTime(time int64) *OperationAttr  {  //设置过期时间
	return &OperationAttr{Name:ATTR_WithTime,Value:time}
}
