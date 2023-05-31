package services

import (
	"errors"

	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

var (
	ErrDuplicateCloudFeedAuth = errors.New("duplicate cloud feed auth")
)

type CloudFeedAuthService struct {
	repository ports.CloudFeedAuthRepository
}

// Create a new CloudFeedAuthService.
func NewCloudFeedAuthService(repository ports.CloudFeedAuthRepository) *CloudFeedAuthService {
	return &CloudFeedAuthService{
		repository: repository,
	}
}

// Create a new cloudFeedAuth.
func (s *CloudFeedAuthService) Create(accountID, cloudFeedID uint, authGrantToken string) (twomes.CloudFeedAuth, error) {
	cloudFeedAuth := twomes.MakeCloudFeedAuth(accountID, cloudFeedID, authGrantToken)
	return s.repository.Create(cloudFeedAuth)
}

// Find a cloudFeedAuth using any field set in the cloudFeedAuth struct.
func (s *CloudFeedAuthService) Find(cloudFeedAuth twomes.CloudFeedAuth) (twomes.CloudFeedAuth, error) {
	return s.repository.Find(cloudFeedAuth)
}
