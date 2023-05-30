package services

import (
	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

type CloudFeedService struct {
	repository ports.CloudFeedRepository
}

// Create a new CloudFeedService.
func NewCloudFeedService(repository ports.CloudFeedRepository) *CloudFeedService {
	return &CloudFeedService{
		repository: repository,
	}
}

func (s *CloudFeedService) Create(name string, authorizationURL string, tokenURL string, clientID string, clientSecret string, scope string, redirectURL string) (twomes.CloudFeed, error) {
	cloudFeed := twomes.MakeCloudFeed(name, authorizationURL, tokenURL, clientID, clientSecret, scope, redirectURL)
	return s.repository.Create(cloudFeed)
}

func (s *CloudFeedService) Find(cloudFeed twomes.CloudFeed) (twomes.CloudFeed, error) {
	return s.repository.Find(cloudFeed)
}

func (s *CloudFeedService) GetByID(id uint) (twomes.CloudFeed, error) {
	return s.repository.Find(twomes.CloudFeed{ID: id})
}
