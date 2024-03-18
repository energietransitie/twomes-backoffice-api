package shoppingitemtype

type ShoppingItemTypeRepository interface {
	Create(ShoppingItemType) (ShoppingItemType, error)
	Delete(ShoppingItemType) error
}
