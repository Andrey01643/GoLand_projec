package utils

import (
	"sync"
	"time"
)

var accessTimes = make(map[string]time.Duration)
var accessTimesMutex sync.RWMutex

func GetAccessTime(website string) time.Duration {
	accessTimesMutex.RLock()
	defer accessTimesMutex.RUnlock()

	return accessTimes[website]
}

func GetMinWebsite() string {
	var minWebsite string
	var minTime time.Duration

	accessTimesMutex.RLock()
	defer accessTimesMutex.RUnlock()

	for website, accessTime := range accessTimes {
		if minTime == 0 || accessTime < minTime {
			minTime = accessTime
			minWebsite = website
		}
	}

	return minWebsite
}

func GetMaxWebsite() string {
	var maxWebsite string
	var maxTime time.Duration

	accessTimesMutex.RLock()
	defer accessTimesMutex.RUnlock()

	for website, accessTime := range accessTimes {
		if accessTime > maxTime {
			maxTime = accessTime
			maxWebsite = website
		}
	}

	return maxWebsite
}
