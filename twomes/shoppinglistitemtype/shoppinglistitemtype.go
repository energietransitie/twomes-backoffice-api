package shoppinglistitemtype

// A type can be a device, cloudfeed or energyquery
type ShoppingListItemType struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
