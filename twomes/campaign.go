package twomes

import "time"

// A campaign is a timeframe where we gather measurements with a specific goal.
type Campaign struct {
	ID         uint         `json:"id"`
	Name       string       `json:"name"`
	App        App          `json:"app"`
	InfoURL    string       `json:"info_url"`
	CloudFeeds []*CloudFeed `json:"cloud_feeds"`
	StartTime  *time.Time   `json:"start_time,omitempty"`
	EndTime    *time.Time   `json:"end_time,omitempty"`
}

// Create a new Campaign.
func MakeCampaign(name string, app App, infoURL string, cloudFeeds []*CloudFeed, startTime, endTime *time.Time) Campaign {
	return Campaign{
		Name:      name,
		App:       app,
		InfoURL:   infoURL,
		StartTime: startTime,
		EndTime:   endTime,
	}
}
