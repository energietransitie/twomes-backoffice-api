package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/energietransitie/twomes-api/internal/helpers"
	"github.com/energietransitie/twomes-api/pkg/ports"
	"github.com/energietransitie/twomes-api/pkg/twomes"
	"github.com/sirupsen/logrus"
)

type PropertyHandler struct {
	service ports.PropertyService
}

// Create a new PropertyHandler.
func NewPropertyHandler(service ports.PropertyService) *PropertyHandler {
	return &PropertyHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new property.
func (h *PropertyHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request twomes.Property
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	property, err := h.service.Create(request.Name, request.Unit)
	if err != nil {
		if helpers.IsMySQLDuplicateError(err) {
			return NewHandlerError(err, "duplicate", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(&property)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}
