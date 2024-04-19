package validator

import (
	"github.com/gin-gonic/gin"
	"testing"
)

type Login struct {
	Username string `form:"username" json:"username" binding:"required,email"`
	Password string `form:"password" json:"password" binding:"required"`
}

func TestMap(t *testing.T) {
	r := gin.New()
	r.POST("/login", func(c *gin.Context) {
		var login Login
		if err := c.ShouldBind(&login); err != nil {
			c.JSON(200, gin.H{"message": TranslateJson(err, &login)})
			return
		}

		c.JSON(200, gin.H{"message": "Login success!"})
	})
	r.Run()
}
