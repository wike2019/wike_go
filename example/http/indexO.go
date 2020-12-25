package http

type User struct {
	Name string `form:"name" checking:"required,CheckName" json:"name" gorm:"column:user_name" name:"user_name" vmsg1:"用户名必填"`
	Id string `form:"id" checking:"required" json:"id" gorm:"column:id" name:"id"`
	User_id string `form:"user_id" checking:"required" json:"user_id" gorm:"column:user_id" name:"user_id"`
}

func NewUser() *User {
	return &User{Name:"wike",Id:"11111",User_id:"fdasfasfsafa"}
}

