package main

import (
	"context"
	"fmt"
	"time"

	"github.com/liuyong-go/gin_project/app/core"
	"github.com/liuyong-go/gin_project/app/models"
	"github.com/liuyong-go/gin_project/bootstrap"
	"github.com/liuyong-go/gin_project/config"
	"github.com/liuyong-go/gin_project/libs/logger"
	"github.com/liuyong-go/gin_project/libs/weapp"
	"github.com/liuyong-go/gin_project/libs/wework"
	"github.com/liuyong-go/gin_project/libs/ycommon"
)

func main() {
	bootstrap.TestInit()
	//testTrace()
	//getJson()
	//testweapp()
	//testDB()
	//testRedis()
	//testWeapp()
	//testPage()
}
func testPage() {
	result := ycommon.Paginator(5, 3, 100)
	fmt.Println(result)
}
func testWeapp() {
	wework.SendGroupMsg()
}
func testRedis() {
	cmd := core.RedisCore.Set("test", "test", time.Second*60)
	fmt.Println(cmd.Err())
	value := core.RedisCore.Get("test").Val()
	fmt.Println("value", value)
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
	ctx := context.Background()
	var data = models.NewArticle()
	// data.Uid = 1
	// data.Title = "测试标题1"
	// data.Desc = "简单描述1"
	// data.Content = "正文内容2"
	// data.UniqId = ycommon.UUID()
	// data.State = 1
	// data.Insert(ctx)
	data.GetByUniqID(ctx, "0")
	data.Del(ctx)
	// data.Title = "修改测试标题2"
	// data.Save(ctx)
	//data.Incr(ctx, "view_nums")
	// where := map[string]interface{}{
	// 	"uid": 1,
	// }
	// result := data.PageList(where, 1, 10, "id desc")
	//fmt.Println(result)
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
