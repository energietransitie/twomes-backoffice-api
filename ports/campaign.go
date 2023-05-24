package ports

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

// A CampaignRepository can load, store and delete campaigns.
type CampaignRepository interface {
	Find(campaign twomes.Campaign) (twomes.Campaign, error)
	GetAll() ([]twomes.Campaign, error)
	Create(twomes.Campaign) (twomes.Campaign, error)
	Delete(twomes.Campaign) error
}

// CampaignService exposes all operations that can be performed on a [twomes.Campaign].
type CampaignService interface {
	Create(name string, app twomes.App, infoURL string, cloudFeeds []*twomes.CloudFeed, startTime, endTime *time.Time) (twomes.Campaign, error)
	Find(campaign twomes.Campaign) (twomes.Campaign, error)
	GetByID(id uint) (twomes.Campaign, error)
}
