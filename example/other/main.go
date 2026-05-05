package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"github.com/wike2019/wike_app/config"
	"github.com/wike2019/wike_go/v2/core"
	"github.com/wike2019/wike_go/v2/func/ants_service"
	"github.com/wike2019/wike_go/v2/func/ctl"
	"github.com/wike2019/wike_go/v2/func/hash"
	"github.com/wike2019/wike_go/v2/func/os"
	"github.com/wike2019/wike_go/v2/func/retry"
	"github.com/wike2019/wike_go/v2/utils"
)

type router struct {
	ctl.Controller
}

func NewRouter() *router {
	return &router{}
}

func (this *router) healtzh(context *gin.Context) {
	c := this.SetContext(context)
	//重试
	re := retry.NewRetry(this.job)
	re.SetTimes(5)
	re.SetDelay(time.Second * 1)
	err := re.Do()
	fmt.Println(err)
	c.OKWithMsg("hello world")
}
func (this router) Build(r *gin.RouterGroup, GCore *core.GCore) {
	r.GET("/healthz", this.healtzh)
}
func (this router) Path() string {
	return "/"
}
func (this *router) job() error {
	fmt.Println(time.Now().String())
	return nil
}

// 数据库模型
type UserPO struct {
	ID       int
	Name     string
	Password string
	Age      int
	Email    string
}

// 返回给前端的 VO
type UserVO struct {
	ID    int
	Name  string
	Age   int
	Email string
}
type Test struct {
	Name string
}

type Order struct {
	ID     int
	UserID int
	Amount float64
}

func main() {
	//有锁map
	mapstring := utils.NewMap[string]()
	mapstring.Set("111", "2222")
	mapstring.Set("222", "3333")
	fmt.Println(mapstring.Keys())
	fmt.Println(mapstring.Values())
	fmt.Println(mapstring.Get("111"))
	//密码加密与验证
	hashed, _ := utils.PasswordHash("123456")
	isLogin := utils.PasswordVerify("123456", hashed)
	fmt.Println(isLogin)
	//集合去重
	set := utils.NewSet()
	set.Add("111")
	set.Add("222")
	set.Add("333")
	set.Add("222")
	var str []string
	set.List(&str)
	for i, item := range str {
		fmt.Println(i, item)
	}
	set2 := utils.NewSet()
	obj1 := &Test{
		Name: "a1",
	}
	obj2 := &Test{
		Name: "a1",
	}
	obj3 := &Test{
		Name: "a2",
	}
	set2.Add(obj1)
	set2.Add(obj2)
	set2.Add(obj3)
	var obj []*Test
	set2.List(&obj)
	for i, item := range obj {
		fmt.Println(i, item)
	}
	//生成随机数据
	fmt.Println(utils.RandomString(10))
	fmt.Println(utils.GetRandomNum(1, 99))
	//并发控制
	// 1. 创建池，限制最多 10 个 goroutine 并发
	pool, err := ants_service.NewPool(10)

	// 2. 设置总任务数（内部调用 WaitGroup.Add）
	pool.SetTotal(100)

	// 3. 提交任务
	for i := 0; i < 100; i++ {
		pool.Submit(func() error {
			// 业务逻辑
			return nil
		})
	}

	// 4. 等待所有任务完成
	pool.Wait()

	// 5. 检查结果
	fmt.Println(pool.Ok, pool.Fail, pool.Error())

	//一致性hash

	h := hash.New()

	// 添加节点（比如服务器地址）
	h.Add("192.168.1.1:8080")
	h.Add("192.168.1.2:8080")
	h.Add("192.168.1.3:8080")

	// 根据 key 获取对应的节点
	node, err := h.Get("user_123")
	if err != nil {
		panic(err)
	}
	fmt.Println(node) // 输出某个节点地址
	//系统参数获取
	fmt.Println(os.InitOS())       //获取 Go 运行时信息：系统类型、CPU 核数、编译器、Go 版本、goroutine 数
	fmt.Println(os.InitCPU())      //获取每个 CPU 核心的使用率百分比和物理核心数
	fmt.Println(os.InitRAM())      //获取内存已用 MB、总量 MB、使用百分比
	fmt.Println(os.InitDisk("./")) //获取指定挂载点的磁盘用量（MB/GB）和使用百分比

	//友好打印时间
	fmt.Println(utils.ParseDuration("1d5h20m"))

	//基于redis实现的队列
	client := redis.NewClient(&redis.Options{
		Addr: "192.168.3.2:6379",
	})

	queue := utils.NewRedisQueue(client)

	// ========== 生产者：推送消息 ==========
	order := Order{ID: 1001, UserID: 42, Amount: 99.9}
	err = queue.Push("order_queue", order)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("消息已推送")

	// ========== 消费者：阻塞等待并消费消息 ==========
	var received Order
	err = queue.Pop("order_queue", &received)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("收到订单: ID=%d, 金额=%.1f\n", received.ID, received.Amount)

	lock := utils.NewRedisLock(client)

	// 加锁，key="order_lock"，TTL=10秒
	lockData, err := lock.Lock("order_lock", 10, nil)
	if err != nil {
		log.Fatal("加锁失败:", err)
	}
	fmt.Println("加锁成功") // ===== 执行需要互斥的业务逻辑 =====
	fmt.Println("正在处理订单...")
	// 比如：扣库存、创建订单等

	// 释放锁
	err = lockData.Release(context.Background())
	if err != nil {
		log.Println("释放锁失败:", err)
	}
	fmt.Println("锁已释放")

	// 1. 计算 MD5 值
	data := []byte("hello world")
	md5Str := utils.MD5V(data)
	fmt.Println(md5Str) // 输出: 5eb63bbbe01eeed093cb22bb8f5acdc3

	// 2. 校验分片完整性（文件分片上传场景）
	chunk := []byte("这是文件的第一个分片内容")
	expectedMd5 := utils.MD5V(chunk) // 上传前客户端计算的 MD5

	// 服务端收到分片后校验
	ok := utils.CheckMd5(chunk, expectedMd5)
	fmt.Println("分片校验:", ok) // true

	// 如果数据被篡改或传输损坏
	ok = utils.CheckMd5([]byte("损坏的数据"), expectedMd5)
	fmt.Println("分片校验:", ok) // false
	// 从数据库查出来的完整数据
	user := UserPO{
		ID:       1,
		Name:     "张三",
		Password: "secret123",
		Age:      25,
		Email:    "zhangsan@example.com",
	}

	// 拷贝到 VO，自动忽略 Password 字段（VO 中没有）
	var vo UserVO
	err = utils.CopyProperties(&vo, user)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", vo)
	// 输出: {ID:1 Name:张三 Age:25 Email:zhangsan@example.com}
	//一个最简单的例子
	g := core.God()
	g.Provide(config.Config).MountWithEmpty(NewRouter).Run()
}
