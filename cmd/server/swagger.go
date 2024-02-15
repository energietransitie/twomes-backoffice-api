package main

import (
	"io/fs"
	"net/http"

	"github.com/energietransitie/twomes-backoffice-api/handlers"
	"github.com/energietransitie/twomes-backoffice-api/swaggerdocs"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func setupSwaggerDocs(r *chi.Mux, baseURL string) {
	swaggerUI, err := fs.Sub(swaggerdocs.StaticFiles, "swagger-ui")
	if err != nil {
		logrus.Fatal(err)
	}

	docsHandler, err := handlers.NewDocsHandler(swaggerdocs.StaticFiles, baseURL)
	if err != nil {
		logrus.Fatal(err)
	}

	r.Method("GET", "/openapi.yml", handlers.Handler(docsHandler.OpenAPISpec))                        // Serve openapi.yml
	r.Method("GET", "/docs/*", http.StripPrefix("/docs/", http.FileServer(http.FS(swaggerUI))))       // Serve static files.
	r.Method("GET", "/docs", handlers.Handler(docsHandler.RedirectDocs(http.StatusMovedPermanently))) // Redirect /docs to /docs/
	r.Method("GET", "/", handlers.Handler(docsHandler.RedirectDocs(http.StatusSeeOther)))             // Redirect / to /docs/
}
