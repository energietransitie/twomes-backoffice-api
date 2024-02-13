package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes/devicetype"
	"github.com/sirupsen/logrus"
)

type DeviceTypeHandler struct {
	service ports.DeviceTypeService
}

// Create a new DeviceTypeHandler.
func NewDeviceTypeHandler(service ports.DeviceTypeService) *DeviceTypeHandler {
	return &DeviceTypeHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new device type.
func (h *DeviceTypeHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request devicetype.DeviceType
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	dt, err := h.service.Create(request.Name, request.InstallationManualURL, request.InfoURL)
	if err != nil {
		if helpers.IsMySQLDuplicateError(err) {
			return NewHandlerError(err, "duplicate", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(&dt)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}
