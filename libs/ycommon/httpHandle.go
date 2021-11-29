package ycommon

import (
	"io/ioutil"
	"net/http"
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
