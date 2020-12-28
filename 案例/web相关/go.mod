module demo

go 1.15

require (
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis/v8 v8.4.4
	github.com/go-sql-driver/mysql v1.5.0
	github.com/jinzhu/gorm v1.9.16
	github.com/shenyisyn/goft-expr v0.3.0
	github.com/shenyisyn/goft-ioc v0.5.4
	github.com/wike2019/wike_go v1.0.7
)

replace google.golang.org/grpc v1.27.0 => google.golang.org/grpc v1.26.0