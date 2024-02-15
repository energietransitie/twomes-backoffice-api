package cloudfeedauth

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/device"
)

// A CloudFeedAuthRepository can load, store and delete cloudFeedAuths.
type CloudFeedAuthRepository interface {
	Find(cloudFeedAuth CloudFeedAuth) (CloudFeedAuth, error)
	FindOAuthInfo(accountID uint, cloudFeedID uint) (tokenURL string, refreshToken string, clientID string, clientSecret string, err error)
	FindFirstTokenToExpire() (accountID uint, cloudFeedID uint, expiry time.Time, err error)
	FindDevice(cloudFeedAuth CloudFeedAuth) (*device.Device, error)
	GetAll() ([]CloudFeedAuth, error)
	Create(CloudFeedAuth) (CloudFeedAuth, error)
	Update(CloudFeedAuth) (CloudFeedAuth, error)
	Delete(CloudFeedAuth) error
}
