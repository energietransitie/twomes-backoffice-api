// Package ports exposes ports for interacting with business logic.
package ports

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/account"
	"github.com/energietransitie/twomes-backoffice-api/twomes/campaign"
	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeedauthstatus"
)



// An AccountService exposes all operations we can perform on a [account.Account]
type AccountService interface {
	Create(campaign campaign.Campaign) (account.Account, error)
	Activate(id uint, longtitude, latitude float32, tzName string) (account.Account, error)
	GetByID(id uint) (account.Account, error)
	GetCloudFeedAuthStatuses(id uint) ([]cloudfeedauthstatus.CloudFeedAuthStatus, error)
}
