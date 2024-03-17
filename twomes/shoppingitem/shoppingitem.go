package shoppingitem

// ItemType represents the type of shopping item
type ItemType string

const (
	Device      ItemType = "device"
	CloudFeed   ItemType = "cloudfeed"
	EnergyQuery ItemType = "energyquery"
)

// Interface to allow device, cloudfeed and energyquery in one action.
type ActionModel interface{}

// An item can be a device, cloudfeed or energyquery
type ShoppingItem struct {
	ID       uint        `json:"id"`
	ActionID ActionModel `json:"actionid"`
	Type     ItemType    `json:"type"`
}
