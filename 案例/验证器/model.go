package main


type User struct {
	Name string `form:"name" binding:"required,CheckName" json:"name" gorm:"column:user_name" name:"user_name" vmsg:"required=用户名必填,CheckName=我是提示信息"`
}