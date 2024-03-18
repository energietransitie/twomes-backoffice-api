package shoppingitem

import "github.com/energietransitie/twomes-backoffice-api/twomes/shoppingitemtype"

// Interface to allow device, cloudfeed and energyquery in one action.
type ActionModel interface{}

// An item can be a device, cloudfeed or energyquery
type ShoppingItem struct {
	ID       uint                              `json:"id"`
	ActionID ActionModel                       `json:"actionid"`
	Type     shoppingitemtype.ShoppingItemType `json:"type"`
}
