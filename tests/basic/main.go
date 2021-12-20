package main

import (
	"context"
	"fmt"

	"github.com/liuyong-go/gin_project/app/models"
	"github.com/liuyong-go/gin_project/bootstrap"
	"github.com/liuyong-go/gin_project/libs/logger"
)

func main() {
	bootstrap.TestInit()
	testTrace()
}
func testDB() {
	var data = models.NewApis()
	data.SiteId = 1
	data.DepartmentId = "2"
	data.Create(context.Background())
	logger.Info(context.TODO(), "dataid", data.ID)
}
func testTrace() {
	ctx := context.Background()
	stateStr := "key1=value1,key2=value2"
	traceID := "1234"

	ltrace, err := logger.InjectTraceContext(ctx, traceID, stateStr)
	ctx = context.WithValue(ctx, "trace", ltrace)
	logger.Info(ctx, "测试log")
	fmt.Println("error", err)
}
