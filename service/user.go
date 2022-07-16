package service

import (
	"gin-IM/db/mysql"
	"gin-IM/model"
	"gin-IM/pkg/util"
	"gin-IM/request"
	"gin-IM/response"
	"net/http"

	"github.com/jinzhu/gorm"
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

func UserLogin(login request.UserRegisterRequest) response.Response {
	var user model.User
	err := mysql.MysqlDB.Where(&model.User{UserName: login.UserName}).First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) { // 数据库中没有找到记录
			return response.Response{
				Status: http.StatusBadRequest,
				Msg:    "该用户不存在，请先注册！",
				Error:  err.Error(),
			}
		}
		// 不是用户不存在却还是出错，其他不可抗拒的因素
		return response.Response{
			Status: http.StatusInternalServerError,
			Msg:    "查询数据库出现错误！",
			Error:  err.Error(),
		}
	}
	// 用户从数据库中找到了，检验密码
	ok, err := user.CheckPassword(login.Password)
	if err != nil {
		return response.Response{
			Status: http.StatusInternalServerError,
			Msg:    "登录失败！",
			Error:  err.Error(),
		}
	}
	if !ok {
		return response.Response{
			Status: http.StatusForbidden,
			Msg:    "密码错误，登录失败！",
		}
	}
	// 登录成功要分发token（其他功能需要身份验证，给前端存储的）
	token, err := util.GenerateToken(user.ID, user.UserName)
	if err != nil {
		return response.Response{
			Status: http.StatusInternalServerError,
			Msg:    "token签发失败！",
			Error:  err.Error(),
		}
	}
	return response.Response{
		Status: http.StatusOK,
		Msg:    "登录成功！",
		Data:   map[string]string{"token": token},
	}
}
