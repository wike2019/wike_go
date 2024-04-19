# wike-go 是一款web框架，基于成熟的gin框架进行二次封装，提供了一些常用的功能，让开发者更专注于业务开发。

## 特性 

  1 基于gin完美兼容，支持更优化的api

  2 依赖注入，基于fx的二次封装，更简单的使用

  3 丰富的工具函数

## 申明
  框架刚成立，需要很多测试，和完善，功能也会慢慢完善，欢迎大家提出宝贵的意见，一起完善这个框架

## APi说明
### coreHttp.God 初始化核心模块支持连式调用
### coreHttp.Mount  挂载路由，里面可以传入多个路由的构造函数 (构造函数返回的参数必须实现 Controller 接口)
demo:

    type Router struct {
      log *zaplog.Logger  //依赖注入的日志
      *controller.Controller //基础控制器提供了一些常用的方法
    }
    
    func (r *Router) Get(ctx *gin.Context) {
        r.ParsePage(ctx)  //解析分页参数
        fmt.Println(r.Page)
        fmt.Println(r.PageSize)
        r.log.Info("这个是依赖注入的日志")
        r.List(ctx, []int{1, 2, 3, 4, 5}, 5)
    }
    func (r *Router) Post(ctx *gin.Context) {
        r.ParsePage(ctx)  //解析分页参数
        fmt.Println(r.Page)
        fmt.Println(r.PageSize)
        r.log.Info("这个是依赖注入的日志")
        r.List(ctx, []int{1, 2, 3, 4, 5}, 5)
    }
    func (r *Router) Init() {
         //必须实现的方法
         //这个函数是在初始化的时候调用的，可以在这里给属性赋值
         // 默认给基础控制器赋值,你可以在这里加入需要的逻辑
         r.Controller = controller.BaseController
    }
    func (r *Router) Path() string {
        //必须实现的方法
        //返回gin.group的path,这个path会拼接在Build的path上
         return "/api/v2"
    }
    func (this *Router) Build(r *gin.RouterGroup) {
        //必须实现的方法
        //路由注册
        r.GET("/hello", this.Get)
        r.POST("/hello2", this.Post)
    }
    func RouterInit(log *zaplog.Logger) *Router {
        //初始化函数 这个函数会在初始化的时候调用 log 是依赖注入自动添加的
        r := &Router{log: log}
        r.Init() //给必要的属性赋值
        return r
    }

### coreHttp.Provide  注册依赖注册函数，这个函数会在初始化的时候调用，可以在这里注册依赖
demo:

    Provide(coreHttp.InitSvc，coreHttp.InitSvc2)

    type Svc struct {
        name string
    }
    func InitSvc() *Svc {
        return &Svc{
          name: "111123123123",
        } //必须返回指针
    }
    type Svc2 struct {
         name string
    }
    func InitSvc2() *Svc2 {
         return &Svc2{
             name: "111123123123",
         } //必须返回指针
    }

### coreHttp.Config  注册依赖注册函数，这个函数会在初始化的时候调用，和Provide 的区别在于 Provide 是传入构造函数,Config是传入结构体
#### 传入的结构体必须是指针,并且里面的方法必须返回指针，我们遍历所有方法把返回的指针都注册到依赖中

demo

    Config(&coreHttp.Conf{}) //调用 注册svc3 和 svc2

    type Conf struct {
    }
    
    func (this *Conf) M1() *Svc3 {
         return &Svc3{}
    }
    func (this *Conf) M2() *Svc2 {
         return &Svc2{}
    }
### coreHttp.GlobalUse
    注册全局中间件 
### coreHttp.Use
    注册路由中间件，这个中间件只会在当前路由注册时匹配到控制器path路由一致的时，进行注册

### 控制器build的函数的详细说明

    func (this *Router2) Build(r *gin.RouterGroup) {
        r.GET("/ping", this.Get) // 只有全局中间件，和路由中间件
        r.GET("/wike2", TestMiddleware, this.Get) //方法级中间件
        r.Use(func(ctx *gin.Context) {
             fmt.Println("进入了中间件*****************")
        })  //后续的路由都会进入这个中间件
        r.GET("/wike", this.Get)
        r.Use(func(ctx *gin.Context) {
          fmt.Println("第二个*****************")
        })
    }

### 日志
  基于zap库二次封装，加入了 trace_id
#### 特殊方法
    GetLoggerWithGinCtx(ctx *gin.Context) *zaplog.Logger //获取gin的上下文的日志
    GetLoggerWithTraceID(ctx *context.Context) *zap.SugaredLogger    //获取trace_id的日志

调用这2个方法会自动添加trace_id的日志，如果没有trace_id会自动添加一个trace_id

#### 内存日志，会保存最近2000条日志，只会收集info，和warn级别日志 
    templog.LogInfo.Show() //打印日志
    templog.LogInfo.All() //返回日志列表

## 工具函数
    CopyProperties(src interface{}, dest interface{}) //结构体属性拷贝 copy.go
    MapSync //map并发安全 map.go
    Decimal //精度计算（） number.go
    PasswordHash //密码加密 password.go
    PasswordVerify //密码验证 password.go
    RandomString 生成随机长度的字符串 random.go
    GetRandomNum 生成随机数字(包括边界) random.go

### 待完善

1. 缓存
2. http客户端
3. 定时任务
4. 消息队列
5. orm
6. redis
7. 服务熔断
8. 服务限流
9. 负载均衡
10. 更多...