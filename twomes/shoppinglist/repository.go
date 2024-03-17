package shoppinglist

type ShoppingListRepository interface {
	Create(ShoppingList) (ShoppingList, error)
	Delete(ShoppingList) error
}
