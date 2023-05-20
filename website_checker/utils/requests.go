package utils

import "sync"

var userRequests = make(map[int]int)
var userRequestsMutex sync.RWMutex

var adminRequests = make(map[int]int)
var adminRequestsMutex sync.RWMutex

func GetUserRequests() map[int]int {
	userRequestsMutex.RLock()
	defer userRequestsMutex.RUnlock()

	requests := make(map[int]int)
	for k, v := range userRequests {
		requests[k] = v
	}

	return requests
}

func UpdateUserRequests(requestType int) {
	userRequestsMutex.Lock()
	defer userRequestsMutex.Unlock()

	userRequests[requestType]++
}

func GetAdminRequests() map[int]int {
	adminRequestsMutex.RLock()
	defer adminRequestsMutex.RUnlock()

	requests := make(map[int]int)
	for k, v := range adminRequests {
		requests[k] = v
	}

	return requests
}
