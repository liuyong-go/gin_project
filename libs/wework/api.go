package wework

import (
	"context"
	"fmt"
	"time"

	"github.com/liuyong-go/gin_project/config"
	"github.com/liuyong-go/gin_project/libs/ycommon"
	"github.com/tidwall/gjson"
)

var tokenData *accessToken

type accessToken struct {
	accessToken string
	expireTime  int64
}

//获取accessToken ，单体应用可以这么用，多个应用需要借助redis
func getAccessToken() (string, error) {
	if tokenData != nil {
		if time.Now().Unix() < tokenData.expireTime {
			return tokenData.accessToken, nil
		}
	}
	//获取token，设置过期时间
	url := fmt.Sprintf(accessTokenUrl, config.AppConfig.WeWork.Corpid, config.AppConfig.WeWork.Corpsecret)
	body, err := ycommon.HttpGet(url)
	if err != nil {
		return "", err
	}
	token := gjson.GetBytes(body, "access_token").String()
	if token != "" {
		expireIn := gjson.GetBytes(body, "expires_in").Int()
		expireTime := time.Now().Unix() + expireIn - 300
		tokenData = &accessToken{accessToken: token, expireTime: expireTime}
		return token, nil
	} else {
		errorCode := gjson.GetBytes(body, "errcode").String()
		errorMessage := gjson.GetBytes(body, "errmsg").String()
		return "", fmt.Errorf(errorCode + " " + errorMessage)
	}
}
func AddGroup() {
	token, err := getAccessToken()
	if err != nil {
		fmt.Println(err)
	}
	url := fmt.Sprintf(groupCreateUrl, token)
	header := map[string]interface{}{
		"Content-Type": "application/json",
	}
	var data = map[string]interface{}{
		"name":  "创建测试",
		"owner": "LiuYong",
		"userlist": []string{
			"LiuYong", "loyal",
		},
	}
	body, err := ycommon.PostRow(context.Background(), url, header, data, 1)
	fmt.Println(string(body))
}
func SendGroupMsg() {
	chatID := "wr4q6BCQAAd4QUa0i1Qz0F5WEzZLpGfA"
	token, err := getAccessToken()
	if err != nil {
		fmt.Println(err)
	}
	url := fmt.Sprintf(sendGroupMsg, token)
	header := map[string]interface{}{
		"Content-Type": "application/json",
	}
	var data = map[string]interface{}{
		"chatid":  chatID,
		"msgtype": "text",
		"text": map[string]string{
			"content": "你的快递已到\n请携带工卡前往邮件中心领取",
		},
		"safe": 0,
	}
	body, err := ycommon.PostRow(context.Background(), url, header, data, 1)
	fmt.Println(string(body))

}
func GroupList() {
	token, err := getAccessToken()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(token)
}
