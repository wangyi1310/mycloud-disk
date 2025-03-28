package bootstrap

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/wangyi1310/mycloud-disk/conf"
	"github.com/wangyi1310/mycloud-disk/pkg/log"
	"github.com/wangyi1310/mycloud-disk/pkg/request"
)

type GitHubRelease struct {
	URL  string `json:"html_url"`
	Name string `json:"name"`
	Tag  string `json:"tag_name"`
}

func InitApplication() {
	fmt.Print(`
  __  ____     __   _____ _      ____  _    _ _____    _____ _____  _____ _  __
 |  \/  \ \   / /  / ____| |    / __ \| |  | |  __ \  |  __ \_   _|/ ____| |/ /
 | \  / |\ \_/ /  | |    | |   | |  | | |  | | |  | | | |  | || | | (___ | ' / 
 | |\/| | \   /   | |    | |   | |  | | |  | | |  | | | |  | || |  \___ \|  <  
 | |  | |  | |    | |____| |___| |__| | |__| | |__| | | |__| || |_ ____) | . \ 
 |_|  |_|  |_|     \_____|______\____/ \____/|_____/  |_____/_____|_____/|_|\_\
                                                                               
 V` + conf.BackendVersion + `  Commit #` + conf.LastCommit + `  Pro=` + conf.IsPro + `
================================================

`)
	go CheckUpdate()
}

// CheckUpdate 检查更新
func CheckUpdate() {
	client := request.NewClient()
	res, err := client.Request("GET", "https://api.github.com/repos/wangyi1310/mycloud-disk/releases", nil).GetResponse()
	if err != nil {
		log.Log().Warning("更新检查失败, %s", err)
		return
	}

	var list []GitHubRelease
	if err := json.Unmarshal([]byte(res), &list); err != nil {
		log.Log().Warning("更新检查失败, %s", err)
		return
	}

	if len(list) > 0 {
		present, err1 := version.NewVersion(conf.BackendVersion)
		latest, err2 := version.NewVersion(list[0].Tag)
		if err1 == nil && err2 == nil && latest.GreaterThan(present) {
			log.Log().Info("有新的版本 [%s] 可用，下载：%s", list[0].Name, list[0].URL)
		}
	}

}
