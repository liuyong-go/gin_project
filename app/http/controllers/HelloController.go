package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Hello struct {
}

func NewHello() *Hello {
	return &Hello{}
}
func (h *Hello) World(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  "hello world",
	})
}
