package bootstrap

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/liuyong-go/gin_project/app/models"
	"github.com/liuyong-go/gin_project/config"
	"github.com/liuyong-go/gin_project/libs/logger"
	"github.com/liuyong-go/gin_project/libs/ydefer"
	"github.com/liuyong-go/gin_project/libs/yhttp"
)

func Start() {
	err := config.ParseConfig()
	if err != nil {
		panic(err)
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	logger.InitLogger(config.Config.Logger)
	err = models.InitMysql(config.Config.MySQL)
	if err != nil {
		logger.Warn("db init fail", err)
	}
	logger.Info("ctx 后续使用", ctx)
	yhttp.Start()
	endingProc(cancelFunc)
}

//监听停止服务信号
func endingProc(cancelFunc context.CancelFunc) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	<-c
	fmt.Printf("stop signal caught, stopping... pid=%d\n", os.Getpid())
	ydefer.Clean()
	// 执行清理工作
	cancelFunc()
	yhttp.Shutdown()

	fmt.Println("process stopped successfully")
}
