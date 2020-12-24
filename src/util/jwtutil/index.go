package jwtutil

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
)

func GetJWTToken(priKeyPath string,call func() *jwt.Token) string  {
	priKeyBytes,err:=ioutil.ReadFile(priKeyPath)
	if err!=nil{
		log.Fatal("私钥文件读取失败")
	}
	priKey,err:=jwt.ParseRSAPrivateKeyFromPEM(priKeyBytes)
	if err!=nil{
		log.Fatal("私钥文件不正确")
	}

	token_obj :=call()
	token,_:=token_obj.SignedString(priKey)
	return token
}

func PraseJWTToken(priKeyPath string,token string) (jwt.MapClaims,error){
	pubKeyBytes,err:=ioutil.ReadFile(priKeyPath)
	if err!=nil{
		log.Fatal("公钥文件读取失败")
	}
	pubKey,err:=jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err!=nil{
		log.Fatal("公钥文件不正确")
	}

	getToken,_:= jwt.Parse(token, func(token *jwt.Token) (i interface{}, e error) {
		return pubKey,nil
	})
	if getToken.Valid{
		return getToken.Claims.(jwt.MapClaims),nil
	}
	return nil,fmt.Errorf("数据不合法")
}