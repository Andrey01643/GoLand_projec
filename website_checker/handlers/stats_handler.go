package handlers

import (
	"fmt"
	"net/http"

	"website_checker/utils"
)

func StatsHandler(w http.ResponseWriter, r *http.Request) {
	adminRequests := utils.GetAdminRequests()

	fmt.Fprintf(w, "User requests: %v\n", utils.GetUserRequests())
	fmt.Fprintf(w, "Admin requests: %v\n", adminRequests)
}
