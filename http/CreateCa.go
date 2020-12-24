package http

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/wike2019/wike_go/src/util/help"
	"github.com/wike2019/wike_go/src/util/jwtutil"
	"time"
)
type UserClaim struct {
	Age float64  //这里是坑 jwt转换过后 Int被转换成float64 所以 int改成float64
	Uname string
	jwt.StandardClaims
}

func CreateCA() {
	help.CreateCA([]string{"wike.com", "192.168.3.3"},"./keys")
	help.GenRSAPubAndPri(1024,"./pem")
	user:=UserClaim{Uname:"wike",Age:111};
	user.ExpiresAt=time.Now().Add(time.Second*5).Unix()
	token:=jwtutil.GetJWTToken("./pem/private.pem", func() *jwt.Token {
		return jwt.NewWithClaims(jwt.SigningMethodRS256,user)
	})
	user2:=&UserClaim{};
	data,err:=jwtutil.PraseJWTToken("./pem/public.pem",token)
	if err !=nil{
		fmt.Println(err)
	}
	help.Map2Struct(data,user2)
	fmt.Println(user2)
}
