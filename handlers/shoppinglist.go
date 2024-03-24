package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/services"
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglist"
	"github.com/sirupsen/logrus"
)

type ShoppingListHandler struct {
	service *services.ShoppingListService
}

// Create a new ShoppingListHandler.
func NewShoppingListHandler(service *services.ShoppingListService) *ShoppingListHandler {
	return &ShoppingListHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new ShoppingList.
func (h *ShoppingListHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request shoppinglist.ShoppingList
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	shoppinglist, err := h.service.Create(request.Description, request.Items)

	if err != nil {
		if helpers.IsMySQLRecordNotFoundError(err) {
			return NewHandlerError(err, "not found", http.StatusNotFound)
		}

		if helpers.IsMySQLDuplicateError(err) {
			return NewHandlerError(err, "duplicate", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(shoppinglist)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}
