package Crypto

import (
"bytes"
"crypto/aes"
"crypto/cipher"
"crypto/dsa"
"crypto/md5"
"crypto/rand"
"crypto/rsa"
"crypto/sha256"
"crypto/x509"
"encoding/base64"
"encoding/pem"
"fmt"
"math/big"
"os"
)

type  Crypto struct {

}

func New() *Crypto {
	return &Crypto{}
}

func pkcs5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}
func pkcs5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}


func (this *Crypto) AesEncryptCBC(origData []byte) (encrypted []byte) {
	// 分组秘钥
	// NewCipher该函数限制了输入k的长度必须为16, 24或者32
	key := []byte("ABCDEFGAIJKLMNOP")
	block, _ := aes.NewCipher(key)
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	origData = pkcs5Padding(origData, blockSize)                // 补全码
	blockMode := cipher.NewCBCEncrypter(block, key[:blockSize]) // 加密模式
	encrypted = make([]byte, len(origData))                     // 创建数组
	blockMode.CryptBlocks(encrypted, origData)                  // 加密
	return encrypted
}

func (this *Crypto) AesDecryptCBC(encrypted []byte) (decrypted []byte) {
	key := []byte("ABCDEFGAIJKLMNOP")
	block, _ := aes.NewCipher(key)                              // 分组秘钥
	blockSize := block.BlockSize()                              // 获取秘钥块的长度
	blockMode := cipher.NewCBCDecrypter(block, key[:blockSize]) // 加密模式
	decrypted = make([]byte, len(encrypted))                    // 创建数组
	blockMode.CryptBlocks(decrypted, encrypted)                 // 解密
	decrypted = pkcs5UnPadding(decrypted)                       // 去除补全码
	return decrypted
}
func  (this *Crypto) CreateDsa(str string)(*big.Int,*big.Int,*dsa.PublicKey){
	var param dsa.Parameters
	// L1024N160是一个枚举，根据L1024N160来决定私钥的长度（L N）
	dsa.GenerateParameters(&param, rand.Reader, dsa.L1024N160)
	// 定义私钥的变量
	var privateKey dsa.PrivateKey
	// 设置私钥的参数
	privateKey.Parameters = param
	// 生成密钥对
	dsa.GenerateKey(&privateKey, rand.Reader)
	publicKey := privateKey.PublicKey
	message := []byte(str)

	// 进入签名操作
	r, s, _ := dsa.Sign(rand.Reader, &privateKey, message)
	return r,s,&publicKey
}
func  (this *Crypto) DsaVerify(str string,r *big.Int,s *big.Int, publicKey *dsa.PublicKey) bool{
	return dsa.Verify(publicKey, []byte(str), r, s)
}
func  (this *Crypto)Sha256(str string) string{
	hash:=sha256.Sum256([]byte(str))
	return fmt.Sprintf("%x",hash)
}
func  (this *Crypto)Md5(str string) string{
	hash:=md5.Sum([]byte(str))
	return fmt.Sprintf("%x",hash)
}


func  (this *Crypto) SessionId() string {
	b := make([]byte, 32)
	//rand.Reader是一个全局、共享的密码用强随机数生成器
	rand.Read(b);
	return base64.URLEncoding.EncodeToString(b)//将生成的随机数b编码后返回字符串,该值则作为session ID
}

func(this *Crypto) RSAGenKey(bits int,folderPath string) error {
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {

		// 必须分成两步：先创建文件夹、再修改权限
		os.MkdirAll(folderPath, 0777) //0777也可以os.ModePerm
		os.Chmod(folderPath, 0777)
	}
	/*
		生成私钥
	*/
	//1、使用RSA中的GenerateKey方法生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	//2、通过X509标准将得到的RAS私钥序列化为：ASN.1 的DER编码字符串
	privateStream := x509.MarshalPKCS1PrivateKey(privateKey)
	//3、将私钥字符串设置到pem格式块中
	block1 := pem.Block{
		Type:  "private key",
		Bytes: privateStream,
	}
	//4、通过pem将设置的数据进行编码，并写入磁盘文件
	fPrivate, err := os.Create(folderPath+"/privateKey.pem")
	if err != nil {
		return err
	}
	defer fPrivate.Close()
	err = pem.Encode(fPrivate, &block1)
	if err != nil {
		return err
	}

	/*
		生成公钥
	*/
	publicKey:=privateKey.PublicKey
	publicStream,err:=x509.MarshalPKIXPublicKey(&publicKey)
	//publicStream:=x509.MarshalPKCS1PublicKey(&publicKey)
	block2:=pem.Block{
		Type:"public key",
		Bytes:publicStream,
	}
	fPublic,err:=os.Create(folderPath+"/publicKey.pem")
	if err!=nil {
		return  err
	}
	defer fPublic.Close()
	pem.Encode(fPublic,&block2)
	return nil
}

func(this *Crypto)   EncyptogRSA(src []byte,path string) (res []byte,err error) {
	//1.获取秘钥（从本地磁盘读取）
	f,err:=os.Open(path)
	if err!=nil {
		return
	}
	defer f.Close()
	fileInfo,_:=f.Stat()
	b:=make([]byte,fileInfo.Size())
	f.Read(b)
	// 2、将得到的字符串解码
	block,_:=pem.Decode(b)

	// 使用X509将解码之后的数据 解析出来
	//x509.MarshalPKCS1PublicKey(block):解析之后无法用，所以采用以下方法：ParsePKIXPublicKey
	keyInit,err:=x509.ParsePKIXPublicKey(block.Bytes)  //对应于生成秘钥的x509.MarshalPKIXPublicKey(&publicKey)
	//keyInit1,err:=x509.ParsePKCS1PublicKey(block.Bytes)
	if err!=nil {
		return
	}
	//4.使用公钥加密数据
	pubKey:=keyInit.(*rsa.PublicKey)
	res,err=rsa.EncryptPKCS1v15(rand.Reader,pubKey,src)
	return
}


//对数据进行加密操作
func(this *Crypto)  DecrptogRSA(src []byte,path string)(res []byte,err error)  {
	//1.获取秘钥（从本地磁盘读取）
	f,err:=os.Open(path)
	if err!=nil {
		return
	}
	defer f.Close()
	fileInfo,_:=f.Stat()
	b:=make([]byte,fileInfo.Size())
	f.Read(b)
	block,_:=pem.Decode(b)//解码
	privateKey,err:=x509.ParsePKCS1PrivateKey(block.Bytes)//还原数据
	res,err=rsa.DecryptPKCS1v15(rand.Reader,privateKey,src)
	return
}