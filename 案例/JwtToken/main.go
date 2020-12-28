package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/wike2019/wike_go/src/util/Help"
	"github.com/wike2019/wike_go/src/util/Jwtutil"
	"time"
)





type UserClaim struct {
	Age float64  //这里是坑 jwt转换过后 Int被转换成float64 所以 int改成float64
	Uname string
	jwt.StandardClaims
}

func main()  {
	Help.GenRSAPubAndPri(1024,"./pem")
	user:=UserClaim{Uname:"wike",Age:111};
	user.ExpiresAt=time.Now().Add(time.Second*5).Unix() //添加过期时间

	token:= Jwtutil.GetJWTToken("./pem/private.pem", func() *jwt.Token {
		return jwt.NewWithClaims(jwt.SigningMethodRS256,user)
	})

	user2:=&UserClaim{};

	data,err:= Jwtutil.PraseJWTToken("./pem/public.pem",token) //解密
	if err !=nil{
		fmt.Println(err)
	}
	Help.Map2Struct(data,user2) //将map 映射为结构体

	fmt.Println(user2)
}
