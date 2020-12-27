package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/wike2019/wike_go/example/services"
	"github.com/wike2019/wike_go/src/Grpc"
	"github.com/wike2019/wike_go/src/Web"
	"github.com/wike2019/wike_go/src/core/Cache"
	"github.com/wike2019/wike_go/src/core/Etcd"
	"github.com/wike2019/wike_go/src/core/Redis"
	"github.com/wike2019/wike_go/src/core/sql"
	"github.com/wike2019/wike_go/src/util/LoadBalance"
	"io"
	"log"
	"time"
)


type IndexController struct {
	Db *gorm.DB                          `inject:"-"`
	RedisSet *Redis.RedisStringOperation `inject:"-"`
	Etcdctl *Etcd.EtcdCtl                `inject:"-"`
}
func NewIndexController() *IndexController {
	return &IndexController{}
}
func(this *IndexController) Index(ctx *gin.Context) string   {
	this.RedisSet.Set("keys","我是第一个", Redis.WithExpire(time.Second*10)).Unwrap()
	this.RedisSet.Set("keys","我是第一个", Redis.WithExpire(time.Second*10), Redis.WithNX()).Unwrap()
	this.RedisSet.Set("keys","我是第一个", Redis.WithExpire(time.Second*10), Redis.WithXX()).Unwrap()
	fmt.Println()
    return "111"
}
func(this *IndexController) Index2(ctx *gin.Context) string   {

	fmt.Println(this.Etcdctl.EtcdClient)
	return "111"
}
func(this *IndexController) Name () string   {
	return "IndexController"
}
func(this *IndexController) Wike (ctx *gin.Context) string   {
	return this.RedisSet.Get("keys").Unwrap().(string)
}
func(this *IndexController) A (ctx *gin.Context) string   {

	return "1111111"
}
func(this *IndexController) A2 (ctx *gin.Context) Web.Json {
    m,_:=this.Etcdctl.LoadService("wike3")
    fmt.Println(m)
    info,_:=this.Etcdctl.Seletor(m,LoadBalance.RoundRobinByWeight,"192.168.127.1")
	return info
}
func(this *IndexController) B (ctx *gin.Context) Web.Json {

	return gin.H{"resut":"test"}
}
func(this *IndexController) C (ctx *gin.Context) Web.Json {
	Web.Throw("aaaaa",500,ctx)
	return gin.H{"resut":"test"}
}
func(this *IndexController) D (ctx *gin.Context) Web.Json {

	return gin.H{"resut":"test"}
}
func(this *IndexController) Time () {

	fmt.Println(time.Now())
}
func(this *IndexController) E (ctx *gin.Context) Web.Json {
	//
	req,_:=ctx.Get("_req")
	fmt.Println(req)
	u:=req.(*User)
	u.Id+="wiiadwadasd"
	return gin.H{"resut":u}
}
func(this *IndexController) F (ctx *gin.Context) Web.Json {
	u:=&User{}
	fmt.Println(this.Db)
	Web.Error(this.Db.Table("users").Where("user_id=?",2).Find(u).Error)
	return gin.H{"resut":u}
}
func(this *IndexController) F2 (ctx *gin.Context) sql.SimpleQuery   {
	return  "select * from users"
}
func(this *IndexController) F23 (ctx *gin.Context) sql.Query   {
	return  sql.SimpleQuery("select * from users where user_id = ?").WithArgs(1)
}
func(this *IndexController) F234 (ctx *gin.Context) sql.Query   {
	return  sql.SimpleQuery("select * from users where user_id = ?").WithArgs(1).WithFirst()
}
func(this *IndexController) F2345 (ctx *gin.Context) sql.Query   {
	return  sql.SimpleQuery("select * from users where user_id = ?").WithArgs(1).WithFirst().WithKey("result")
}
func(this *IndexController) F23456 (ctx *gin.Context) Web.Json {
	//u:=&User{}
	ret:=  sql.SimpleQuery("select * from users where user_id = ?").WithArgs(1).WithFirst().Get()
	m := ret.(map[string]interface{})
	m["additive"]="111111111"
	return  m
}
func(this *IndexController) F234567 (ctx *gin.Context) Web.Json {
	// 1、从对象池 获取新闻缓存 对象
	newsCache:= Cache.RedisCache()
	defer Cache.ReleaseRedisCache(newsCache)
	id:=ctx.Query("id")
	newsCache.DBGetter= NewUserGetter(id) //一旦缓存没有，则需要从哪去取
	// 3、取缓存输出（一旦没有，上面的DBGetter会被调用)
	usr:=&User{}
	newsCache.GetCacheForObject("user"+id,usr)
	return  usr
}
func(this *IndexController) F2345679 (ctx *gin.Context) string {
	// 1、从对象池 获取新闻缓存 对象
	newsCache:= Cache.RedisCache()
	defer Cache.ReleaseRedisCache(newsCache)
	id:=ctx.Query("id")
	newsCache.DBGetter= NewUserGetter2(id) //一旦缓存没有，则需要从哪去取
	// 3、取缓存输出（一旦没有，上面的DBGetter会被调用)
	var usr string
	newsCache.GetCacheForObject("news"+id,&usr)
	return  usr
}
func(this *IndexController) Grpc (ctx *gin.Context) Web.Json {
	// 1、从对象池 获取新闻缓存 对象
	client:=Grpc.NewClient(Grpc.KeyPath("./keys"),Grpc.Host("wike.com"),Grpc.Ip("192.168.3.3:8081"))
	defer client.Close()
	svc:= services.NewUserServiceClient(client)
	users:=make([]*services.UserInfo,0)
	users=append(users,&services.UserInfo{UserId: 1})
	users=append(users,&services.UserInfo{UserId: 2})
	rsp,_:=svc.GetUserScore(ctx,&services.UserScoreRequest{Users: users})
	return  rsp
}

func(this *IndexController) Grpc5(ctx *gin.Context) Web.Json {
	// 1、从对象池 获取新闻缓存 对象
	client:=Grpc.NewClient(Grpc.KeyPath("./keys"),Grpc.WithEtcd("wike3",LoadBalance.RoundRobinByWeight,ctx.ClientIP()))
	defer client.Close()
	svc:= services.NewUserServiceClient(client)
	users:=make([]*services.UserInfo,0)
	users=append(users,&services.UserInfo{UserId: 1})
	users=append(users,&services.UserInfo{UserId: 2})
	rsp,_:=svc.GetUserScore(ctx,&services.UserScoreRequest{Users: users})
	return  rsp
}
func(this *IndexController) Grpc2 (ctx *gin.Context) string {
	// 1、从对象池 获取新闻缓存 对象
	client:=Grpc.NewClient(Grpc.KeyPath("./keys"),Grpc.Host("wike.com"),Grpc.Ip("192.168.3.3:8081"))
	defer client.Close()
	svc:= services.NewUserServiceClient(client)
	users:=make([]*services.UserInfo,0)
	users=append(users,&services.UserInfo{UserId: 1})
	users=append(users,&services.UserInfo{UserId: 2})
	stream,_:=svc.GetUserScoreByServerStream(ctx,&services.UserScoreRequest{Users: users})
	for{
		res,err:=stream.Recv()
		if err==io.EOF{
			break
		}
		fmt.Println(res.Users)

	}
	return  "结束操作"
}

func(this *IndexController) Grpc3 (ctx *gin.Context) Web.Json {
	// 1、从对象池 获取新闻缓存 对象
	client:=Grpc.NewClient(Grpc.KeyPath("./keys"),Grpc.Host("wike.com"),Grpc.Ip("192.168.3.3:8081"))
	defer client.Close()
	svc:= services.NewUserServiceClient(client)
	users:=make([]*services.UserInfo,0)
	users=append(users,&services.UserInfo{UserId: 1})
	users=append(users,&services.UserInfo{UserId: 2})
	stream,_:=svc.GetUserScoreByClientStream(ctx)
	//可以多次发送数据
	stream.Send(&services.UserScoreRequest{Users: users})
	//结束操作
	res,_:=stream.CloseAndRecv()
	return  res
}
func(this *IndexController) Grpc4 (ctx *gin.Context) string {
	// 1、从对象池 获取新闻缓存 对象
	client:=Grpc.NewClient(Grpc.KeyPath("./keys"),Grpc.Host("wike.com"),Grpc.Ip("192.168.3.3:8081"))
	defer client.Close()
	svc:= services.NewUserServiceClient(client)
	for j:=1;j<=3;j++{
		users:=make([]*services.UserInfo,0)
		users=append(users,&services.UserInfo{UserId: 1})
		users=append(users,&services.UserInfo{UserId: 2})
		stream,_:=svc.GetUserScoreByTWS(ctx)
		stream.Send(&services.UserScoreRequest{Users: users})
		res,err:=stream.Recv()
		if err==io.EOF{
			break
		}
		if err!=nil{
			log.Println(err)
		}
		fmt.Println(res.Users)
	}

	return  "结束操作"
}
func(this *IndexController) Build(goft *Web.Goft){
	goft.Handle("GET","/",this.Index).
		Handle("GET","/wike",this.Wike).
	Handle("GET","/A",this.A).
	Handle("GET","/B",this.B).
	Handle("GET","/C",this.C).
		HandleWithFairing("GET","/D",this.D, NewTest()).
		Handle("GET","/F",this.F).
		Handle("GET","/F2",this.F2).
		Handle("GET","/F23",this.F23).
		Handle("GET","/F234",this.F234).
		Handle("GET","/F2345",this.F2345).
		Handle("GET","/F23456",this.F23456).
		Handle("GET","/F234567",this.F234567).
		Handle("GET","/F2345679",this.F2345679).
		Handle("GET","/Index2",this.Index2).
		Handle("GET","/A2",this.A2).
		Handle("GET","/grpc",this.Grpc).
		Handle("GET","/grpc2",this.Grpc2).
		Handle("GET","/grpc3",this.Grpc3).
		Handle("GET","/grpc4",this.Grpc4).
		Handle("GET","/grpc5",this.Grpc5).

	HandleWithFairing("GET","/E",this.E, NewTestIn())
}

