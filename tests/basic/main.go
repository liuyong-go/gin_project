package main

import (
	"context"
	"fmt"

	"github.com/liuyong-go/gin_project/app/models"
	"github.com/liuyong-go/gin_project/bootstrap"
	"github.com/liuyong-go/gin_project/config"
	"github.com/liuyong-go/gin_project/libs/logger"
	"github.com/liuyong-go/gin_project/libs/weapp"
)

func main() {
	bootstrap.TestInit()
	//testTrace()
	//getJson()
	//testweapp()
	testDB()
}
func testweapp() {
	var ctx = context.Background()
	weapp, err := weapp.NewWeapp(ctx, "test")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(weapp)
}
func testDB() {
	var data = models.NewApis()
	data.SiteId = 1
	data.DepartmentId = "2"
	t := data.Get(context.Background())
	logger.Info(context.TODO(), "dataid", t)
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
func getJson() {
	test := config.GetTest(context.Background())
	fmt.Println(test)
	test = config.GetTest(context.Background())
	fmt.Println(test)
}
