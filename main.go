package main

import (
	"flag"

	"github.com/liuyong-go/gin_project/bootstrap"
	"github.com/liuyong-go/gin_project/config"
	"github.com/toolkits/pkg/runner"
)

var (
	f *string
)

func init() {
	f = flag.String("f", "", "config path")
	flag.Parse()
	runner.Init()
	config.InitBaseInfo()
	if *f != "" {
		config.BaseInfo.ConfigPath = *f
	}

}

//接收配置参数设置一些全局变量
//调用bootstrap启动项目
func main() {
	bootstrap.Start()
	//time.Sleep(5 * time.Second)
}
