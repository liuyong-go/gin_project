package config

import (
	"errors"
	"io/ioutil"
	"os"
	"path"

	"github.com/liuyong-go/gin_project/app/core"
	"github.com/liuyong-go/gin_project/libs/logger"
	"github.com/toolkits/pkg/file"
	"github.com/toolkits/pkg/runner"
	"gopkg.in/yaml.v2"
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
	if env == "" {
		env = "local"
	}
	rootPath := runner.Cwd
	var configPath string
	if env == "local" {
		rootPath = "/Users/liuyong/go/src/gin_project"
	}
	configPath = rootPath + "/etc"
	BaseInfo = &Base{
		RootPath:   rootPath,
		ConfigPath: configPath,
		Env:        env,
	}
	err := ParseConfig()
	if err != nil {
		panic(err)
	}
	err = ParseAppConfig()
	if err != nil {
		panic(err)
	}
}

type ConfigStruct struct {
	Logger logger.LoggerStruct
	HTTP   httpStruct
	RPC    rpcStruct
	MySQL  core.MysqlConfig
	Redis  core.RedisMap
	Es     EsConfig
}
type EsConfig struct {
	Address  []string `yaml:"address"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
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

//获取基础配置项
func ParseConfig() error {
	ymlFile := GetYmlFile("config")
	if ymlFile == "" {
		return errors.New("configuration file of config not found")
	}
	data, _ := ioutil.ReadFile(ymlFile)
	yaml.Unmarshal(data, &Config)
	if Config.HTTP.Mode == "" {
		return errors.New("configuration file of config parse fail")
	}
	return nil
}
func GetYmlFile(module string) string {
	filename := module + "." + BaseInfo.Env + ".yaml"
	ymlFile := path.Join(BaseInfo.ConfigPath, filename)
	if file.IsExist(ymlFile) {
		return ymlFile
	}
	return ""
}
