package ports

import (
	"context"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeedauth"
	"github.com/energietransitie/twomes-backoffice-api/twomes/device"
)

// A CloudFeedAuthRepository can load, store and delete cloudFeedAuths.
type CloudFeedAuthRepository interface {
	Find(cloudFeedAuth cloudfeedauth.CloudFeedAuth) (cloudfeedauth.CloudFeedAuth, error)
	FindOAuthInfo(accountID uint, cloudFeedID uint) (tokenURL string, refreshToken string, clientID string, clientSecret string, err error)
	FindFirstTokenToExpire() (accountID uint, cloudFeedID uint, expiry time.Time, err error)
	FindDevice(cloudFeedAuth cloudfeedauth.CloudFeedAuth) (*device.Device, error)
	GetAll() ([]cloudfeedauth.CloudFeedAuth, error)
	Create(cloudfeedauth.CloudFeedAuth) (cloudfeedauth.CloudFeedAuth, error)
	Update(cloudfeedauth.CloudFeedAuth) (cloudfeedauth.CloudFeedAuth, error)
	Delete(cloudfeedauth.CloudFeedAuth) error
}

// CloudFeedAuthService exposes all operations that can be performed on a [cloudfeedauth.CloudFeedAuth].
type CloudFeedAuthService interface {
	Create(ctx context.Context, accountID, cloudFeedID uint, authGrantToken string) (cloudfeedauth.CloudFeedAuth, error)
	Find(cloudFeedAuth cloudfeedauth.CloudFeedAuth) (cloudfeedauth.CloudFeedAuth, error)
}
