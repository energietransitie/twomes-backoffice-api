package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/services"
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglistitemtype"
	"github.com/sirupsen/logrus"
)

type ShoppingListItemTypeHandler struct {
	service *services.ShoppingListItemTypeService
}

// Create a new ShoppingListItemTypeHandler.
func NewShoppingListItemTypeHandler(service *services.ShoppingListItemTypeService) *ShoppingListItemTypeHandler {
	return &ShoppingListItemTypeHandler{
		service: service,
	}
}

// Handle API endpoint for creating a new ShoppingListItemType.
func (h *ShoppingListItemTypeHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var request shoppinglistitemtype.ShoppingListItemType
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		return NewHandlerError(err, "bad request", http.StatusBadRequest).WithLevel(logrus.ErrorLevel)
	}

	shoppinglistitemtype, err := h.service.Create(request.Name)

	if err != nil {
		if helpers.IsMySQLRecordNotFoundError(err) {
			return NewHandlerError(err, "not found", http.StatusNotFound)
		}

		if helpers.IsMySQLDuplicateError(err) {
			return NewHandlerError(err, "duplicate", http.StatusBadRequest)
		}

		return NewHandlerError(err, "internal server error", http.StatusInternalServerError)
	}

	err = json.NewEncoder(w).Encode(shoppinglistitemtype)
	if err != nil {
		return NewHandlerError(err, "internal server error", http.StatusInternalServerError).WithLevel(logrus.ErrorLevel)
	}

	return nil
}
