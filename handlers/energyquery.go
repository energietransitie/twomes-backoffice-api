package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/services"
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyquery"
	"github.com/sirupsen/logrus"
)

type EnergyQueryHandler struct {
	service *services.EnergyQueryService
}

// Create a new EnergyQueryHandler.
func NewEnergyQueryHandler(service *services.EnergyQueryService) *EnergyQueryHandler {
	return &EnergyQueryHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new EnergyQuery.
func (h *EnergyQueryHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request energyquery.EnergyQuery
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	EnergyQuery, err := h.service.Create(request.Name, request.Formula)

	if err != nil {
		if helpers.IsMySQLRecordNotFoundError(err) {
			return NewHandlerError(err, "not found", http.StatusNotFound)
		}

		if helpers.IsMySQLDuplicateError(err) {
			return NewHandlerError(err, "duplicate", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(EnergyQuery)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}
