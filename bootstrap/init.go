package bootstrap

import (
	"github.com/wangyi1310/mycloud-disk/conf"
	"github.com/wangyi1310/mycloud-disk/models"
)

func Init(path string) {
	InitApplication()

	if !conf.SystemConfig.Debug {
	}

	dependencies := []struct {
		mode    string
		factory func()
	}{
		{
			"master",
			func() {
				models.Init()
			},
		},
	}

	for _, dependency := range dependencies {
		if dependency.mode == conf.SystemConfig.Mode || dependency.mode == "both" {
			dependency.factory()
		}
	}

}
