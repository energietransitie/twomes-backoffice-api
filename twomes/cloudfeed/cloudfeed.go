package cloudfeed

import "github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeedauth"

// A CloudFeed is an external online data source.
type CloudFeed struct {
	ID               uint                          `json:"id"`
	Name             string                        `json:"name"`
	AuthorizationURL string                        `json:"authorization_url"`
	TokenURL         string                        `json:"token_url"`
	ClientID         string                        `json:"client_id"`
	ClientSecret     string                        `json:"client_secret,omitempty"`
	Scope            string                        `json:"scope"`
	RedirectURL      string                        `json:"redirect_url"`
	CloudFeedAuths   []cloudfeedauth.CloudFeedAuth `json:"-"`
}

// Create a new CloudFeed.
func MakeCloudFeed(name, authorizationURL, tokenURL, clientID, clientSecret, scope, redirectURL string) CloudFeed {
	return CloudFeed{
		Name:             name,
		AuthorizationURL: authorizationURL,
		TokenURL:         tokenURL,
		ClientID:         clientID,
		ClientSecret:     clientSecret,
		Scope:            scope,
		RedirectURL:      redirectURL,
	}
}
