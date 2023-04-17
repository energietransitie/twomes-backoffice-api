// Package services exposes a services as entrypoints for business logic.
package services

import (
	"errors"
	"net/url"
	"regexp"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"github.com/sirupsen/logrus"
)

var (
	ErrTokenSigningMethodInvalid = errors.New("unexpected signing method")
	ErrTokenInvalid              = errors.New("token is invalid")
)

type AccountService struct {
	repository ports.AccountRepository

	// Services used when activating an account.
	authService     ports.AuthorizationService
	appService      ports.AppService
	campaignService ports.CampaignService
	buildingService ports.BuildingService

	// Regular expression used for pattern matching in a provisioning_url_template.
	activationTokenRegex *regexp.Regexp
}

// Create a new AccountService
func NewAccountService(repository ports.AccountRepository, authService ports.AuthorizationService, appService ports.AppService, campaignService ports.CampaignService, buildingService ports.BuildingService) *AccountService {
	activationTokenRegex, err := regexp.Compile(`<account_activation_token>`)
	if err != nil {
		logrus.WithField("error", err).Fatal("account activation token regex did not compile")
	}

	return &AccountService{
		repository:           repository,
		authService:          authService,
		appService:           appService,
		campaignService:      campaignService,
		buildingService:      buildingService,
		activationTokenRegex: activationTokenRegex,
	}
}

// Create a new account.
func (s *AccountService) Create(campaign twomes.Campaign) (twomes.Account, error) {
	campaign, err := s.campaignService.Find(campaign)
	if err != nil {
		return twomes.Account{}, err
	}

	account := twomes.MakeAccount(campaign)

	account, err = s.repository.Create(account)
	if err != nil {
		return twomes.Account{}, err
	}

	account.InvitationToken, err = s.authService.CreateToken(twomes.AccountActivationToken, account.ID, time.Time{})
	if err != nil {
		return twomes.Account{}, err
	}

	account.InvitationURL = s.activationTokenRegex.ReplaceAllString(campaign.App.ProvisioningURLTemplate, url.PathEscape(account.InvitationToken))

	return account, nil
}

// Activate an account.
func (s *AccountService) Activate(id uint, longtitude, latitude float32, tzName string) (twomes.Account, error) {
	account, err := s.repository.Find(twomes.Account{ID: id})
	if err != nil {
		return twomes.Account{}, err
	}

	err = account.Activate()
	if err != nil {
		return account, err
	}

	account, err = s.repository.Update(account)
	if err != nil {
		return twomes.Account{}, err
	}

	building, err := s.buildingService.Create(account.ID, longtitude, latitude, tzName)
	if err != nil {
		return twomes.Account{}, err
	}

	account.Buildings = append(account.Buildings, building)

	account.AuthorizationToken, err = s.authService.CreateToken(twomes.AccountToken, account.ID, time.Time{})
	if err != nil {
		return twomes.Account{}, err
	}

	return account, nil
}
