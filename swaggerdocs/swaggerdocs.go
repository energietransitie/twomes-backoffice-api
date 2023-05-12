// Package swaggerdocs embeds swagger UI and OpenAPI spec files
// that can be server by and http.FileServer.
package swaggerdocs

import (
	"embed"
)

// Static files for swagger UI and OpenAPI spec.
//
//go:embed swagger-ui
//go:embed openapi.yml
var StaticFiles embed.FS
