package campaign

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/app"
	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeed"
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglist"
	"github.com/energietransitie/twomes-backoffice-api/twomes/uploadschedule"
)

// A campaign is a timeframe where we gather measurements with a specific goal.
type Campaign struct {
	ID             uint                          `json:"id"`
	Name           string                        `json:"name"`
	App            app.App                       `json:"app"`
	InfoURL        string                        `json:"info_url"`
	CloudFeeds     []cloudfeed.CloudFeed         `json:"cloud_feeds"`
	StartTime      *time.Time                    `json:"start_time,omitempty"`
	EndTime        *time.Time                    `json:"end_time,omitempty"`
	ShoppingList   shoppinglist.ShoppingList     `json:"shoppinglist,omitempty"`
	UploadSchedule uploadschedule.UploadSchedule `json:"uploadschedule"`
}

// Create a new Campaign.
func MakeCampaign(name string, app app.App, infoURL string, cloudFeeds []cloudfeed.CloudFeed, startTime, endTime *time.Time) Campaign {
	return Campaign{
		Name:       name,
		App:        app,
		InfoURL:    infoURL,
		CloudFeeds: cloudFeeds,
		StartTime:  startTime,
		EndTime:    endTime,
	}
}
