package Jwtutil

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"io/ioutil"
	"log"
)

func GetJWTToken(priKeyPath string,claims jwt.Claims) string  {
	priKeyBytes,err:=ioutil.ReadFile(priKeyPath)
	if err!=nil{
		log.Fatal("私钥文件读取失败")
	}
	priKey,err:=jwt.ParseRSAPrivateKeyFromPEM(priKeyBytes)
	if err!=nil{
		log.Fatal("私钥文件不正确")
	}

	token_obj :=jwt.NewWithClaims(jwt.SigningMethodRS256,claims)
	token,_:=token_obj.SignedString(priKey)
	return token
}

func PraseJWTToken(priKeyPath string,token string, claims jwt.Claims) error{
	pubKeyBytes,err:=ioutil.ReadFile(priKeyPath)
	if err!=nil{
		log.Fatal("公钥文件读取失败")
	}
	pubKey,err:=jwt.ParseRSAPublicKeyFromPEM(pubKeyBytes)
	if err!=nil{
		log.Fatal("公钥文件不正确")
	}

	getToken,err:= jwt.ParseWithClaims(token,claims, func(token *jwt.Token) (i interface{}, e error) {
		return pubKey,nil
	})
	if err!=nil{
		return err
	}

	if getToken.Valid{
		return nil
	}
	return fmt.Errorf("数据不合法")
}