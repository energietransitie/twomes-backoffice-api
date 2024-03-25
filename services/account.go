// Package services exposes a services as entrypoints for business logic.
package services

import (
	"errors"
	"net/url"
	"regexp"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/twomes/account"
	"github.com/energietransitie/twomes-backoffice-api/twomes/authorization"
	"github.com/energietransitie/twomes-backoffice-api/twomes/campaign"
	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeedauth"
	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeedauthstatus"
	"github.com/sirupsen/logrus"
)

var (
	ErrTokenSigningMethodInvalid = errors.New("unexpected signing method")
	ErrTokenInvalid              = errors.New("token is invalid")
)

type AccountService struct {
	repository account.AccountRepository

	// Services used when activating an account.
	authService     *AuthorizationService
	appService      *AppService
	campaignService *CampaignService
	buildingService *BuildingService

	// Services used for getting cloud feed auth statuses.
	cloudFeedAuthService *CloudFeedAuthService

	// Regular expression used for pattern matching in a provisioning_url_template.
	activationTokenRegex *regexp.Regexp
}

// Create a new AccountService
func NewAccountService(repository account.AccountRepository, authService *AuthorizationService, appService *AppService, campaignService *CampaignService, buildingService *BuildingService, cloudFeedAuthService *CloudFeedAuthService) *AccountService {
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
func (s *AccountService) Create(campaign campaign.Campaign) (account.Account, error) {
	campaign, err := s.campaignService.Find(campaign)
	if err != nil {
		return account.Account{}, err
	}

	a := account.MakeAccount(campaign)

	a, err = s.repository.Create(a)
	if err != nil {
		return account.Account{}, err
	}

	a.InvitationToken, err = s.authService.CreateToken(authorization.AccountActivationToken, a.ID, time.Time{})
	if err != nil {
		return account.Account{}, err
	}

	a.InvitationURL = s.activationTokenRegex.ReplaceAllString(campaign.App.ProvisioningURLTemplate, url.PathEscape(a.InvitationToken))

	return a, nil
}

// Activate an account.
func (s *AccountService) Activate(id uint, longitude, latitude float32, tzName string) (account.Account, error) {
	a, err := s.repository.Find(account.Account{ID: id})
	if err != nil {
		return account.Account{}, err
	}

	err = a.Activate()
	if err != nil {
		return a, err
	}

	a, err = s.repository.Update(a)
	if err != nil {
		return account.Account{}, err
	}

	if len(a.Buildings) < 1 {
		building, err := s.buildingService.Create(a.ID, longitude, latitude, tzName)
		if err != nil {
			return account.Account{}, err
		}

		a.Buildings = append(a.Buildings, building)
	}

	a.AuthorizationToken, err = s.authService.CreateToken(authorization.AccountToken, a.ID, time.Time{})
	if err != nil {
		return account.Account{}, err
	}

	return a, nil
}

// Get an account by ID.
func (s *AccountService) GetByID(id uint) (account.Account, error) {
	return s.repository.Find(account.Account{ID: id})
}

// Get cloud feed auth statuses.
func (s *AccountService) GetCloudFeedAuthStatuses(id uint) ([]cloudfeedauthstatus.CloudFeedAuthStatus, error) {
	var cloudFeedAuthStatuses []cloudfeedauthstatus.CloudFeedAuthStatus

	a, err := s.GetByID(id)
	if err != nil {
		return cloudFeedAuthStatuses, err
	}

	for _, cloudFeed := range a.Campaign.CloudFeeds {
		cloudFeedAuth, err := s.cloudFeedAuthService.Find(cloudfeedauth.CloudFeedAuth{AccountID: id, CloudFeedID: cloudFeed.ID})
		if err != nil && !helpers.IsMySQLRecordNotFoundError(err) {
			return cloudFeedAuthStatuses, err
		}

		cloudFeedAuthStatuses = append(cloudFeedAuthStatuses, cloudfeedauthstatus.MakeCloudFeedAuthStatus(cloudFeed, cloudFeedAuth))
	}

	return cloudFeedAuthStatuses, nil
}
