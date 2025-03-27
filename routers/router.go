package routers

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi1310/mycloud-disk/conf"
	"github.com/wangyi1310/mycloud-disk/pkg/log"
	"github.com/wangyi1310/mycloud-disk/routers/controllers"
)

func Init() *gin.Engine {
	if conf.SystemConfig.Mode == "Master" {
		log.Log().Info("Current runing mode: Master")
		return InitMaster()
	} else {
		log.Log().Info("Current runing mode: Slave")
		return InitSlave()
	}
}

func InitMaster() *gin.Engine {
	r := gin.Default()
	v3 := r.Group("/api/v3")
	site := v3.Group("site")
	{
		site.GET("/ping", controllers.Ping)
		// site.GET("captcha")
		// site.GET("config")
	}
	return r
}

func InitSlave() *gin.Engine {
	r := gin.Default()
	return r
}
