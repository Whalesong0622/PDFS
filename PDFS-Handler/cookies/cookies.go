package cookies

import (
	"math/rand"
	"time"
)

type info struct {
	username    string
	createdtime int64
}

var cookies_map = make(map[string]info)

// 生成Cookie并记录到cookies_map中
func LoginCookie(username string) []byte {
	cookie := make([]byte, 20)
	rand.Read(cookie)
	tmpinfo := info{
		username:    username,
		createdtime: time.Now().Unix(),
	}
	cookies_map[string(cookie)] = tmpinfo

	return cookie
}

func DelectCookie(cookie []byte) {
	delete(cookies_map, string(cookie))
}

// 检查Cookie对应的用户名是否存在，如果有，返回。
func CookieToUsername(cookie []byte) string {
	tmpinfo, ok := cookies_map[string(cookie)]

	if !ok {
		return ""
	} else if (time.Now().Unix()-tmpinfo.createdtime)/3600 >= 1 {
		delete(cookies_map, string(cookie))
		return ""
	}
	return tmpinfo.username
}

// 每十分钟进行一次检查
func CookiesDaemon() {
	for {
		time.Sleep(time.Minute * 10)
		CookiesMapCleaner()
	}
}

// 遍历cookies，对于超过一小时的cookie进行删除
func CookiesMapCleaner() {
	for key, val := range cookies_map {
		if (time.Now().Unix()-val.createdtime)/3600 >= 1 {
			delete(cookies_map, key)
		}
	}
}
