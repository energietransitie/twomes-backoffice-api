package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/energietransitie/twomes-api/pkg/ports"
	"github.com/energietransitie/twomes-api/pkg/twomes"
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
func (h *UploadHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request twomes.Upload
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

	if !auth.IsKind(twomes.DeviceToken) {
		logrus.Infof("%s token was used while %s was required", auth.Kind, twomes.DeviceToken)
		http.Error(w, "wrong token kind", http.StatusForbidden)
		return
	}

	upload, err := h.service.Create(auth.ID, request.DeviceTime, request.Measurements)
	if err != nil {
		logrus.Info(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// We don't need to return all measurements in the upload response.
	upload.Measurements = nil

	err = json.NewEncoder(w).Encode(&upload)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
