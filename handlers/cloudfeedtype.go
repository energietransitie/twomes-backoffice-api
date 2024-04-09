package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/services"
	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeed"
	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeedtype"
	"github.com/sirupsen/logrus"
)

type CloudFeedTypeHandler struct {
	service *services.CloudFeedTypeService
}

// Create a new CloudFeedTypeHandler.
func NewCloudFeedTypeHandler(service *services.CloudFeedTypeService) *CloudFeedTypeHandler {
	return &CloudFeedTypeHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new cloud feed.
func (h *CloudFeedTypeHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request cloudfeedtype.CloudFeedType
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	CloudFeedType, err := h.service.Create(request.Name, request.AuthorizationURL, request.TokenURL, request.ClientID, request.ClientSecret, request.Scope, request.RedirectURL)
	if err != nil {
		if helpers.IsMySQLDuplicateError(err) {
			return NewHandlerError(err, "duplicate", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(&CloudFeedType)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
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
func (h *CloudFeedHandler) Download(args DownloadArgs, reply *string) error {
	cfa, err := h.service.Find(cloudfeed.CloudFeed{AccountID: args.AccountID, CloudFeedTypeID: args.CloudFeedID})
	if err != nil {
		return err
	}

	err = h.service.Download(context.Background(), cfa, args.StartPeriod, args.EndPeriod)
	if err != nil {
		return err
	}

	*reply = "Downloaded data from cloud feed. Check server logs for more information."

	return nil
}
