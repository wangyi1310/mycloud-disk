package main

import (
	"flag"
	"net"
	"net/http"
	"os"

	"github.com/wangyi1310/mycloud-disk/bootstrap"
	"github.com/wangyi1310/mycloud-disk/conf"
	"github.com/wangyi1310/mycloud-disk/pkg/log"
	"github.com/wangyi1310/mycloud-disk/pkg/util"
	"github.com/wangyi1310/mycloud-disk/routers"
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

func runUnix(server *http.Server) error {
	listener, err := net.Listen("unix", conf.UnixConfig.Listen)
	if err != nil {
		return err
	}

	defer listener.Close()
	defer os.Remove(conf.UnixConfig.Listen)

	if conf.UnixConfig.Perm > 0 {
		err = os.Chmod(conf.UnixConfig.Listen, os.FileMode(conf.UnixConfig.Perm))
		if err != nil {
			log.Log().Warning(
				"Failed to set permission to %q for socket file %q: %s",
				conf.UnixConfig.Perm,
				conf.UnixConfig.Listen,
				err,
			)
		}
	}

	return server.Serve(listener)
}

func main() {
	api := routers.Init()
	api.TrustedPlatform = conf.SystemConfig.ProxyHeader
	server := &http.Server{Handler: api}
	if conf.UnixConfig.Listen != "" {
		// delete socket file before listening
		if _, err := os.Stat(conf.UnixConfig.Listen); err == nil {
			if err = os.Remove(conf.UnixConfig.Listen); err != nil {
				log.Log().Error("Failed to delete socket file: %s", err)
				return
			}
		}

		log.Log().Info("Listening to %q", conf.UnixConfig.Listen)
		if err := runUnix(server); err != nil {
			log.Log().Error("Failed to listen to %q: %s", conf.UnixConfig.Listen, err)
		}
		return
	}

	log.Log().Info("Listening to %q", conf.SystemConfig.Listen)
	server.Addr = conf.SystemConfig.Listen
	if err := server.ListenAndServe(); err != nil {
		log.Log().Error("Failed to listen to %q: %s", conf.SystemConfig.Listen, err)
	}

}
