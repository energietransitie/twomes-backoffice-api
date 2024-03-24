package shoppinglist

import "github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglistitem"

// A shoppinglist is a collection of shoppingitems that are part of a campaign.
type ShoppingList struct {
	ID          uint                                `json:"id"`
	Description string                              `json:"description"`
	Items       []shoppinglistitem.ShoppingListItem `json:"items,omitempty"`
}

// Create a new ShoppingList.
func MakeShoppingList(items []shoppinglistitem.ShoppingListItem, description string) ShoppingList {
	return ShoppingList{
		Items:       items,
		Description: description,
	}
}
