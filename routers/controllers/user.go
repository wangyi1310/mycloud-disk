package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi1310/mycloud-disk/serializer"
	"github.com/wangyi1310/mycloud-disk/services"
)

// UserRegister 用户注册
func UserRegister(c *gin.Context) {
	// 注册逻辑
	var register services.RegisterUser
	err := c.BindJSON(&register)
	var res serializer.Response
	if err != nil {
		res = serializer.Err(serializer.CodeParamErr, "参数错误", err)
	}
	res = userService.Register(&register)
	c.JSON(200, res)
}

func UserActive(c *gin.Context) {
	uid, _ := c.Get("object_id")
	res := userService.Activate(&services.ActiveUser{
		Uid: uid,
	})
	c.JSON(200, res)
}
