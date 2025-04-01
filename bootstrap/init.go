package bootstrap

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi1310/mycloud-disk/conf"
	"github.com/wangyi1310/mycloud-disk/models"
	"github.com/wangyi1310/mycloud-disk/pkg/auth"
	"github.com/wangyi1310/mycloud-disk/pkg/cache"
	"github.com/wangyi1310/mycloud-disk/pkg/email"
)

func Init(path string) {
	InitApplication()
	conf.Init(path)
	// Debug 关闭时，切换为生产模式
	if !conf.SystemConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	dependencies := []struct {
		mode    string
		factory func()
	}{
		{
			"both",
			func() {
				cache.Init()
			},
		},
		{
			"master",
			func() {
				models.Init()
			},
		},
		{
			"both",
			func() {
				auth.Init()
			},
		},
		{
			"both",
			func() {
				email.Init()
			},
		},
	}

	for _, dependency := range dependencies {
		if dependency.mode == conf.SystemConfig.Mode || dependency.mode == "both" {
			dependency.factory()
		}
	}

}
