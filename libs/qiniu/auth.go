package qiniu

import (
	"context"
	"time"

	"github.com/liuyong-go/gin_project/config"
)

var authData *Auth

type Auth struct {
	accessKey  string
	secretKey  string
	expireTime time.Time
}

func NewAuth(ctx context.Context) *Auth {
	return &Auth{
		accessKey: config.AppConfig.Qiniu.Accesskey,
		secretKey: config.AppConfig.Qiniu.Secretkey,
	}
}
