package twomes

// A CloudFeedAuthStatus contains all available cloud feeds for an account and if they are connected or not.
type CloudFeedAuthStatus struct {
	CloudFeed CloudFeed `json:"cloud_feed"`
	Connected bool      `json:"connected"`
}

// Create a new CloudFeedStatus.
func MakeCloudFeedAuthStatus(cloudFeed CloudFeed, cloudFeedAuth CloudFeedAuth) CloudFeedAuthStatus {
	// TODO: Change this to check if access and/or refresh token is valid.
	// For now, we use the grant token, since we have not implemented functionality
	// to exchange the auth grant token for an access and refresh token.
	connected := cloudFeedAuth.AuthGrantToken != ""

	return CloudFeedAuthStatus{
		CloudFeed: cloudFeed,
		Connected: connected,
	}
}
