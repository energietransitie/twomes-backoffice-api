package shoppinglist

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppingitem"
)

// A shoppinglist is a collection of shoppingitems that are part of a campaign.
type ShoppingList struct {
	ID    uint                        `json:"id"`
	Items []shoppingitem.ShoppingItem `json:"items,omitempty"`
}
