package handlers

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

func GetAPIDocs(w http.ResponseWriter, router *mux.Router) {
	var routes []string
	err := router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) (err error) {
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		methods, err := route.GetMethods()
		if err != nil {
			return err
		}
		routes = append(routes, fmt.Sprintf("<li><strong>%s %s</strong>: %s</li>", strings.Join(methods, "/"), pathTemplate, ""))
		return nil
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to generate API docs")
		return
	}

	// Генерация и отправка документации API
	docs := fmt.Sprintf(`
    <!DOCTYPE html>
    <html>
    <head>
        <title>Task API Documentation</title>
    </head>
    <body>
        <h1>Task API Documentation</h1>
        <h2>Endpoints</h2>
        <ul>
            %s
        </ul>
    </body>
    </html>
    `, strings.Join(routes, ""))

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(docs))
}
