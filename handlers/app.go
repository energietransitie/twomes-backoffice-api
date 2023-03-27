package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/energietransitie/twomes-api/internal/helpers"
	"github.com/energietransitie/twomes-api/pkg/ports"
	"github.com/energietransitie/twomes-api/pkg/twomes"
	"github.com/sirupsen/logrus"
)

type AppHandler struct {
	service ports.AppService
}

func NewAppHandler(service ports.AppService) *AppHandler {
	return &AppHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new app.
func (h *AppHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request twomes.App
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	app, err := h.service.Create(request.Name, request.ProvisioningURLTemplate)
	if err != nil {
		if helpers.IsMySQLRecordNotFoundError(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		if helpers.IsMySQLDuplicateError(err) {
			http.Error(w, "duplicate", http.StatusBadRequest)
			return
		}

		logrus.Info(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(app)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
