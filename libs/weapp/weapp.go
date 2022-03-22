package weapp

import (
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/liuyong-go/gin_project/libs/logger"
	"github.com/liuyong-go/gin_project/libs/ycommon"
	"github.com/tidwall/gjson"
)

var TokenData *AceessToken

var weappMap = map[string]Weapp{}

type Weapp struct {
	AppName string `json:"app_name"`
	Appid   string `json:"app_id"`
	Secret  string `json:"secret"`
}

//登录授权
type GrantData struct {
	OpenID     string
	UnionID    string
	SessionKey string
}
type AceessToken struct {
	aceessToken string
	expireTime  int64
}

//获得临时登录凭证 code 后传到开发者服务器调用此接口完成登录流程
const CODE_SESSION_URL = "https://api.weixin.qq.com/sns/jscode2session?appid=%s&secret=%s&js_code=%s&grant_type=authorization_code"
const ACCESS_TOKEN_URL = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"

//go:embed weapp.json
var weappjson string

//初始化账号
func NewWeapp(ctx context.Context, appName string) (weapp Weapp, err error) {
	var ok bool
	if weapp, ok = weappMap[appName]; ok {
		return
	}

	err = json.Unmarshal([]byte(weappjson), &weappMap)
	if err != nil {
		logger.Info(ctx, "解析weapp.json失败:", err)
	}

	if weapp, ok = weappMap[appName]; ok {
		return
	} else {
		err = errors.New("获取不到小程序账号信息:" + appName)
	}
	return
}

//登录凭证校验。通过 wx.login 接口获得临时登录凭证 code 后传到开发者服务器调用此接口完成登录流程
func (weapp *Weapp) CodeToSession(jsCode string) (*GrantData, error) {
	url := fmt.Sprintf(CODE_SESSION_URL, weapp.Appid, weapp.Secret, jsCode)
	body, err := ycommon.HttpGet(url)
	if err != nil {
		return &GrantData{}, err
	}
	openID := gjson.GetBytes(body, "openid").String()
	errorCode := gjson.GetBytes(body, "errcode").String()
	if openID == "" {
		errorMessage := gjson.GetBytes(body, "errmsg").String()
		return &GrantData{}, errors.New(errorCode + " " + errorMessage)
	}

	return &GrantData{
		OpenID:     openID,
		UnionID:    gjson.GetBytes(body, "unionid").String(),
		SessionKey: gjson.GetBytes(body, "session_key").String(),
	}, nil
}
func (weapp *Weapp) TestAccessToken() (string, error) {
	return weapp.getAccessToken()
}

//获取accessToken ，单体应用可以这么用，多个应用需要借助redis
func (weapp *Weapp) getAccessToken() (string, error) {
	if TokenData != nil {
		if time.Now().Unix() < TokenData.expireTime {
			return TokenData.aceessToken, nil
		}
	}
	//获取token，设置过期时间
	url := fmt.Sprintf(ACCESS_TOKEN_URL, weapp.Appid, weapp.Secret)
	body, err := ycommon.HttpGet(url)
	if err != nil {
		return "", err
	}
	token := gjson.GetBytes(body, "access_token").String()
	if token != "" {
		expireIn := gjson.GetBytes(body, "access_token").Int()
		expireTime := time.Now().Unix() + expireIn - 300
		TokenData = &AceessToken{aceessToken: token, expireTime: expireTime}
		return token, nil
	} else {
		errorCode := gjson.GetBytes(body, "errcode").String()
		errorMessage := gjson.GetBytes(body, "errmsg").String()
		return "", errors.New(errorCode + " " + errorMessage)
	}
}
