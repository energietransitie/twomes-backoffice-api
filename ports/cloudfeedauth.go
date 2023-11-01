package ports

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

// A CloudFeedAuthRepository can load, store and delete cloudFeedAuths.
type CloudFeedAuthRepository interface {
	Find(cloudFeedAuth twomes.CloudFeedAuth) (twomes.CloudFeedAuth, error)
	FindOAuthInfo(accountID uint, cloudFeedID uint) (tokenURL string, refreshToken string, clientID string, clientSecret string, err error)
	FindFirstTokenToExpire() (accountID uint, cloudFeedID uint, expiry time.Time, err error)
	GetAll() ([]twomes.CloudFeedAuth, error)
	Create(twomes.CloudFeedAuth) (twomes.CloudFeedAuth, error)
	Update(twomes.CloudFeedAuth) (twomes.CloudFeedAuth, error)
	Delete(twomes.CloudFeedAuth) error
}

// CloudFeedAuthService exposes all operations that can be performed on a [twomes.CloudFeedAuth].
type CloudFeedAuthService interface {
	Create(accountID, cloudFeedID uint, authGrantToken string) (twomes.CloudFeedAuth, error)
	Find(cloudFeedAuth twomes.CloudFeedAuth) (twomes.CloudFeedAuth, error)
}
