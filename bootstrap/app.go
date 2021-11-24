package bootstrap

import (
	"github.com/liuyong-go/gin_project/config"
	"github.com/liuyong-go/gin_project/libs/logger"
)

func Start() {
	err := config.ParseConfig()
	if err != nil {
		panic(err)
	}
	logger.InitLogger(config.Config.Logger)
	logger.Info(config.Config)
}
