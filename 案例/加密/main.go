package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/wike2019/wike_go/src/util/Crypto"
)

func main()  {

	crypto := Crypto.New()
	origData := []byte("Hello World 11111233") // 待加密的数据
	fmt.Println("原文：", string(origData))
	fmt.Println("------------------ CBC模式 --------------------")
	encrypted := crypto.AesEncryptCBC(origData)
	fmt.Println("密文(hex)：", hex.EncodeToString(encrypted))
	fmt.Println("密文(base64)：", base64.StdEncoding.EncodeToString(encrypted))
	decrypted := crypto.AesDecryptCBC(encrypted)
	fmt.Println("解密结果：", string(decrypted))


	fmt.Println(crypto.Md5("wike is ok"))

	fmt.Println(crypto.Sha256("wike is ok"))


	crypto.RSAGenKey(1024,"./public/key")

	data,_:=crypto.EncyptogRSA(origData,"./public/key/publicKey.pem")
	fmt.Println("加密之后的数据为：",string(data))
	data,_=crypto.DecrptogRSA(data,"./public/key/privateKey.pem")
	fmt.Println("解密之后的数据为：",string(data))


}

