package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/energietransitie/twomes-api/internal/helpers"
	"github.com/energietransitie/twomes-api/pkg/ports"
	"github.com/energietransitie/twomes-api/pkg/services"
	"github.com/energietransitie/twomes-api/pkg/twomes"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type DeviceHandler struct {
	service ports.DeviceService
}

// Create a new DeviceHandler.
func NewDeviceHandler(service ports.DeviceService) *DeviceHandler {
	return &DeviceHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new device.
func (h *DeviceHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request twomes.Device
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	auth, ok := r.Context().Value(AuthorizationCtxKey).(*twomes.Authorization)
	if !ok {
		logrus.Error("failed when getting authentication context value")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !auth.IsKind(twomes.AccountToken) {
		logrus.Infof("%s token was used while %s was required", auth.Kind, twomes.AccountToken)
		http.Error(w, "wrong token kind", http.StatusForbidden)
		return
	}

	device, err := h.service.Create(request.Name, request.DeviceType, request.BuildingID, auth.ID, request.ActivationSecret)
	if err != nil {
		if helpers.IsMySQLRecordNotFoundError(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		if helpers.IsMySQLDuplicateError(err) {
			http.Error(w, "duplicate", http.StatusBadRequest)
			return
		}

		if errors.Is(err, services.ErrBuildingDoesNotBelongToAccount) {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		logrus.Info(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&device)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// Handle API endpoint for activating a device.
func (h *DeviceHandler) Activate(w http.ResponseWriter, r *http.Request) {
	var request twomes.Device
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logrus.Warning("device name present")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		logrus.Info("authorization header not present")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authHeader = strings.Split(authHeader, "Bearer ")[1]

	if authHeader == "" {
		logrus.Info("authorization malformed")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	device, err := h.service.Activate(request.Name, authHeader)
	if err != nil {
		if errors.Is(err, twomes.ErrDeviceActivationSecretIncorrect) {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		logrus.Info(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(&device)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// Handle API endpoint for getting device information.
func (h *DeviceHandler) GetDeviceByName(w http.ResponseWriter, r *http.Request) {
	deviceName := chi.URLParam(r, "device_name")
	if deviceName == "" {
		http.Error(w, "device_name not specified", http.StatusBadRequest)
		return
	}

	device, err := h.service.GetByName(deviceName)
	if err != nil {
		logrus.Info("device not found")
		http.Error(w, "device not found", http.StatusNotFound)
		return
	}

	auth, ok := r.Context().Value(AuthorizationCtxKey).(*twomes.Authorization)
	if !ok {
		logrus.Error("failed when getting authentication context value")
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	accountID, err := h.service.GetAccountByDeviceID(device.ID)
	if err != nil {
		logrus.Info("device could not be found by ID")
		http.Error(w, "device not found", http.StatusNotFound)
		return
	}

	if auth.ID != accountID {
		logrus.Info("request was made for device not owned by account")
		http.Error(w, "device does not belong to account", http.StatusForbidden)
		return
	}

	// We don't need to share all uploads.
	device.Uploads = nil

	err = json.NewEncoder(w).Encode(&device)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
