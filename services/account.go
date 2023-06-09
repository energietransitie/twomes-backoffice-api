// Package services exposes a services as entrypoints for business logic.
package services

import (
	"errors"
	"net/url"
	"regexp"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
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

	// Services used for getting cloud feed auth statuses.
	cloudFeedAuthService ports.CloudFeedAuthService

	// Regular expression used for pattern matching in a provisioning_url_template.
	activationTokenRegex *regexp.Regexp
}

// Create a new AccountService
func NewAccountService(repository ports.AccountRepository, authService ports.AuthorizationService, appService ports.AppService, campaignService ports.CampaignService, buildingService ports.BuildingService, cloudFeedAuthService ports.CloudFeedAuthService) *AccountService {
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
		cloudFeedAuthService: cloudFeedAuthService,
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
func (s *AccountService) Activate(id uint, longitude, latitude float32, tzName string) (twomes.Account, error) {
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

	if len(account.Buildings) < 1 {
		building, err := s.buildingService.Create(account.ID, longitude, latitude, tzName)
		if err != nil {
			return twomes.Account{}, err
		}

		account.Buildings = append(account.Buildings, building)
	}

	account.AuthorizationToken, err = s.authService.CreateToken(twomes.AccountToken, account.ID, time.Time{})
	if err != nil {
		return twomes.Account{}, err
	}

	return account, nil
}

// Get an account by ID.
func (s *AccountService) GetByID(id uint) (twomes.Account, error) {
	return s.repository.Find(twomes.Account{ID: id})
}

// Get cloud feed auth statuses.
func (s *AccountService) GetCloudFeedAuthStatuses(id uint) ([]twomes.CloudFeedAuthStatus, error) {
	var cloudFeedAuthStatuses []twomes.CloudFeedAuthStatus

	account, err := s.GetByID(id)
	if err != nil {
		return cloudFeedAuthStatuses, err
	}

	for _, cloudFeed := range account.Campaign.CloudFeeds {
		cloudFeedAuth, err := s.cloudFeedAuthService.Find(twomes.CloudFeedAuth{AccountID: id, CloudFeedID: cloudFeed.ID})
		if err != nil && !helpers.IsMySQLRecordNotFoundError(err) {
			return cloudFeedAuthStatuses, err
		}

		cloudFeedAuthStatuses = append(cloudFeedAuthStatuses, twomes.MakeCloudFeedAuthStatus(cloudFeed, cloudFeedAuth))
	}

	return cloudFeedAuthStatuses, nil
}
