package services

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"golang.org/x/oauth2"
)

var (
	ErrDuplicateCloudFeedAuth = errors.New("duplicate cloud feed auth")
)

type CloudFeedAuthService struct {
	cloudFeedAuthRepo ports.CloudFeedAuthRepository
	cloudFeedRepo     ports.CloudFeedRepository
	ctx               context.Context
}

// Create a new CloudFeedAuthService.
func NewCloudFeedAuthService(ctx context.Context, cloudFeedAuthRepo ports.CloudFeedAuthRepository, cloudFeedRepo ports.CloudFeedRepository) *CloudFeedAuthService {
	return &CloudFeedAuthService{
		cloudFeedAuthRepo: cloudFeedAuthRepo,
		cloudFeedRepo:     cloudFeedRepo,
		ctx:               ctx,
	}
}

// Create a new cloudFeedAuth.
// This function exchanges the AuthGrantToken (Code) for a access and refresh token.
func (s *CloudFeedAuthService) Create(accountID, cloudFeedID uint, authGrantToken string) (twomes.CloudFeedAuth, error) {
	cloudFeed, err := s.cloudFeedRepo.Find(twomes.CloudFeed{ID: cloudFeedID})
	if err != nil {
		return twomes.CloudFeedAuth{}, err
	}

	scopes := strings.Split(cloudFeed.Scope, " ")

	conf := &oauth2.Config{
		ClientID:     cloudFeed.ClientID,
		ClientSecret: cloudFeed.ClientSecret,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  cloudFeed.AuthorizationURL,
			TokenURL: cloudFeed.TokenURL,
		},
		RedirectURL: cloudFeed.RedirectURL,
	}

	accessToken, refreshToken, expiry, err := exchangeAuthCode(s.ctx, conf, authGrantToken)
	if err != nil {
		return twomes.CloudFeedAuth{}, err
	}

	cloudFeedAuth := twomes.MakeCloudFeedAuth(accountID, cloudFeedID, accessToken, refreshToken, expiry, authGrantToken)
	return s.cloudFeedAuthRepo.Create(cloudFeedAuth)
}

// Find a cloudFeedAuth using any field set in the cloudFeedAuth struct.
func (s *CloudFeedAuthService) Find(cloudFeedAuth twomes.CloudFeedAuth) (twomes.CloudFeedAuth, error) {
	return s.cloudFeedAuthRepo.Find(cloudFeedAuth)
}

func exchangeAuthCode(ctx context.Context, conf *oauth2.Config, code string) (string, string, time.Time, error) {
	token, err := conf.Exchange(ctx, code, oauth2.AccessTypeOffline)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return token.AccessToken, token.RefreshToken, token.Expiry, nil
}
