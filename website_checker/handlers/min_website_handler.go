package handlers

import (
	"fmt"
	"net/http"

	"website_checker/utils"
)

func MinWebsiteHandler(w http.ResponseWriter, r *http.Request) {
	minWebsite := utils.GetMinWebsite()

	fmt.Fprintf(w, "Website with minimum access time: %s", minWebsite)

	// Обновляем статистику запросов пользователей
	utils.UpdateUserRequests(2)
}
