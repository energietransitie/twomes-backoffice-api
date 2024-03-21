package shoppinglistitem

import "github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglistitemtype"

// An item can be a device, cloudfeed or energyquery
type ShoppingListItem struct {
	ID       uint                                      `json:"id"`
	SourceID uint                                      `json:"source_id"`
	Schedule []string                                  `json:"schedule"`
	Type     shoppinglistitemtype.ShoppingListItemType `json:"type"`
}

func MakeShoppingListItem(SourceID uint, Schedule []string, Type shoppinglistitemtype.ShoppingListItemType) ShoppingListItem {
	return ShoppingListItem{
		SourceID: SourceID,
		Schedule: Schedule,
		Type:     Type,
	}
}
