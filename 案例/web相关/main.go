package main

import (

	"github.com/wike2019/wike_go/src/Web"
	"github.com/wike2019/wike_go/src/core/Ioc"
	"os"
	"os/signal"
	"syscall"
)

func main()  {
	signalChan := make(chan os.Signal, 1)
	Ioc.New().Config(NewDBConfig())
	app:= Web.New(). //初始化脚手架
		Attach(NewTokenCheck(), NewAddVersion()).
		Mount("/string", NewStringController()).
	    Mount("/json", NewJsonController()).
		Mount("/sql", NewSqlController()).
		Mount("/middle", NewMiddleController())

	Ioc.New().ApplyAll()
	go func() {
		app.Launch()
	}()
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	//关闭工作
	<-signalChan
}

