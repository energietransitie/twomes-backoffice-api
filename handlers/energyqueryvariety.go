package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/services"
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryvariety"
	"github.com/sirupsen/logrus"
)

type EnergyQueryVarietyHandler struct {
	service *services.EnergyQueryVarietyService
}

// Create a new EnergyQueryVarietyHandler.
func NewEnergyQueryVarietyHandler(service *services.EnergyQueryVarietyService) *EnergyQueryVarietyHandler {
	return &EnergyQueryVarietyHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new EnergyQueryVariety.
func (h *EnergyQueryVarietyHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request energyqueryvariety.EnergyQueryVariety
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	EnergyQueryVariety, err := h.service.Create(request.Name)

	if err != nil {
		if helpers.IsMySQLRecordNotFoundError(err) {
			return NewHandlerError(err, "not found", http.StatusNotFound)
		}

		if helpers.IsMySQLDuplicateError(err) {
			return NewHandlerError(err, "duplicate", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(EnergyQueryVariety)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}
