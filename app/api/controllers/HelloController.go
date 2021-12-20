package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/liuyong-go/gin_project/libs/logger"
)

type Hello struct {
}

func NewHello() *Hello {
	return &Hello{}
}
func (h *Hello) World(c *gin.Context) {
	logger.Info(c, "请求接口")
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "hello world",
	})
}
