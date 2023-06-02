package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"github.com/sirupsen/logrus"
)

type CloudFeedHandler struct {
	service ports.CloudFeedService
}

// Create a new CloudFeedHandler.
func NewCloudFeedHandler(service ports.CloudFeedService) *CloudFeedHandler {
	return &CloudFeedHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new cloud feed.
func (h *CloudFeedHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request twomes.CloudFeed
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	cloudFeed, err := h.service.Create(request.Name, request.AuthorizationURL, request.TokenURL, request.ClientID, request.ClientSecret, request.Scope, request.RedirectURL)
	if err != nil {
		if helpers.IsMySQLDuplicateError(err) {
			return NewHandlerError(err, "duplicate", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(&cloudFeed)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}
