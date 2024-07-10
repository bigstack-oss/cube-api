package main

import (
	"flag"

	"github.com/bigstack-oss/cube-api/internal/config"
	"github.com/bigstack-oss/cube-api/internal/runtime"
	"github.com/bigstack-oss/cube-api/internal/service"
	log "go-micro.dev/v5/logger"
)

var (
	conf string
)

func init() {
	flag.StringVar(&conf, "conf", "", "")
	flag.Parse()
}

// @title     Cube API
// @version   1.0
// @BasePath  /api
func main() {
	conf, err := config.Load(conf)
	if err != nil {
		log.Errorf("failed to load config: %s", err.Error())
		return
	}

	runtime, err := runtime.NewRuntime(conf)
	if err != nil {
		log.Errorf("failed to init runtime: %s", err.Error())
		return
	}

	err = service.WrapGoMicro(runtime).Run()
	if err != nil {
		log.Errorf("failed to run service: %s", err.Error())
	}
}
