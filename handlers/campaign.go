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
func (h *CampaignHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request twomes.Campaign
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	campaign, err := h.service.Create(request.Name, request.App, request.InfoURL, request.StartTime, request.EndTime)
	if err != nil {
		if helpers.IsMySQLRecordNotFoundError(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		if helpers.IsMySQLDuplicateError(err) {
			http.Error(w, "duplicate", http.StatusBadRequest)
			return
		}

		logrus.Info(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(campaign)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
