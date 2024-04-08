package cloudfeedstatus

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeed"
	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeedtype"
)

// A CloudFeedStatus contains all available cloud feeds for an account and if they are connected or not.
type CloudFeedStatus struct {
	CloudFeed cloudfeed.CloudFeed `json:"cloud_feed"`
	Connected bool                `json:"connected"`
}

// Create a new CloudFeedStatus.
func MakeCloudFeedStatus(cloudFeedType cloudfeedtype.CloudFeedType, cloudFeed cloudfeed.CloudFeed) CloudFeedStatus {
	// TODO: Change this to check if access and/or refresh token is valid.
	// For now, we use the grant token, since we have not implemented functionality
	// to exchange the  grant token for an access and refresh token.
	connected := cloudFeed.AuthGrantToken != ""

	return CloudFeedStatus{
		CloudFeed: cloudFeed,
		Connected: connected,
	}
}
