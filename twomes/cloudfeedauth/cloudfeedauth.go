package cloudfeedauth

import "time"

// A CloudFeedAuth stores auth information about CloudFeeds authorized by an account.
type CloudFeedAuth struct {
	AccountID      uint      `json:"account_id"`
	CloudFeedID    uint      `json:"cloud_feed_id"`
	AccessToken    string    `json:"-"`
	RefreshToken   string    `json:"-"`
	Expiry         time.Time `json:"-"`
	AuthGrantToken string    `json:"auth_grant_token"`
}

// Create a new CloudFeedAuth.
func MakeCloudFeedAuth(accountID, cloudFeedID uint, accessToken string, refreshToken string, expiry time.Time, authGrantToken string) CloudFeedAuth {
	return CloudFeedAuth{
		AccountID:      accountID,
		CloudFeedID:    cloudFeedID,
		AccessToken:    accessToken,
		RefreshToken:   refreshToken,
		Expiry:         expiry,
		AuthGrantToken: authGrantToken,
	}
}
