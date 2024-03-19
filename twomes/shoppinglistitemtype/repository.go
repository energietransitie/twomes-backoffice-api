package shoppinglistitemtype

type ShoppingListItemTypeRepository interface {
	Create(ShoppingListItemType) (ShoppingListItemType, error)
	Delete(ShoppingListItemType) error
}
