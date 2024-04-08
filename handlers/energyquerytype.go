package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/services"
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyquerytype"
	"github.com/sirupsen/logrus"
)

type EnergyQueryTypeHandler struct {
	service *services.EnergyQueryTypeService
}

// Create a new EnergyQueryTypeHandler.
func NewEnergyQueryTypeHandler(service *services.EnergyQueryTypeService) *EnergyQueryTypeHandler {
	return &EnergyQueryTypeHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new EnergyQueryType.
func (h *EnergyQueryTypeHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request energyquerytype.EnergyQueryType
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	EnergyQueryType, err := h.service.Create(request.QueryVariety, request.Formula)

	if err != nil {
		if helpers.IsMySQLRecordNotFoundError(err) {
			return NewHandlerError(err, "not found", http.StatusNotFound)
		}

		if helpers.IsMySQLDuplicateError(err) {
			return NewHandlerError(err, "duplicate", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(EnergyQueryType)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}
