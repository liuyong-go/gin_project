package ycommon

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/idoubi/goz"
	"github.com/liuyong-go/gin_project/libs/logger"
)

func HttpGet(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}
	return body, nil
}
func PostForm(ctx context.Context, url string, header map[string]interface{}, formData map[string]interface{}, timeout float32) (result []byte, err error) {
	cli := goz.NewClient()

	resp, err := cli.Post(url, goz.Options{
		Headers:    header,
		FormParams: formData,
		Timeout:    timeout,
	})
	if err != nil {
		logger.Info(ctx, "call_url_api_fail", url, err)
		return
	}

	result, _ = resp.GetBody()
	return
}
func PostRow(ctx context.Context, url string, header map[string]interface{}, rowData interface{}, timeout float32) (result []byte, err error) {
	if err != nil {
		return
	}
	cli := goz.NewClient()
	resp, err := cli.Request(http.MethodPost, url, goz.Options{
		Headers: header,
		JSON:    rowData, //提交map,最终解析为json
		Timeout: timeout,
	})
	if err != nil {
		logger.Info(ctx, "call_url_api_fail", url, err)
		return
	}

	result, _ = resp.GetBody()
	return
}

// 发送GET请求
// url：         请求地址
// response：    请求返回的内容
func Get(ctx context.Context, url string, timeout time.Duration) string {

	// 超时时间：秒
	client := &http.Client{Timeout: timeout * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		logger.Info(ctx, "call_url_api_fail", url, err)
	}
	defer resp.Body.Close()
	var buffer [512]byte
	result := bytes.NewBuffer(nil)
	for {
		n, err := resp.Body.Read(buffer[0:])
		result.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			logger.Info(ctx, "call_url_api_fail", url, err)
		}
	}

	return result.String()
}
