package cloudfeed

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/device"
)

// A CloudFeedRepository can load, store and delete CloudFeeds.
type CloudFeedRepository interface {
	Find(CloudFeed CloudFeed) (CloudFeed, error)
	FindOAuthInfo(accountID uint, cloudFeedID uint) (tokenURL string, refreshToken string, clientID string, clientSecret string, err error)
	FindFirstTokenToExpire() (accountID uint, cloudFeedID uint, expiry time.Time, err error)
	FindDevice(CloudFeed CloudFeed) (*device.Device, error)
	GetAll() ([]CloudFeed, error)
	Create(CloudFeed) (CloudFeed, error)
	Update(CloudFeed) (CloudFeed, error)
	Delete(CloudFeed) error
}
