package main

import (
	"flag"

	"github.com/wangyi1310/mycloud-disk/bootstrap"
	"github.com/wangyi1310/mycloud-disk/pkg/log"
	"github.com/wangyi1310/mycloud-disk/pkg/util"
)

var (
	isEject    bool
	confPath   string
	scriptName string
)

func init() {
	flag.StringVar(&confPath, "c", util.RelativePath("conf.ini"), "Path to the config file.")
	flag.Parse()
	bootstrap.Init(confPath)
}

func main() {
	log.Log().Warning("Hello World")
}
