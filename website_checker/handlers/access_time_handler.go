package handlers

import (
	"fmt"
	"net/http"

	"website_checker/utils"
)

func AccessTimeHandler(w http.ResponseWriter, r *http.Request) {
	website := r.URL.Query().Get("website")

	accessTime := utils.GetAccessTime(website)

	if accessTime > 0 {
		fmt.Fprintf(w, "Access time to %s: %s", website, accessTime)
	} else {
		fmt.Fprintf(w, "Website %s is not accessible", website)
	}

	// Обновляем статистику запросов пользователей
	utils.UpdateUserRequests(1)
}
