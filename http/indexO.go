package http

type User struct {
	Name string `form:"name" binding:"required" json:"name" gorm:"column:user_name" name:"user_name"`
	Id string `form:"id" binding:"required" json:"id" gorm:"column:id" name:"id"`
	User_id string `form:"user_id" binding:"required" json:"user_id" gorm:"column:user_id" name:"user_id"`
}

func NewUser() *User {
	return &User{Name:"wike",Id:"11111",User_id:"fdasfasfsafa"}
}