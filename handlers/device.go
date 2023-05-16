package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/services"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
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
func (h *DeviceHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request twomes.Device
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	auth, ok := r.Context().Value(AuthorizationCtxKey).(*twomes.Authorization)
	if !ok {
		return NewHandlerError(err, "unauthorized", http.StatusUnauthorized).WithMessage("failed when getting authentication context value").WithLevel(logrus.ErrorLevel)
	}

	if !auth.IsKind(twomes.AccountToken) {
		return NewHandlerError(err, "wrong token kind", http.StatusForbidden).WithMessage("wrong token kind was used")
	}

	device, err := h.service.Create(request.Name, request.DeviceType, request.BuildingID, auth.ID, request.ActivationSecret)
	if err != nil {
		if helpers.IsMySQLRecordNotFoundError(err) {
			return NewHandlerError(err, "not found", http.StatusNotFound)
		}

		if helpers.IsMySQLDuplicateError(err) {
			return NewHandlerError(err, "duplicate", http.StatusBadRequest)
		}

		if errors.Is(err, services.ErrBuildingDoesNotBelongToAccount) {
			return NewHandlerError(err, err.Error(), http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(&device)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}

// Handle API endpoint for activating a device.
func (h *DeviceHandler) Activate(w http.ResponseWriter, r *http.Request) error {
	var request twomes.Device
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "unauthorized", http.StatusUnauthorized).WithMessage("device name present")
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return NewHandlerError(err, "unauthorized", http.StatusUnauthorized).WithMessage("authorization header not present")
	}

	authHeader = strings.Split(authHeader, "Bearer ")[1]

	if authHeader == "" {
		logrus.Info("authorization malformed")
		return NewHandlerError(err, "unauthorized", http.StatusUnauthorized).WithMessage("authorization malformed")
	}

	device, err := h.service.Activate(request.Name, authHeader)
	if err != nil {
		if errors.Is(err, twomes.ErrDeviceActivationSecretIncorrect) {
			return NewHandlerError(err, "forbidden", http.StatusForbidden)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	// We don't need to share all uploads.
	device.Uploads = nil

	err = json.NewEncoder(w).Encode(&device)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}

// Handle API endpoint for getting device information.
func (h *DeviceHandler) GetDeviceByName(w http.ResponseWriter, r *http.Request) error {
	deviceName := chi.URLParam(r, "device_name")
	if deviceName == "" {
		return NewHandlerError(nil, "device_name not specified", http.StatusBadRequest)
	}

	device, err := h.service.GetByName(deviceName)
	if err != nil {
		return NewHandlerError(err, "device not found", http.StatusNotFound).WithMessage("device not found")
	}

	auth, ok := r.Context().Value(AuthorizationCtxKey).(*twomes.Authorization)
	if !ok {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithMessage("failed when getting authentication context value")
	}

	accountID, err := h.service.GetAccountByDeviceID(device.ID)
	if err != nil {
		return NewHandlerError(err, "device not found", http.StatusNotFound).WithMessage("device could not be found by ID")
	}

	if auth.ID != accountID {
		return NewHandlerError(err, "device does not belong to account", http.StatusForbidden).WithMessage("request was made for device not owned by account")
	}

	// We don't need to share all uploads.
	device.Uploads = nil

	err = json.NewEncoder(w).Encode(&device)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}

// Handle API endpoint for getting device measurements
func (h *DeviceHandler) GetDeviceMeasurements(w http.ResponseWriter, r *http.Request) error {
	deviceName := chi.URLParam(r, "device_name")

	auth, ok := r.Context().Value(AuthorizationCtxKey).(*twomes.Authorization)
	if !ok {
		return NewHandlerError(nil, "internal server error", http.StatusInternalServerError).WithMessage("failed when getting authentication context value")
	}

	device, err := h.getDeviceByName(deviceName, auth.ID)
	if err != nil {
		return err
	}

	measurements, err := h.service.GetMeasurementsByDeviceID(device.ID)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithMessage("failed when getting measurements")
	}

	err = json.NewEncoder(w).Encode(&measurements)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}

// Handle API endpoint for getting device properties
func (h *DeviceHandler) GetDeviceProperties(w http.ResponseWriter, r *http.Request) error {
	deviceName := chi.URLParam(r, "device_name")

	auth, ok := r.Context().Value(AuthorizationCtxKey).(*twomes.Authorization)
	if !ok {
		return NewHandlerError(nil, "internal server error", http.StatusInternalServerError).WithMessage("failed when getting authentication context value")
	}

	device, err := h.getDeviceByName(deviceName, auth.ID)
	if err != nil {
		return err
	}

	properties, err := h.service.GetPropertiesByDeviceID(device.ID)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithMessage("failed when getting measurements")
	}

	err = json.NewEncoder(w).Encode(&properties)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}

func (h *DeviceHandler) getDeviceByName(deviceName string, accountId uint) (*twomes.Device, error) {
	if deviceName == "" {
		return nil, NewHandlerError(nil, "device_name not specified", http.StatusBadRequest)
	}

	device, err := h.service.GetByName(deviceName)
	if err != nil {
		return nil, NewHandlerError(err, "device not found", http.StatusNotFound).WithMessage("device not found")
	}

	deviceAccountId, err := h.service.GetAccountByDeviceID(device.ID)
	if err != nil {
		return nil, NewHandlerError(err, "device not found", http.StatusNotFound).WithMessage("device could not be found by ID")
	}

	if deviceAccountId != accountId {
		return nil, NewHandlerError(nil, "device does not belong to account", http.StatusForbidden).WithMessage("request was made for device not owned by account")
	}

	return &device, nil
}
