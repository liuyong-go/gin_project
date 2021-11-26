package main

import (
	"github.com/liuyong-go/gin_project/app/models"
	"github.com/liuyong-go/gin_project/bootstrap"
	"github.com/liuyong-go/gin_project/libs/logger"
)

func main() {
	bootstrap.TestInit()
	var data = models.NewApis()
	data.SiteId = 1
	data.DepartmentId = "2"
	data.Create()
	logger.Info("dataid", data.ID)
}
