package utils

import (
	"net/http"
	"time"
)

var websites = []string{
	"google.com",
	"youtube.com",
	"facebook.com",
	"baidu.com",
	"wikipedia.org",
	"qq.com",
	"taobao.com",
	"yahoo.com",
	"tmall.com",
	"amazon.com",
	"google.co.in",
	"twitter.com",
	"sohu.com",
	"jd.com",
	"live.com",
	"instagram.com",
	"sina.com.cn",
	"weibo.com",
	"google.co.jp",
	"reddit.com",
	"vk.com",
	"360.cn",
	"login.tmall.com",
	"blogspot.com",
	"yandex.ru",
	"google.com.hk",
	"netflix.com",
	"linkedin.com",
	"pornhub.com",
	"google.com.br",
	"twitch.tv",
	"pages.tmall.com",
	"csdn.net",
	"yahoo.co.jp",
	"mail.ru",
	"aliexpress.com",
	"alipay.com",
	"office.com",
	"google.fr",
	"google.ru",
	"google.co.uk",
	"microsoftonline.com",
	"google.de",
	"ebay.com",
	"microsoft.com",
	"livejasmin.com",
	"t.co",
	"bing.com",
	"xvideos.com",
	"google.ca",
}

func CheckWebsites() {
	for {
		for _, website := range websites {
			startTime := time.Now()
			_, err := http.Get("https://" + website)
			elapsedTime := time.Since(startTime)
			if err == nil {
				accessTimesMutex.Lock()
				accessTimes[website] = elapsedTime
				accessTimesMutex.Unlock()
			}
		}
		time.Sleep(1 * time.Minute)
	}
}
