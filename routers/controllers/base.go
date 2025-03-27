package controllers

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi1310/mycloud-disk/conf"
	"github.com/wangyi1310/mycloud-disk/serializer"
)

func Ping(c *gin.Context) {
	version := conf.BackendVersion
	c.JSON(200, serializer.Response{
		Code: 0,
		Data: version,
	})
}
