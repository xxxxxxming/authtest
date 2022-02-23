package routers

import (
	services "main/apps/test"

	"github.com/gin-gonic/gin"
)

var s services.TestServices

func testGroupRouter(c *gin.Engine) *gin.RouterGroup {
	return c.Group("test")
}

func TestRouter(c *gin.Engine) {
	group := testGroupRouter(c)
	group.GET("/:id", s.Test1)
	group.GET("/t1/:id", s.Test2)
	group.GET("/t2", s.Test3)
	group.GET("/t/:id1", s.Test4)
	group.GET("/t2/:id1/:id2", s.Test5)
	group.GET("/", s.Test6)

	group.POST("/:id", s.Test7)
	group.POST("/t1/:id", s.Test8)
	group.POST("/t2", s.Test9)
	group.POST("/t/:id1", s.Test10)
	group.POST("/t2/:id1/:id2", s.Test11)
	group.POST("/", s.Test12)

	group.PATCH("/:id", s.Test13)
	group.PATCH("/t1/:id", s.Test14)
	group.PATCH("/t2", s.Test15)
	group.PATCH("/t/:id1", s.Test16)
	group.PATCH("/t2/:id1/:id2", s.Test17)
	group.PATCH("/", s.Test18)

	group.DELETE("/:id", s.Test19)
	group.DELETE("/t1/:id", s.Test20)
	group.DELETE("/t2", s.Test21)
	group.DELETE("/t/:id1", s.Test22)
	group.DELETE("/t2/:id1/:id2", s.Test23)
	group.DELETE("/", s.Test24)
}
