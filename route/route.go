package route

import (
	"gin-IM/api"
	"gin-IM/middleware"
	ws2 "gin-IM/service/ws"
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRoute() *gin.Engine {
	// gin.Default()和gin.New()的区别：gin.Default()默认使用了2个中间件Logger(), Recovery()
	r := gin.Default()
	user := r.Group("/api/user")
	{
		user.POST("/register", api.UserRegister)
		user.POST("/login", api.UserLogin)
	}
	test := r.Group("api/test")
	test.Use(middleware.JWTAuth())
	{
		test.GET("/test", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"msg": "test",
			})
		})
	}

	ws := r.Group("/ws")
	{
		ws.GET("/", ws2.WsHandler)
	}
	return r
}
