package api

import (
	"gin-IM/pkg/util"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Test(c *gin.Context) {
	token := c.GetHeader("token")
	claims, err := util.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status": http.StatusInternalServerError,
			"error":  err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "欢迎：" + claims.Username,
	})
}
