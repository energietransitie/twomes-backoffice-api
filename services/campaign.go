package services

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/app"
	"github.com/energietransitie/twomes-backoffice-api/twomes/campaign"
	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeed"
)

type CampaignService struct {
	repository campaign.CampaignRepository

	// Service used when creating a campaign.
	appService       *AppService
	cloudFeedService *CloudFeedService
}

// Create a new CampaignService.
func NewCampaignService(repository campaign.CampaignRepository, appService *AppService, cloudFeedService *CloudFeedService) *CampaignService {
	return &CampaignService{
		repository:       repository,
		appService:       appService,
		cloudFeedService: cloudFeedService,
	}
}

// Create a new campaign.
func (s *CampaignService) Create(name string, app app.App, infoURL string, cloudFeeds []cloudfeed.CloudFeed, startTime, endTime *time.Time) (campaign.Campaign, error) {
	app, err := s.appService.Find(app)
	if err != nil {
		return campaign.Campaign{}, err
	}

	for i, cloudFeed := range cloudFeeds {
		cloudFeeds[i], err = s.cloudFeedService.Find(cloudFeed)
		if err != nil {
			return campaign.Campaign{}, err
		}
	}

	campaign := campaign.MakeCampaign(name, app, infoURL, cloudFeeds, startTime, endTime)
	return s.repository.Create(campaign)
}

// Find a campaign using any field set in the campaign struct.
func (s *CampaignService) Find(campaign campaign.Campaign) (campaign.Campaign, error) {
	return s.repository.Find(campaign)
}

// Get a campaign by its ID.
func (s *CampaignService) GetByID(id uint) (campaign.Campaign, error) {
	return s.repository.Find(campaign.Campaign{ID: id})
}
