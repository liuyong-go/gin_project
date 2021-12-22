package config

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/liuyong-go/gin_project/libs/logger"
)

//go:embed test.json
var testJson string
var testResult = &Test{}

type Test struct {
	Test string `json:"test"`
}

func GetTest(ctx context.Context) *Test {
	if testResult.Test != "" {
		return testResult
	}
	fmt.Println("result", testResult)
	fmt.Println("json", []byte(testJson))

	err := json.Unmarshal([]byte(testJson), testResult)
	if err != nil {
		logger.Info(ctx, err)
	}
	return testResult
}
