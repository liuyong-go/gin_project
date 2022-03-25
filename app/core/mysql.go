package core

import (
	"context"

	"github.com/liuyong-go/gin_project/libs/logger"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

var DB *gorm.DB

type MysqlConfig struct {
	Master MysqlSingle `yaml:"master"`
	Slave  MysqlSingle `yaml:"slave"`
}
type MysqlSingle struct {
	Addr  string `yaml:"addr"`
	Max   int    `yaml:"max"`
	Idle  int    `yaml:"idle"`
	Debug bool   `yaml:"debug"`
}

func InitMysql(conf MysqlConfig) (err error) {
	DB, err = gorm.Open(mysql.Open(conf.Master.Addr), &gorm.Config{})
	if err != nil {
		logger.Info(context.Background(), "mysql connect fail", err)
	}
	DB.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{mysql.Open(conf.Slave.Addr)},
	}))
	return
}
