package shoppingitem

type ShoppingItemRepository interface {
	Create(ShoppingItem) (ShoppingItem, error)
	Delete(ShoppingItem) error
}
