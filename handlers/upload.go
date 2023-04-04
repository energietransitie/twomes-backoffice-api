package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/energietransitie/twomes-backoffice-api/pkg/ports"
	"github.com/energietransitie/twomes-backoffice-api/pkg/services"
	"github.com/energietransitie/twomes-backoffice-api/pkg/twomes"
	"github.com/sirupsen/logrus"
)

type UploadHandler struct {
	service ports.UploadService
}

// Create a new UploadHandler.
func NewUploadHandler(service ports.UploadService) *UploadHandler {
	return &UploadHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new upload.
func (h *UploadHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request twomes.Upload
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	auth, ok := r.Context().Value(AuthorizationCtxKey).(*twomes.Authorization)
	if !ok {
		return NewHandlerError(err, "unauthorized", http.StatusUnauthorized).WithMessage("failed when getting authentication context value")
	}

	if !auth.IsKind(twomes.DeviceToken) {
		return NewHandlerError(err, "wrong token kind", http.StatusForbidden).WithMessage("wrong token kind was used")
	}

	upload, err := h.service.Create(auth.ID, request.DeviceTime, request.Measurements)
	if err != nil {
		if errors.Is(err, services.ErrEmptyUpload) {
			return NewHandlerError(err, "empty upload", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	// We don't need to return all measurements in the upload response.
	upload.Measurements = nil

	err = json.NewEncoder(w).Encode(&upload)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}
