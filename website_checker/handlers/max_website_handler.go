package handlers

import (
	"fmt"
	"net/http"

	"website_checker/utils"
)

func MaxWebsiteHandler(w http.ResponseWriter, r *http.Request) {
	maxWebsite := utils.GetMaxWebsite()

	fmt.Fprintf(w, "Website with maximum access time: %s", maxWebsite)

	// Обновляем статистику запросов пользователей
	utils.UpdateUserRequests(3)
}
