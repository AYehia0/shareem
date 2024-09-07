package app

import (
	"html/template"
	"net/http"
	"shareem/internal/handler"
)

var TemplatePath = "./templates/*"

// reload the application routes
func (a *App) reloadRoutes() {
	tmpl := template.Must(template.New("").ParseGlob("./templates/*"))

	appHandler := handler.NewServer(a.logger, tmpl, a.db)

	files := http.FileServer(http.Dir("./static"))

	// handle the static files
	a.router.Handle("GET /static/", http.StripPrefix("/static/", files))

	// handle the index page
	a.router.HandleFunc("GET /{$}", appHandler.Index)

	// handle the share creation
	a.router.HandleFunc("POST /{$}", appHandler.Insert)
}
