package shoppinglistitemtype

type ShoppingListItemTypeRepository interface {
	Find(shoppingListItemType ShoppingListItemType) (ShoppingListItemType, error)
	GetAll() ([]ShoppingListItemType, error)
	Create(ShoppingListItemType) (ShoppingListItemType, error)
	Delete(ShoppingListItemType) error
}
