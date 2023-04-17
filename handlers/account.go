// Package handlers defines available HTTP handlers
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"github.com/sirupsen/logrus"
)

type AccountHandler struct {
	accountService ports.AccountService
}

func NewAccountHandler(accountService ports.AccountService) *AccountHandler {
	return &AccountHandler{
		accountService: accountService,
	}
}

// Handle API endpoint for creating a new account.
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request twomes.Account
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	account, err := h.accountService.Create(request.Campaign)
	if err != nil {
		if helpers.IsMySQLRecordNotFoundError(err) {
			return NewHandlerError(err, "not found", http.StatusNotFound)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}

// Handle API endpoint for activating an account.
// This endpoint should be protected with an account activation token.
func (h *AccountHandler) Activate(w http.ResponseWriter, r *http.Request) error {
	var request twomes.Building
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	auth, ok := r.Context().Value(AuthorizationCtxKey).(*twomes.Authorization)
	if !ok {
		return NewHandlerError(err, "unauthorized", http.StatusUnauthorized).WithMessage("failed when getting authentication context value")
	}

	if !auth.IsKind(twomes.AccountActivationToken) {
		return NewHandlerError(err, "wrong token kind", http.StatusForbidden).WithMessage("wrong token kind was used")
	}

	account, err := h.accountService.Activate(auth.ID, request.Longtitude, request.Latitude, request.TZName)
	if err != nil {
		if errors.Is(err, twomes.ErrAccountAlreadyActivated) {
			return NewHandlerError(err, "account already activated", http.StatusBadRequest)
		}

		return NewHandlerError(err, "account activation failed", http.StatusBadRequest)
	}

	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}
