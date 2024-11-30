package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wike2019/wike_go/pkg/core"
	controller "github.com/wike2019/wike_go/pkg/service/http"
	"github.com/wike2019/wike_go/pkg/service/memorylog"
	"log"
	"time"
)

type router struct {
	controller.Controller
}
type Invoke struct {
}
type Query struct {
	Id   int    `form:"id" required:"true" desc:"主键"`
	Name string `form:"name" desc:"姓名"`
}
type Body struct {
	Id  int `json:"id" required:"true" desc:"主键"`
	Age int `json:"age" desc:"年龄"`
}

type Header struct {
	Token string `header:"token" desc:"验证"`
}

func NewInvoke() *Invoke {
	return &Invoke{}
}
func (this *Invoke) job() error {
	fmt.Println(time.Now().String())
	return nil
}
func NewRouter() *router {
	return &router{}
}

func (this *router) healtzh(context *gin.Context) {
	c := this.SetContext(context)
	c.Success("修改成功", nil)
}
func (this *router) log(context *gin.Context) {
	c := this.SetContext(context)
	c.Success("修改成功", memorylog.LogInfo.All())
}
func (this *router) job() error {
	fmt.Println(time.Now().String())
	return nil
}
func (this *router) Name() string {
	return "测试路由"
}
func (this *router) stop() error {
	fmt.Println("这里做局部清理")
	return nil
}

type Game struct {
	ID    uint64 `gorm:"primaryKey;column:id;comment:主键" json:"id"`
	Game  string `gorm:"column:game;type:varchar(255);comment:游戏" json:"game"`
	Type  string `gorm:"column:type;type:varchar(255);comment:分类" json:"type"`
	Order int    `gorm:"column:order;type:int;comment:排序" json:"order"`
	Show  int    `gorm:"column:show;type:int;comment:冗余字段" json:"show"` //当初用来区分不同账号登入，显示游戏的，后来没用到
}

func (this *router) Build(r *gin.RouterGroup, GCore *core.GCore) {
	GCore.Stop(this.stop)
	this.SetDoc(Query{}, Body{}, Header{}, controller.PageDoc[Game]())
	GCore.GetWithRbac(r, this, "healthz1", this.healtzh, "游戏列表1")
	this.SetDoc(Query{}, nil, Header{}, controller.PageDoc[Game]())
	GCore.GetWithRbac(r, this, "healthz2", this.healtzh, "游戏列表2")
	this.SetDoc(Query{}, nil, nil, controller.PageDoc[Game]())
	GCore.GetWithRbac(r, this, "healthz3", this.healtzh, "游戏列表3")
	this.SetDoc(nil, nil, nil, controller.PageDoc[Game]())
	GCore.GetWithRbac(r, this, "healthz4", this.healtzh, "游戏列表4")
	this.SetDoc(Query{}, Body{}, Header{}, nil)
	GCore.GetWithRbac(r, this, "healthz5", this.healtzh, "游戏列表5")
	this.SetDoc(Query{}, nil, Header{}, controller.PageDoc[Game]())
	GCore.GetWithRbac(r, this, "healthz6", this.healtzh, "游戏列表6")
	this.SetDoc(nil, nil, nil, nil)
	GCore.GetWithRbac(r, this, "healthz7", this.healtzh, "游戏列表7")

	this.SetDoc(nil, nil, nil, controller.HttpDoc[Game]{})
	GCore.GetWithRbac(r, this, "healthz8", this.healtzh, "游戏列表8")

	this.SetDoc(nil, Body{}, nil, controller.HttpDoc[Game]{})
	GCore.GetWithRbac(r, this, "healthz9", this.healtzh, "游戏列表9")

}
func (this router) Path() string {
	return "/"
}

//	type Job struct {
//		i2 int
//	}
//
//	func (this *Job) Job() error {
//		if this.i2 == 3 {
//			return errors.New("错误")
//		}
//		time.Sleep(time.Duration(this.i2) * time.Second)
//		fmt.Println("time: ", this.i2, time.Now().Format("2006-01-02 15:04:05"))
//		return nil
//	}
func MyViper(core *core.GCore) *viper.Viper {
	//fmt.Println(core)
	viper.SetDefault("port", "8888")
	viper.SetDefault("logPath", "./logs/app.log")
	viper.SetDefault("development", true)
	viper.SetConfigFile("config.yaml") // 指定配置文件路径
	viper.SetConfigName("config")      // 配置文件名称(无扩展名)
	viper.SetConfigType("yaml")        // 如果配置文件的名称中没有扩展名，则需要配置此项
	viper.AddConfigPath(".")           // 还可以在工作目录中查找配置
	err := viper.ReadInConfig()        // 查找并读取配置文件
	if err != nil {                    // 处理读取配置文件的错误
		log.Fatalf("Fatal error config file: %s \n", err.Error())
	}
	//业务参数
	viper.SetDefault("Timeout", 3000)
	core.Stop(func() error {
		fmt.Println("这里是viper清理")
		return nil
	})
	return viper.GetViper()
}
func main() {
	//
	//	pool, _ := ants_service.NewPool(2)
	//	defer pool.Release()
	//	pool.SetTotal(6)
	//	for i := 0; i < 6; i++ {
	//		func(k int) {
	//			j := &Job{i2: k}
	//			err := pool.Submit(j.Job)
	//			if err != nil {
	//				fmt.Println(err)
	//			}
	//		}(i)
	//
	//	}
	//	//time.Sleep(1 * time.Second)
	//
	//	fmt.Println("finish-1")
	//	pool.Wait() // 等待所有任务完成
	//	fmt.Println(pool.Ok, pool.Fail, pool.Error())
	//	fmt.Println("finish")
	//	return
	//
	//	//一个最简单的例子
	g := core.God()
	g.Stop(func() error {
		fmt.Println("这里做全局清理")
		return nil
	})

	g.Default().GlobalUse(core.CORSMiddleware()) //选择redis作为缓存服务的存储
	g.Provide(MyViper).Provide(NewInvoke).MountWithEmpty(NewRouter).Invokes(func(r *Invoke) {
		go r.job() //这里不能阻塞 所以最好用 go xxx
	}).Run()
}
