package config

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

var AppConfig *AppStruct

type AppStruct struct {
	WeWork WeWork
}
type WeWork struct {
	Corpid     string `yaml:"corpid"`
	Corpsecret string `yaml:"corpsecret"`
}

//获取基础配置项
func ParseAppConfig() error {
	ymlFile := GetYmlFile("app")
	if ymlFile == "" {
		return errors.New("configuration file of config not found")
	}
	data, _ := ioutil.ReadFile(ymlFile)
	yaml.Unmarshal(data, &AppConfig)
	if Config.HTTP.Mode == "" {
		return errors.New("configuration file of config parse fail")
	}
	return nil
}
