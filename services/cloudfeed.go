package services

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeed"
)

type CloudFeedService struct {
	repository cloudfeed.CloudFeedRepository
}

// Create a new CloudFeedService.
func NewCloudFeedService(repository cloudfeed.CloudFeedRepository) *CloudFeedService {
	return &CloudFeedService{
		repository: repository,
	}
}

func (s *CloudFeedService) Create(name string, authorizationURL string, tokenURL string, clientID string, clientSecret string, scope string, redirectURL string) (cloudfeed.CloudFeed, error) {
	cloudFeed := cloudfeed.MakeCloudFeed(name, authorizationURL, tokenURL, clientID, clientSecret, scope, redirectURL)
	return s.repository.Create(cloudFeed)
}

func (s *CloudFeedService) Find(cloudFeed cloudfeed.CloudFeed) (cloudfeed.CloudFeed, error) {
	return s.repository.Find(cloudFeed)
}

func (s *CloudFeedService) GetByID(id uint) (cloudfeed.CloudFeed, error) {
	return s.repository.Find(cloudfeed.CloudFeed{ID: id})
}

func (s *CloudFeedService) GetByIDForShoppingList(id uint) (interface{}, error) {
	return s.repository.Find(cloudfeed.CloudFeed{ID: id})
}

func (s *CloudFeedService) GetTableName() string {
	return "cloud_feed"
}
