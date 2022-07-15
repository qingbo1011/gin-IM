package service

import (
	"gin-IM/db/mysql"
	"gin-IM/model"
	"gin-IM/request"
	"gin-IM/response"
	"net/http"
)

func UserRegister(register request.UserRegisterRequest) response.Response {
	var user model.User
	var count int
	mysql.MysqlDB.Where(&model.User{UserName: register.UserName}).First(&user).Count(&count)
	if count == 1 {
		return response.Response{
			Status: http.StatusForbidden,
			Msg:    "用户名重复！",
		}
	}
	// 如果数据库中没有该用户，那么就开始注册
	user.UserName = register.UserName
	err := user.SetPassword(register.Password)
	if err != nil {
		return response.Response{
			Status: http.StatusInternalServerError,
			Msg:    "数据库插入数据出错！",
			Error:  err.Error(),
		}
	}
	// 加密成功就可以创建用户了
	err = mysql.MysqlDB.Create(&user).Error
	if err != nil {
		return response.Response{
			Status: http.StatusInternalServerError,
			Msg:    "数据库添加数据出错！",
			Error:  err.Error(),
		}
	}
	return response.Response{
		Status: http.StatusOK,
		Msg:    "用户注册成功！",
	}
}
