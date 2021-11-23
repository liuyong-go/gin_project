package config

import (
	"fmt"
	"os"

	"github.com/liuyong-go/gin_project/libs"
	"github.com/toolkits/pkg/runner"
)

var BaseInfo *Base
var Config *ConfigStruct

type Base struct {
	RootPath   string //根目录
	ConfigPath string //config目录
	Env        string //环境local,prod
}

func InitBaseInfo() {
	env := os.Getenv("APP_ENV")
	rootPath := runner.Cwd
	fmt.Println("root", rootPath)
	var configPath string
	if env == "local" {
		rootPath = "/Users/liuyong/go/src/gin_project"
	}
	configPath = rootPath + "/config"
	BaseInfo = &Base{
		RootPath:   rootPath,
		ConfigPath: configPath,
		Env:        env,
	}
}

type ConfigStruct struct {
	Logger libs.LoggerStruct
	HTTP   httpStruct
	RPC    rpcStruct
	MySQL  libs.MysqlConfig
}

type httpStruct struct {
	Mode           string `yaml:"mode"`
	Listen         string `yaml:"listen"`
	Pprof          bool   `yaml:"pprof"`
	CookieName     string `yaml:"cookieName"`
	CookieDomain   string `yaml:"cookieDomain"`
	CookieSecure   bool   `yaml:"cookieSecure"`
	CookieHttpOnly bool   `yaml:"cookieHttpOnly"`
	CookieMaxAge   int    `yaml:"cookieMaxAge"`
	CookieSecret   string `yaml:"cookieSecret"`
	CsrfSecret     string `yaml:"csrfSecret"`
}
type rpcStruct struct {
	Listen string `yaml:"listen"`
}

func Parse() error {

	return nil
}
