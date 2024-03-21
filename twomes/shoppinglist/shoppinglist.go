package shoppinglist

import "github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglistitem"

// A shoppinglist is a collection of shoppingitems that are part of a campaign.
type ShoppingList struct {
	ID           uint                                   `json:"id"`
	Items        []shoppinglistitem.ShoppingListItem    `json:"items,omitempty"`
	Dependencies [][2]shoppinglistitem.ShoppingListItem `json:"dependencies"` //Example: [[1,2],[3,2]], 1 and 3 should be done before 2
}

// Create a new ShoppingList.
func MakeShoppingList(Items []shoppinglistitem.ShoppingListItem, Dependencies [][2]shoppinglistitem.ShoppingListItem) ShoppingList {
	return ShoppingList{
		Items:        Items,
		Dependencies: Dependencies,
	}
}
