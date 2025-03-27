package main

import (
	"flag"
	"fmt"

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
}

func main() {
	fmt.Println("Hello World")
}