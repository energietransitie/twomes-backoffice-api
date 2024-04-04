package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/services"
	"github.com/energietransitie/twomes-backoffice-api/twomes/authorization"
	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeedauth"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

type CloudFeedAuthHandler struct {
	service *services.CloudFeedAuthService
}

// Create a new CloudFeedAuthHandler.
func NewCloudFeedAuthHandler(service *services.CloudFeedAuthService) *CloudFeedAuthHandler {
	return &CloudFeedAuthHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new cloud feed.
func (h *CloudFeedAuthHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request cloudfeedauth.CloudFeedAuth
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	auth, ok := r.Context().Value(AuthorizationCtxKey).(*authorization.Authorization)
	if !ok {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithMessage("failed when getting authentication context value")
	}

	_, err = h.service.Create(r.Context(), auth.ID, request.CloudFeedID, request.AuthGrantToken)
	if err != nil {
		if helpers.IsMySQLDuplicateError(err) {
			return NewHandlerError(err, "duplicate", http.StatusBadRequest)
		}

		if _, ok := err.(*oauth2.RetrieveError); ok {
			return NewHandlerError(err, "invalid auth code exchange", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	return nil
}

type DownloadArgs struct {
	AccountID   uint
	CloudFeedID uint
	StartPeriod time.Time
	EndPeriod   time.Time
}

// Handle RPC endpoint for downloading data from a cloud feed.
func (h *CloudFeedAuthHandler) Download(args DownloadArgs, reply *string) error {
	cfa, err := h.service.Find(cloudfeedauth.CloudFeedAuth{AccountID: args.AccountID, CloudFeedID: args.CloudFeedID})
	if err != nil {
		return err
	}

	return h.service.Download(context.Background(), cfa, args.StartPeriod, args.EndPeriod)
}
