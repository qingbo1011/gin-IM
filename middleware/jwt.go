package middleware

import (
	"gin-IM/pkg/util"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// JWTAuth 定义一个JWTAuth的中间件
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 通过http header中的token解析来认证
		token := c.GetHeader("token")
		if token == "" {
			c.JSON(http.StatusForbidden, gin.H{
				"status": http.StatusForbidden,
				"msg":    "请求未携带token，无权访问！",
			})
			c.Abort() // Abort(): 在被调用的函数中阻止后续中间件的执行(http://www.codebaoku.com/gin/gin-abort.html)
			return
		}
		// 解析token中包含的相关信息（有效载荷）
		claims, err := util.ParseToken(token)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"status": http.StatusForbidden,
				"msg":    "token解析失败！",
				"error":  err.Error(),
			})
			c.Abort()
			return
		}
		if time.Now().Unix() > claims.ExpiresAt { // token过期了
			c.JSON(http.StatusForbidden, gin.H{
				"status": http.StatusForbidden,
				"msg":    "token已过期！",
			})
			c.Abort()
			return
		}
	}
}
