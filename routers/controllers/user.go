package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi1310/mycloud-disk/pkg/session"
	"github.com/wangyi1310/mycloud-disk/serializer"
	"github.com/wangyi1310/mycloud-disk/services"
)

func UserLogin(c *gin.Context) {
	var login services.LoginUser
	if err := c.ShouldBindJSON(&login); err != nil {
		c.JSON(200, serializer.Err(serializer.CodeParamErr, "Failed to parse params", err))
		return
	}
	user, err := services.Login(&login)
	if err != nil {
		c.JSON(200, serializer.Err(serializer.CodeCredentialInvalid, "Wrong email or password", err))
		return
	}

	session.SetSession(c, map[string]interface{}{
		"user_id": user.ID,
	})
	c.JSON(200, user)
}

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	// 注册逻辑
	var register services.RegisterUser
	err := c.BindJSON(&register)
	var res serializer.Response
	if err != nil {
		c.JSON(200, serializer.Err(serializer.CodeParamErr, "Failed to parse params", err))
		return
	}
	res = services.Register(&register)
	c.JSON(200, res)
}

func UserActive(c *gin.Context) {
	uid, _ := c.Get("object_id")
	res := services.Activate(&services.ActiveUser{
		Uid: uid,
	})
	c.JSON(200, res)
}
