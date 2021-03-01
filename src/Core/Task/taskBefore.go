package Task

import (
	"fmt"
	"time"
)
type InChan chan interface{}

func TaskMustBefore(job func(in InChan),d time.Duration) (interface{},  error) {
	ret:=make(InChan)
	go job(ret)
	select {
	case r:=<-ret:
		return r,nil
	case <-time.After(d):
		return nil,fmt.Errorf("time out")
	}
}