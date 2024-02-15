package ports

import (
	"context"

	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeedauth"
)

// CloudFeedAuthService exposes all operations that can be performed on a [cloudfeedauth.CloudFeedAuth].
type CloudFeedAuthService interface {
	Create(ctx context.Context, accountID, cloudFeedID uint, authGrantToken string) (cloudfeedauth.CloudFeedAuth, error)
	Find(cloudFeedAuth cloudfeedauth.CloudFeedAuth) (cloudfeedauth.CloudFeedAuth, error)
}
