package shoppinglist

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppingitem"
)

// A shoppinglist is a collection of shoppingitems that are part of a campaign.
type ShoppingList struct {
	ID           uint                           `json:"id"`
	Items        []shoppingitem.ShoppingItem    `json:"items,omitempty"`
	Dependencies [][2]shoppingitem.ShoppingItem `json:"dependencies"` //Example: [[1,2],[3,2]], 1 and 3 should be done before 2
}
