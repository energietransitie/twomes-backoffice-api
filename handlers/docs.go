package handlers

import (
	"io/fs"
	"net/http"
	"text/template"

	"github.com/sirupsen/logrus"
)

// Type used to fill openapi.template.yml template.
type templateData struct {
	BaseURL string
}

type DocsHandler struct {
	template *template.Template
	baseURL  templateData
}

// Create a new DocsHandler.
func NewDocsHandler(files fs.FS, baseURL string) (*DocsHandler, error) {
	t, err := template.ParseFS(files, "openapi.template.yml")
	if err != nil {
		return nil, err
	}

	return &DocsHandler{
		template: t,
		baseURL:  templateData{baseURL},
	}, nil
}

// Handle API endpoint for displaying OpenAPI spec.
// This file should be displayed as openapi.yml.
func (h *DocsHandler) OpenAPISpec(w http.ResponseWriter, r *http.Request) error {
	err := h.template.Execute(w, h.baseURL)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}

// Handle redirection from /docs to /docs/ to serve static files.
func (h *DocsHandler) RedirectDocs(redirectCode int) func(http.ResponseWriter, *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		redirectURL := "/docs/"

		// Fix prefix stripped by Traefik proxy to ensure correcy redirection.
		redirectURL = r.Header.Get("X-Forwarded-Prefix") + redirectURL

		http.Redirect(w, r, redirectURL, redirectCode)

		return nil
	}
}
