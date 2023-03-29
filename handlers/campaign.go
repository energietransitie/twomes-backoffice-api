package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/energietransitie/twomes-api/internal/helpers"
	"github.com/energietransitie/twomes-api/pkg/ports"
	"github.com/energietransitie/twomes-api/pkg/twomes"
	"github.com/sirupsen/logrus"
)

type CampaignHandler struct {
	service ports.CampaignService
}

// Create a new CampaignHandler.
func NewCampaignHandler(service ports.CampaignService) *CampaignHandler {
	return &CampaignHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new campaign.
func (h *CampaignHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request twomes.Campaign
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	campaign, err := h.service.Create(request.Name, request.App, request.InfoURL, request.StartTime, request.EndTime)
	if err != nil {
		if helpers.IsMySQLRecordNotFoundError(err) {
			return NewHandlerError(err, "not found", http.StatusNotFound)
		}

		if helpers.IsMySQLDuplicateError(err) {
			return NewHandlerError(err, "duplicate", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(campaign)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}
