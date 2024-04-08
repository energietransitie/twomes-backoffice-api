package account

import (
	"errors"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/building"
	"github.com/energietransitie/twomes-backoffice-api/twomes/campaign"
	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeed"
)

var (
	ErrAccountAlreadyActivated = errors.New("account is already activated")
)

// An Account is registered to a research subject.
type Account struct {
	ID                 uint                  `json:"id"`
	Campaign           campaign.Campaign     `json:"campaign"`
	ActivatedAt        *time.Time            `json:"activated_at"`
	InvitationToken    string                `json:"invitation_token,omitempty"`
	InvitationURL      string                `json:"invitation_url,omitempty"`
	AuthorizationToken string                `json:"authorization_token,omitempty"`
	Buildings          []building.Building   `json:"buildings,omitempty"`
	CloudFeeds         []cloudfeed.CloudFeed `json:"cloud_feeds,omitempty"`
	// Maybe use separate pseudonym field,
	// but right now we can derive a pseudonym
	// using the ID or the campaign ID + account ID.
}

// Create a new Account.
func MakeAccount(campaign campaign.Campaign) Account {
	return Account{
		Campaign: campaign,
	}
}

// Activate an account.
// An error will be returned if the account is already activated.
func (a *Account) Activate() error {
	if a.ActivatedAt != nil {
		return ErrAccountAlreadyActivated
	}

	now := time.Now().UTC()
	a.ActivatedAt = &now

	return nil
}
