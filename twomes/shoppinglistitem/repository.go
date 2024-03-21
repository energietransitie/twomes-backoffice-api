package shoppinglistitem

type ShoppingListItemRepository interface {
	Find(shoppingListItem ShoppingListItem) (ShoppingListItem, error)
	GetAll() ([]ShoppingListItem, error)
	Create(ShoppingListItem) (ShoppingListItem, error)
	Delete(ShoppingListItem) error
}
