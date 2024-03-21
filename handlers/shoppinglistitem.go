package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/services"
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglistitem"
	"github.com/sirupsen/logrus"
)

type ShoppingListItemHandler struct {
	service *services.ShoppingListItemService
}

// Create a new ShoppingListItemHandler.
func NewShoppingListItemHandler(service *services.ShoppingListItemService) *ShoppingListItemHandler {
	return &ShoppingListItemHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new ShoppingListItem.
func (h *ShoppingListItemHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request shoppinglistitem.ShoppingListItem
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	shoppinglistitem, err := h.service.Create(request.SourceID, request.Schedule, request.Type)

	if err != nil {
		if helpers.IsMySQLRecordNotFoundError(err) {
			return NewHandlerError(err, "not found", http.StatusNotFound)
		}

		if helpers.IsMySQLDuplicateError(err) {
			return NewHandlerError(err, "duplicate", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(shoppinglistitem)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}
