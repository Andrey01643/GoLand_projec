package main

import (
	"fmt"
	"net/http"

	"website_checker/handlers"
	"website_checker/utils"
)

func main() {
	go utils.CheckWebsites()

	//http://localhost:8080/access-time?website=google.com

	http.HandleFunc("/access-time", handlers.AccessTimeHandler)
	http.HandleFunc("/min-website", handlers.MinWebsiteHandler)
	http.HandleFunc("/max-website", handlers.MaxWebsiteHandler)
	http.HandleFunc("/stats", handlers.StatsHandler)

	fmt.Println("Server listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
