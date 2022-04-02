package qiniu

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"strings"
	"time"
)

//获取图片的privateurl
func (a *Auth) GetPrivateUrl(baseUrl string, expire time.Duration) (privateUrl string) {
	if !strings.Contains(baseUrl, "private") {
		privateUrl = baseUrl
		return
	}
	if strings.Contains(baseUrl, "?") {
		baseUrl += "&e="
	} else {
		baseUrl += "?e="
	}
	deadline := fmt.Sprintf("%d", time.Now().Add(expire*time.Second).Unix())
	baseUrl += deadline
	h := hmac.New(sha1.New, []byte(a.secretKey))
	h.Write([]byte(baseUrl))
	sign := base64.URLEncoding.EncodeToString(h.Sum(nil))
	signUrl := fmt.Sprintf("%s:%s", a.accessKey, sign)
	privateUrl = baseUrl + "&token=" + signUrl
	return
}
