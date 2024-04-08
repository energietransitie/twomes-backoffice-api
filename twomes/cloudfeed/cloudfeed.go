package cloudfeed

import "time"

// A CloudFeed stores auth information about CloudFeeds authorized by an account.
type CloudFeed struct {
	AccountID       uint       `json:"account_id"`
	CloudFeedTypeID uint       `json:"cloud_feed_id"`
	AccessToken     string     `json:"-"`
	RefreshToken    string     `json:"-"`
	Expiry          time.Time  `json:"-"`
	AuthGrantToken  string     `json:"auth_grant_token"`
	ActivatedAt     *time.Time `json:"activated_at"`
}

// Create a new CloudFeed.
func MakeCloudFeed(accountID, cloudFeedTypeID uint, accessToken string, refreshToken string, expiry time.Time, authGrantToken string) CloudFeed {
	now := time.Now().UTC()
	return CloudFeed{
		AccountID:       accountID,
		CloudFeedTypeID: cloudFeedTypeID,
		AccessToken:     accessToken,
		RefreshToken:    refreshToken,
		Expiry:          expiry,
		AuthGrantToken:  authGrantToken,
		ActivatedAt:     &now,
	}
}
