package shoppinglistitem

import "github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglistitemtype"

// An item can be a device, cloudfeed or energyquery
type ShoppingListItem struct {
	ID                    uint                                      `json:"id"`
	SourceID              uint                                      `json:"source_id"`
	Type                  shoppinglistitemtype.ShoppingListItemType `json:"type"`
	Precedes              []ShoppingListItem                        `json:"precedes"`
	UploadSchedule        []string                                  `json:"upload_schedule"`
	MeasurementSchedule   []string                                  `json:"measurement_schedule"`
	NotificationThreshold string                                    `json:"notification_threshold"`
}

func MakeShoppingListItem(
	sourceID uint,
	itemType shoppinglistitemtype.ShoppingListItemType,
	precedes []ShoppingListItem,
	uploadSchedule []string,
	measurementSchedule []string,
	notificationThreshold string,
) ShoppingListItem {
	return ShoppingListItem{
		SourceID:              sourceID,
		Type:                  itemType,
		Precedes:              precedes,
		UploadSchedule:        uploadSchedule,
		MeasurementSchedule:   measurementSchedule,
		NotificationThreshold: notificationThreshold,
	}
}
