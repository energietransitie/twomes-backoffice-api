package ports

import "github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeed"

// CloudFeedService exposes all operations that can be performed on a [cloudfeed.CloudFeed].
type CloudFeedService interface {
	Create(name, authorizationURL, tokenURL, clientID, clientSecret, scope, redirectURL string) (cloudfeed.CloudFeed, error)
	Find(cloudFeed cloudfeed.CloudFeed) (cloudfeed.CloudFeed, error)
	GetByID(id uint) (cloudfeed.CloudFeed, error)
}
