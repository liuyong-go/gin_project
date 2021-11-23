package libs

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

type MysqlConfig struct {
	Addr  string `yaml:"addr"`
	Max   int    `yaml:"max"`
	Idle  int    `yaml:"idle"`
	Debug bool   `yaml:"debug"`
}

func InitMysql(conf MysqlConfig) (err error) {
	DB, err = gorm.Open(mysql.Open(conf.Addr), &gorm.Config{})
	return
}
