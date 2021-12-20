package main

import (
	"context"

	"github.com/liuyong-go/gin_project/app/models"
	"github.com/liuyong-go/gin_project/config"
	"github.com/liuyong-go/gin_project/libs/logger"
	"github.com/liuyong-go/gin_project/libs/modelsgenerate"
	"github.com/toolkits/pkg/runner"
)

func main() {
	geneModels()
}
func geneModels() {
	runner.Init()
	config.InitBaseInfo()
	err := config.ParseConfig()
	if err != nil {
		panic(err)
	}
	logger.InitLogger(config.Config.Logger)
	err = models.InitMysql(config.Config.MySQL)
	if err != nil {
		logger.Warn(context.TODO(), "db init fail", err)
	}
	modelPath := config.BaseInfo.RootPath + "/app/models/"
	var tables = []string{"apis"}
	modelsgenerate.NewGenerator(models.DB, modelPath, "test").Genertate(tables...)
}
