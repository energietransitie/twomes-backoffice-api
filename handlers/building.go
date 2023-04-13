package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type BuildingHandler struct {
	service ports.BuildingService
}

// Create a new BuildingHandler.
func NewBuildingHandler(service ports.BuildingService) *BuildingHandler {
	return &BuildingHandler{
		service: service,
	}
}

// Handle API endpoint for getting building information.
func (h *BuildingHandler) GetBuildingByID(w http.ResponseWriter, r *http.Request) error {
	buildingIDParam := chi.URLParam(r, "building_id")
	if buildingIDParam == "" {
		return NewHandlerError(nil, "building_id not specified", http.StatusBadRequest)
	}

	buildingID, err := strconv.ParseUint(buildingIDParam, 10, 64)
	if err != nil {
		return NewHandlerError(err, "building_id not a number", http.StatusBadRequest)
	}

	auth, ok := r.Context().Value(AuthorizationCtxKey).(*twomes.Authorization)
	if !ok {
		return NewHandlerError(nil, "internal server error", http.StatusInternalServerError).WithMessage("failed when getting authentication context value")
	}

	building, err := h.service.GetByID(uint(buildingID))
	if err != nil {
		return NewHandlerError(err, "not found", http.StatusNotFound)
	}

	if building.AccountID != auth.ID {
		return NewHandlerError(nil, "building does not belong to account", http.StatusForbidden).WithMessage("request was made for building not owned by account")
	}

	err = json.NewEncoder(w).Encode(&building)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}
