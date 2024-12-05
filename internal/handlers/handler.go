package handlers

import (
	"html/template"
	"net/http"
)

func NotFoundHandler(templatesDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		tmpl, err := template.ParseFiles(templatesDir + "/404-admin.html")
		if err != nil {
			http.Error(w, "Error loading 404 page", http.StatusInternalServerError)
			return
		}
		tmpl.Execute(w, nil)
	}
}
