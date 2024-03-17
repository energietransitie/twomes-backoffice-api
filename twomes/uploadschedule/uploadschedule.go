package uploadschedule

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppingitem"
)

// A shoppinglist is a collection of shoppingitems that are part of a campaign.
type UploadSchedule struct {
	ID       uint                      `json:"id"`
	Item     shoppingitem.ShoppingItem `json:"item"`
	Schedule []time.Time               `json:"schedule"` //Schedule in seconds. So for every 10, 30, 60 minutes. you'd need to set it to [60, 1800, 3600]
}
