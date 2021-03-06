package Grpc

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/wike2019/wike_go/src/Result"
	"github.com/wike2019/wike_go/src/Core/Etcd"
	"github.com/wike2019/wike_go/src/Core/Bean"
	"github.com/wike2019/wike_go/src/Util/LoadBalance"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
)

type  Grpc struct {
   KeyPath string
   ChainStreamServer []grpc.StreamServerInterceptor
   ChainServer []grpc.UnaryServerInterceptor
   Ip string
   Host string
   Selector int
   ServerName string
   ClientIp string
   Etcdctl  *Etcd.EtcdCtl
}

func init() {
	Bean.New().Beans(new(Grpc))
}

func (this *Grpc) Name()string  {
	return  "Grpc"
}

func NewServer(fs ...GrpcAttrFunc) *grpc.Server {
	u:= new(Grpc)
	u.ChainServer=make([]grpc.UnaryServerInterceptor,0)
	u.KeyPath="./keys"
	u.ChainStreamServer=make([]grpc.StreamServerInterceptor,0)
	GrpcAttrFuncs(fs).apply(u)
	u.ChainStreamServer=append(u.ChainStreamServer,grpc_recovery.StreamServerInterceptor())
	u.ChainServer=append(u.ChainServer,RecoveryInterceptor)

	c:=u.GetServercert(u.KeyPath)
	opts := []grpc.ServerOption{
		grpc.Creds(c),
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			u.ChainStreamServer...,
			)),
		grpc.UnaryInterceptor(grpc_middleware.ChainUnaryServer(
			u.ChainServer...,
			)),
	}
	server := grpc.NewServer(opts...)
	return  server
}
func NewClient(fs ...GrpcAttrFunc) *grpc.ClientConn {
	u:= new(Grpc)
	u.KeyPath="./keys"
	u.Host="localhost"
	u.Etcdctl=Etcd.EtcdCache()
	defer  Etcd.ReleaseEtcdCache(u.Etcdctl)
	GrpcAttrFuncs(fs).apply(u)
	if u.ServerName !=""{
		m:=u.Etcdctl.LoadService(u.ServerName)
		info := Result.Result(u.Etcdctl.Seletor(m, LoadBalance.RoundRobinByWeight, u.ClientIp)).Unwrap()[0].(LoadBalance.NodeBalance)
		u.Host=info.GetNode().(*Etcd.ServiceInfo).ServiceHost
		u.Ip=info.GetNode().(*Etcd.ServiceInfo).ServiceAddr
		fmt.Println(u)
	}
	conn,err:=grpc.Dial(u.Ip,grpc.WithTransportCredentials(u.GetClientcert(u.KeyPath,u.Host)))
		if err!=nil{
		log.Fatal(err)
	}
	return  conn
}
func  (this *Grpc)  GetServercert(path string) credentials.TransportCredentials {
	cert,_:=tls.LoadX509KeyPair(path+"/server.pem",path+"/server.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile(path+"/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds:=credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},//服务端证书
		ClientAuth:   tls.RequireAndVerifyClientCert,
		ClientCAs:    certPool,
	})
	return creds
}
func (this *Grpc)  GetClientcert(path string,host string) credentials.TransportCredentials {
	cert,_:=tls.LoadX509KeyPair(path+"/client.pem",path+"/client.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile(path+"/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds:=credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},//客户端证书
		ServerName: host,
		RootCAs:      certPool,
	})
	return  creds
}