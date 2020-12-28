package main

import (
	"fmt"
	"github.com/wike2019/wike_go/src/util/Archive"
	"github.com/wike2019/wike_go/src/util/Compress"
)

func main()  {
	instance:=Archive.New()
	fmt.Println(instance.TarFile("./用于压缩测试/my.txt","./用于压缩测试.tar"))
	instance.UnTarFile("./用于压缩测试.tar","./用于压缩解压测试")
	instance.TarDir("./用于压缩测试","./用于压缩测试目录.tar")
	instance.UnTarDir("./用于压缩测试目录.tar","./压缩解压测试tmp")
	instance.Zip("./用于压缩测试","./压缩解压测试tmp.zip")
	instance.UnZip("./压缩解压测试tmp.zip","./压缩解压测试zip")


	compress:=Compress.New()

	compress.Gzip("./用于压缩测试","./用于压缩测试.gzip")
	compress.UnGzip("./用于压缩测试.gzip","./用于压缩测试gzip")

	data:=compress.Zlib([]byte("key"))
	fmt.Println(data)
	data,_=compress.UnZlib(data)
	fmt.Println(string(data))



}
