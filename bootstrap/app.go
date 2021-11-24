package bootstrap

import (
	"fmt"

	"github.com/liuyong-go/gin_project/config"
)

func Start() {
	err := config.ParseConfig()
	if err != nil {
		fmt.Println(err)
	}
}
