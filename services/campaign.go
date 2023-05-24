package services

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

type CampaignService struct {
	repository ports.CampaignRepository

	// Service used when creating a campaign.
	appService ports.AppService
}

// Create a new CampaignService.
func NewCampaignService(repository ports.CampaignRepository, appService ports.AppService) *CampaignService {
	return &CampaignService{
		repository: repository,
		appService: appService,
	}
}

// Create a new campaign.
func (s *CampaignService) Create(name string, app twomes.App, infoURL string, cloudFeeds []*twomes.CloudFeed, startTime, endTime *time.Time) (twomes.Campaign, error) {
	app, err := s.appService.Find(app)
	if err != nil {
		return twomes.Campaign{}, err
	}

	campaign := twomes.MakeCampaign(name, app, infoURL, cloudFeeds, startTime, endTime)
	return s.repository.Create(campaign)
}

// Find a campaign using any field set in the campaign struct.
func (s *CampaignService) Find(campaign twomes.Campaign) (twomes.Campaign, error) {
	return s.repository.Find(campaign)
}

// Get a campaign by its ID.
func (s *CampaignService) GetByID(id uint) (twomes.Campaign, error) {
	return s.repository.Find(twomes.Campaign{ID: id})
}
