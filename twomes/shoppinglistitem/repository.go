package shoppinglistitem

type ShoppingListItemRepository interface {
	Create(ShoppingListItem) (ShoppingListItem, error)
	Delete(ShoppingListItem) error
}
