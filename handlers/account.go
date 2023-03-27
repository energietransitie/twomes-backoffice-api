// Package handlers defines available HTTP handlers
package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/energietransitie/twomes-api/internal/helpers"
	"github.com/energietransitie/twomes-api/pkg/ports"
	"github.com/energietransitie/twomes-api/pkg/twomes"
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
func (h *AccountHandler) Create(w http.ResponseWriter, r *http.Request) {
	var request twomes.Account
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	account, err := h.accountService.Create(request.Campaign)
	if err != nil {
		if helpers.IsMySQLRecordNotFoundError(err) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		logrus.Error(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

// Handle API endpoint for activating an account.
// This endpoint should be protected with an account activation token.
func (h *AccountHandler) Activate(w http.ResponseWriter, r *http.Request) {
	var request twomes.Building
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	auth, ok := r.Context().Value(AuthorizationCtxKey).(*twomes.Authorization)
	if !ok {
		logrus.Info("failed when getting authentication context value")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	if !auth.IsKind(twomes.AccountActivationToken) {
		logrus.Info("%s token was used while %s was required", auth.Kind, twomes.AccountActivationToken)
		http.Error(w, "wrong token kind", http.StatusForbidden)
		return
	}

	account, err := h.accountService.Activate(auth.ID, request.Longtitude, request.Latitude, request.TZName)
	if err != nil {
		logrus.Info(err)

		if errors.Is(err, twomes.ErrAccountAlreadyActivated) {
			http.Error(w, "account already activated", http.StatusBadRequest)
			return
		}

		http.Error(w, "account activation failed", http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(account)
	if err != nil {
		logrus.Error(err)
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
