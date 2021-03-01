package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"


	"time"
)





type UserClaim struct {
	Age float64  //这里是坑 jwt转换过后 Int被转换成float64 所以 int改成float64
	Uname string
	jwt.StandardClaims
}

func main()  {
	Help.GenRSAPubAndPri(1024,"./pem")
	user:=UserClaim{Uname:"wike",Age:33333};
	user.ExpiresAt=time.Now().Add(time.Second*5).Unix() //添加过期时间

	token:= Jwtutil.GetJWTToken("./pem/private.pem",user )
	fmt.Println(token)
	user2:=&UserClaim{};

	err, _ := Jwtutil.PraseJWTToken("./pem/public.pem",token,user2) //解密
	if err !=nil{
		fmt.Println(err)
	}
	fmt.Println(user2)
}
