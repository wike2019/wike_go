package http

import (
	"fmt"
	"github.com/wike2019/wike_go/src/util/Pip"
	"time"
)

func GetPage(args ...interface{}) Pip.InChan  {
	in:=make(Pip.InChan)
	go func() {
		defer close(in)
		//模拟数据库取数据
		for i:=11;i>0;i--{
			in<-i
		}
	}()
	return in
}

func DoData(in Pip.InChan) Pip.OutChan{
	out:=make(Pip.OutChan)
	go func() {
		defer close(out)
		for d:=range in {
			//模拟导入到es
			time.Sleep(time.Second*1)
			out<-fmt.Sprintf("处理了%d条数据,%d\n",d.(int)*1000,time.Now().Unix())
		}
	}()
	return out
}
func Do()  Pip.OutChan{
	p2:= Pip.New()
	p2.SetCmd(GetPage)
	p2.SetPipeCmd(DoData,5)
	out2:=p2.Exec()
	return out2
}






