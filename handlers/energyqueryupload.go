package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/energietransitie/twomes-backoffice-api/services"
	"github.com/energietransitie/twomes-backoffice-api/twomes/authorization"
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryupload"
	"github.com/sirupsen/logrus"
)

type EnergyQueryUploadHandler struct {
	service *services.EnergyQueryUploadService
}

// Create a new EnergyQueryUploadHandler.
func NewEnergyQueryUploadHandler(service *services.EnergyQueryUploadService) *EnergyQueryUploadHandler {
	return &EnergyQueryUploadHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new EnergyQueryUpload.
func (h *EnergyQueryUploadHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request energyqueryupload.EnergyQueryUpload
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	auth, ok := r.Context().Value(AuthorizationCtxKey).(*authorization.Authorization)
	if !ok {
		return NewHandlerError(err, "unauthorized", http.StatusUnauthorized).WithMessage("failed when getting authentication context value")
	}

	if !auth.IsKind(authorization.AccountToken) {
		return NewHandlerError(err, "wrong token kind", http.StatusForbidden).WithMessage("wrong token kind was used")
	}

	EnergyQueryUpload, err := h.service.Create(request.QueryID, request.BuildingID, request.EnergyQueryValues)
	if err != nil {
		if errors.Is(err, services.ErrEmptyQueryUpload) {
			return NewHandlerError(err, "empty EnergyQueryUpload", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	// We don't need to return all measurements in the EnergyQueryUpload response.
	EnergyQueryUpload.EnergyQueryValues = nil

	err = json.NewEncoder(w).Encode(&EnergyQueryUpload)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}
