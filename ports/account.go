// Package ports exposes ports for interacting with business logic.
package ports

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/account"
	"github.com/energietransitie/twomes-backoffice-api/twomes/campaign"
	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeedauthstatus"
)

// An AccountRepository can load, store and delete accounts.
type AccountRepository interface {
	Find(account account.Account) (account.Account, error)
	GetAll() ([]account.Account, error)
	Create(account.Account) (account.Account, error)
	Update(account.Account) (account.Account, error)
	Delete(account.Account) error
}

// An AccountService exposes all operations we can perform on a [account.Account]
type AccountService interface {
	Create(campaign campaign.Campaign) (account.Account, error)
	Activate(id uint, longtitude, latitude float32, tzName string) (account.Account, error)
	GetByID(id uint) (account.Account, error)
	GetCloudFeedAuthStatuses(id uint) ([]cloudfeedauthstatus.CloudFeedAuthStatus, error)
}
