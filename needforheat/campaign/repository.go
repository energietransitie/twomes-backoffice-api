package campaign

// A CampaignRepository can load, store and delete campaigns.
type CampaignRepository interface {
	Find(campaign Campaign) (Campaign, error)
	GetAll() ([]Campaign, error)
	Create(Campaign) (Campaign, error)
	Delete(Campaign) error
}
