package api

import (
	"gin-IM/request"
	"gin-IM/service"
	"net/http"

	"github.com/gin-gonic/gin"
	logging "github.com/sirupsen/logrus"
)

func UserRegister(c *gin.Context) {
	var userRegister request.UserRegisterRequest
	err := c.ShouldBind(&userRegister)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "json数据解析失败！",
			"error":  err.Error(),
		})
		logging.Info(err)
		return
	}
	res := service.UserRegister(userRegister)
	// gin.H其实就是个map[string]any嘛。
	// 其实map也好，struct也罢，c.JSON函数签名中第二个参数为obj any
	// 只要能被序列化成json即可
	c.JSON(http.StatusOK, res)
}

func UserLogin(c *gin.Context) {
	var userRegister request.UserRegisterRequest
	err := c.ShouldBind(&userRegister)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": http.StatusBadRequest,
			"msg":    "json数据解析失败！",
			"error":  err.Error(),
		})
		logging.Info(err)
		return
	}
	ua := c.GetHeader("User-Agent")
	res := service.UserLogin(ua, userRegister)
	c.JSON(http.StatusOK, res)
}
