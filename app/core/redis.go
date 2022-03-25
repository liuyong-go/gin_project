package core

import (
	"github.com/go-redis/redis"
)

var rdbs = make(map[string]*redis.Client)
var RedisCore *redis.Client

const RedisCoreName = "test"

type RedisMap map[string]RedisConfig
type RedisConfig struct {
	Name     string `yaml:"name"`
	Addr     string `yaml:"addr"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func InitRedis(redisConf RedisMap) {
	if conf, ok := redisConf[RedisCoreName]; ok {
		RedisCore = newRedis(conf)
	}

}

func newRedis(redisConfig RedisConfig) *redis.Client {
	if rdbs[redisConfig.Name] != nil {
		return rdbs[redisConfig.Name]
	}
	rdbs[redisConfig.Name] = redis.NewClient(&redis.Options{
		Addr:     redisConfig.Addr,
		Password: redisConfig.Password, // no password set
		DB:       redisConfig.DB,       // use default DB
	})
	return rdbs[redisConfig.Name]
}
