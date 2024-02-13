package ports

import "github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeed"

// A CloudFeedRepository can load, store and delete cloud feeds.
type CloudFeedRepository interface {
	Find(cloudfeed.CloudFeed) (cloudfeed.CloudFeed, error)
	GetAll() ([]cloudfeed.CloudFeed, error)
	Create(cloudfeed.CloudFeed) (cloudfeed.CloudFeed, error)
	Delete(cloudfeed.CloudFeed) error
}

// CloudFeedService exposes all operations that can be performed on a [cloudfeed.CloudFeed].
type CloudFeedService interface {
	Create(name, authorizationURL, tokenURL, clientID, clientSecret, scope, redirectURL string) (cloudfeed.CloudFeed, error)
	Find(cloudFeed cloudfeed.CloudFeed) (cloudfeed.CloudFeed, error)
	GetByID(id uint) (cloudfeed.CloudFeed, error)
}
