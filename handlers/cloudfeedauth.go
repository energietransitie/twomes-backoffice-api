package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"github.com/sirupsen/logrus"
)

type CloudFeedAuthHandler struct {
	service ports.CloudFeedAuthService
}

// Create a new CloudFeedAuthHandler.
func NewCloudFeedAuthHandler(service ports.CloudFeedAuthService) *CloudFeedAuthHandler {
	return &CloudFeedAuthHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new cloud feed.
func (h *CloudFeedAuthHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request twomes.CloudFeedAuth
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	auth, ok := r.Context().Value(AuthorizationCtxKey).(*twomes.Authorization)
	if !ok {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithMessage("failed when getting authentication context value")
	}

	if auth.ID != request.AccountID {
		return NewHandlerError(nil, "account ID in response does not match authorization token", http.StatusBadRequest)
	}

	_, err = h.service.Create(auth.ID, request.CloudFeedID, request.AuthGrantToken)
	if err != nil {
		if helpers.IsMySQLDuplicateError(err) {
			return NewHandlerError(err, "duplicate", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	return nil
}
