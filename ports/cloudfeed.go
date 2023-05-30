package ports

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

// A CloudFeedRepository can load, store and delete cloud feeds.
type CloudFeedRepository interface {
	Find(campaign twomes.CloudFeed) (twomes.CloudFeed, error)
	GetAll() ([]twomes.CloudFeed, error)
	Create(twomes.CloudFeed) (twomes.CloudFeed, error)
	Delete(twomes.CloudFeed) error
}

// CloudFeedService exposes all operations that can be performed on a [twomes.CloudFeed].
type CloudFeedService interface {
	Create(name, authorizationURL, tokenURL, clientID, clientSecret, scope, redirectURL string) (twomes.CloudFeed, error)
	Find(cloudFeed twomes.CloudFeed) (twomes.CloudFeed, error)
	GetByID(id uint) (twomes.CloudFeed, error)
}
