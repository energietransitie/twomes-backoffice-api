package cloudfeed

// A CloudFeedRepository can load, store and delete cloud feeds.
type CloudFeedRepository interface {
	Find(CloudFeed) (CloudFeed, error)
	GetAll() ([]CloudFeed, error)
	Create(CloudFeed) (CloudFeed, error)
	Delete(CloudFeed) error
}
