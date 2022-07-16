package route

import (
	"gin-IM/api"
	"gin-IM/middleware"
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
	agent := r.Group("/api/agent")
	{
		agent.GET("/", func(c *gin.Context) {
			ua := c.GetHeader("User-Agent")
			c.JSON(http.StatusOK, gin.H{
				"User-Agent": ua,
			})
		})
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
	return r
}
