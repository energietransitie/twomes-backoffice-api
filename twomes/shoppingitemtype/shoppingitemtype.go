package shoppingitemtype

// A type can be a device, cloudfeed or energyquery
type ShoppingItemType struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
