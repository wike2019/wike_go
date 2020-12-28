package main

import (
	"github.com/wike2019/wike_go/src/core/Ioc"
)

func main()  {
	Ioc.New().Config(NewEtcdConfig())
	Ioc.New().ApplyAll()
	call:=new(EtcdCall)
	call=Ioc.New().Get(call).(*EtcdCall)

	//续期
	call.KeepAlive()//续期

	//删除
	call.DelLease()//删除租约
	call.DelWithPrevKV()//

	//设置
	call.SetWithTimeOut() //有过期时间
	call.Set() //没有过期时间
	call.SetLeaseId()//设置租约
	call.SetLeaseId()//组合使用

	//取值
	call.Get() //没有prev
	call.GetWithPrevKV()  //有prev
}
