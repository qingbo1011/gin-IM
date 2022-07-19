package route

import (
	"gin-IM/api"
	"gin-IM/middleware"
	ws2 "gin-IM/service/ws"

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
		test.GET("/test", api.Test)
	}

	ws := r.Group("/ws")
	//ws.Use(middleware.JWTAuth()) // 方便测试这里就没有使用jwt,该有的鉴权功能已经实现类（正常来说uid应该在token中获取）
	{
		ws.GET("/", ws2.WsHandler)
	}
	return r
}
