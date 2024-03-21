package shoppinglist

type ShoppingListRepository interface {
	Find(shoppingList ShoppingList) (ShoppingList, error)
	GetAll() ([]ShoppingList, error)
	Create(ShoppingList) (ShoppingList, error)
	Delete(ShoppingList) error
}
