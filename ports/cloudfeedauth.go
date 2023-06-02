package ports

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

// A CloudFeedAuthRepository can load, store and delete cloudFeedAuths.
type CloudFeedAuthRepository interface {
	Find(cloudFeedAuth twomes.CloudFeedAuth) (twomes.CloudFeedAuth, error)
	GetAll() ([]twomes.CloudFeedAuth, error)
	Create(twomes.CloudFeedAuth) (twomes.CloudFeedAuth, error)
	Delete(twomes.CloudFeedAuth) error
}

// CloudFeedAuthService exposes all operations that can be performed on a [twomes.CloudFeedAuth].
type CloudFeedAuthService interface {
	Create(accountID, cloudFeedID uint, authGrantToken string) (twomes.CloudFeedAuth, error)
	Find(cloudFeedAuth twomes.CloudFeedAuth) (twomes.CloudFeedAuth, error)
}
