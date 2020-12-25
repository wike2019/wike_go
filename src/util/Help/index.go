package Help

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net"
	"os"
	"reflect"
	"time"
)

func ExistDir(dirname string) bool {
	fi, err := os.Stat(dirname)
	return (err == nil || os.IsExist(err)) && fi.IsDir()
}

func Map2Struct(m map[string]interface{},u interface{})  {
	v:=reflect.ValueOf(u)
	if v.Kind()==reflect.Ptr{
		v=v.Elem()
		if v.Kind()!=reflect.Struct{
			panic("must struct")
		}
		findFromMap:= func(key string,nameTag string ) interface {}{
			for k,v:=range m{
				if k==key || k==nameTag {
					return v
				}
			}
			return nil
		}
		for i:=0;i<v.NumField();i++{
			get_value:=findFromMap(v.Type().Field(i).Name,v.Type().Field(i).Tag.Get("name"))
			if get_value!=nil && reflect.ValueOf(get_value).Kind()==v.Field(i).Kind(){
				v.Field(i).Set(reflect.ValueOf(get_value))
			}
		}
	}else{
		panic("must ptr")
	}
}


func CreateCA(hosts []string, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.MkdirAll(path, 0644) //0777也可以os.ModePerm
	}
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(1653),
		Subject: pkix.Name{
			Organization: []string{"Acme Co"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(10, 0, 0),
		SubjectKeyId:          []byte{1, 2, 3, 4, 5},
		BasicConstraintsValid: true,
		IsCA:        true,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:    x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	privCa, _ := rsa.GenerateKey(rand.Reader, 1024)
	createCertificateFile(path+"/"+"ca", ca, privCa, ca, nil)
	server := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"SERVER"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 6},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			server.IPAddresses = append(server.IPAddresses, ip)
		} else {
			server.DNSNames = append(server.DNSNames, h)
		}
	}
	privSer, _ := rsa.GenerateKey(rand.Reader, 1024)
	createCertificateFile(path+"/"+"server", server, privSer, ca, privCa)
	client := &x509.Certificate{
		SerialNumber: big.NewInt(1658),
		Subject: pkix.Name{
			Organization: []string{"CLIENT"},
		},
		NotBefore:    time.Now(),
		NotAfter:     time.Now().AddDate(10, 0, 0),
		SubjectKeyId: []byte{1, 2, 3, 4, 7},
		ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:     x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
	}
	privCli, _ := rsa.GenerateKey(rand.Reader, 1024)
	createCertificateFile(path+"/"+"client", client, privCli, ca, privCa)
}

func createCertificateFile(name string, cert *x509.Certificate, key *rsa.PrivateKey, caCert *x509.Certificate, caKey *rsa.PrivateKey) {
	priv := key
	pub := &priv.PublicKey
	privPm := priv
	if caKey != nil {
		privPm = caKey
	}
	ca_b, err := x509.CreateCertificate(rand.Reader, cert, caCert, pub, privPm)
	if err != nil {
		log.Println("create failed", err)
		return
	}
	ca_f := name + ".pem"
	log.Println("write to pem", ca_f)
	var certificate = &pem.Block{Type: "CERTIFICATE",
		Headers: map[string]string{},
		Bytes:   ca_b}
	ca_b64 := pem.EncodeToMemory(certificate)
	ioutil.WriteFile(ca_f, ca_b64, 0777)

	priv_f := name + ".key"
	priv_b := x509.MarshalPKCS1PrivateKey(priv)
	log.Println("write to key", priv_f)
	ioutil.WriteFile(priv_f, priv_b, 0777)
	var privateKey = &pem.Block{Type: "PRIVATE KEY",
		Headers: map[string]string{},
		Bytes:   priv_b}
	priv_b64 := pem.EncodeToMemory(privateKey)
	ioutil.WriteFile(priv_f, priv_b64, 0777)
}

func GenRSAPubAndPri(bits int,filepath string ) error {

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		// 必须分成两步：先创建文件夹、再修改权限
		os.MkdirAll(filepath, 0644) //0777也可以os.ModePerm
	}
	// 生成私钥文件
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	priBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}

	err= ioutil.WriteFile(filepath+"/private.pem",pem.EncodeToMemory(priBlock), 0644)
	if err!=nil{
		return err
	}
	fmt.Println("=======私钥文件创建成功========")
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	derPkix, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	publicBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}

	err= ioutil.WriteFile(filepath+"/public.pem",pem.EncodeToMemory(publicBlock), 0644)
	if err!=nil{
		return err
	}
	fmt.Println("=======公钥文件创建成功=========")

	return nil
}
