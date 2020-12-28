package main

import (
	"fmt"
	"github.com/wike2019/wike_go/src/core/Etcd"
	"github.com/wike2019/wike_go/src/core/Ioc"
	"sync"
	"time"
)

func main() {
	Ioc.New().Config(NewEtcdConfig())
	Ioc.New().ApplyAll()
	Instance := Etcd.EtcdCache()
	defer Etcd.ReleaseEtcdCache(Instance)
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		time.Sleep(time.Second * 1)
		defer wg.Done()
		err := Instance.Lock("lock", func(params ...interface{}) {

			fmt.Println(params)
			time.Sleep(5 * time.Second)
		},
			"wike", 1)
		fmt.Println(err)
	}()
	wg.Add(1)
	go func() {

		defer wg.Done()
		err := Instance.Lock("lock", func(params ...interface{}) {

			fmt.Println(params)
			time.Sleep(5 * time.Second)
		},
			"wike", 1)
		fmt.Println(err)

	}()

	wg.Wait()
}