package uploadschedule

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppingitem"
)

// A shoppinglist is a collection of shoppingitems that are part of a campaign.
type UploadSchedule struct {
	ID       uint                      `json:"id"`
	Item     shoppingitem.ShoppingItem `json:"item"`
	Schedule []string                  `json:"schedule"` //Schedule in cronjob format, can be multiple
}
