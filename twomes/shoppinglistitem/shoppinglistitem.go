package shoppinglistitem

import "github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglistitemtype"

// An item can be a device, cloudfeed or energyquery
type ShoppingListItem struct {
	ID       uint                                      `json:"id"`
	SourceID uint                                      `json:"source_id"`
	Type     shoppinglistitemtype.ShoppingListItemType `json:"type"`
	Schedule []string                                  `json:"schedule"`
}
