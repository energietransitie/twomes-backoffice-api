package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/energietransitie/twomes-api/internal/helpers"
	"github.com/energietransitie/twomes-api/pkg/ports"
	"github.com/energietransitie/twomes-api/pkg/twomes"
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
func (h *DeviceTypeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request twomes.DeviceType
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	deviceType, err := h.service.Create(request.Name, request.InstallationManualURL, request.InfoURL, request.Properties, request.UploadInterval)
	if err != nil {
		if helpers.IsMySQLDuplicateError(err) {
			http.Error(w, "duplicate", http.StatusBadRequest)
			return
		}

		logrus.Info(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&deviceType)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
