package ycommon

import "time"

//获取之间加上指定时长后的时间 time.Now(), "15s" ms,ns,s,m,h格式
func GetCalculateTime(currentTimer time.Time, d string) (time.Time, error) {
	duration, err := time.ParseDuration(d) //解析字符串获取时间
	if err != nil {
		return time.Time{}, err
	}

	return currentTimer.Add(duration), nil
}
