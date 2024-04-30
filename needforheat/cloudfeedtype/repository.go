package cloudfeedtype

// A CloudFeedTypeRepository can load, store and delete cloud feeds.
type CloudFeedTypeRepository interface {
	Find(CloudFeedType) (CloudFeedType, error)
	GetAll() ([]CloudFeedType, error)
	Create(CloudFeedType) (CloudFeedType, error)
	Delete(CloudFeedType) error
}
