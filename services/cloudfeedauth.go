package services

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/services/cloudfeeds/enelogic"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const (
	CloudFeedDownloadInterval = time.Hour * 24
)

var (
	ErrDuplicateCloudFeedAuth = errors.New("duplicate cloud feed auth")
)

type CloudFeedAuthService struct {
	cloudFeedAuthRepo ports.CloudFeedAuthRepository
	cloudFeedRepo     ports.CloudFeedRepository
	uploadService     ports.UploadService
	updateChan        chan struct{}
}

// Create a new CloudFeedAuthService.
func NewCloudFeedAuthService(cloudFeedAuthRepo ports.CloudFeedAuthRepository, cloudFeedRepo ports.CloudFeedRepository, uploadService ports.UploadService) *CloudFeedAuthService {
	return &CloudFeedAuthService{
		cloudFeedAuthRepo: cloudFeedAuthRepo,
		cloudFeedRepo:     cloudFeedRepo,
		uploadService:     uploadService,
		updateChan:        make(chan struct{}, 1),
	}
}

// Create a new cloudFeedAuth.
// This function exchanges the AuthGrantToken (Code) for a access and refresh token.
func (s *CloudFeedAuthService) Create(ctx context.Context, accountID, cloudFeedID uint, authGrantToken string) (twomes.CloudFeedAuth, error) {
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

	accessToken, refreshToken, expiry, err := exchangeAuthCode(ctx, conf, authGrantToken)
	if err != nil {
		return twomes.CloudFeedAuth{}, err
	}

	cloudFeedAuth := twomes.MakeCloudFeedAuth(accountID, cloudFeedID, accessToken, refreshToken, expiry, authGrantToken)

	cloudFeedAuth, err = s.cloudFeedAuthRepo.Create(cloudFeedAuth)
	if err != nil {
		return cloudFeedAuth, err
	}

	// Signal an update
	s.updateChan <- struct{}{}

	return cloudFeedAuth, nil
}

// Find a cloudFeedAuth using any field set in the cloudFeedAuth struct.
func (s *CloudFeedAuthService) Find(cloudFeedAuth twomes.CloudFeedAuth) (twomes.CloudFeedAuth, error) {
	return s.cloudFeedAuthRepo.Find(cloudFeedAuth)
}

// Refresh the tokens for the CloudFeedAuth corresponding to accountID and cloudFeedID.
func (s *CloudFeedAuthService) RefreshTokens(ctx context.Context, accountID uint, cloudFeedID uint) (twomes.CloudFeedAuth, error) {
	logrus.Infoln("refreshing token for accountID", accountID, "cloudFeedID", cloudFeedID)

	cloudFeedAuth, err := s.cloudFeedAuthRepo.Find(twomes.CloudFeedAuth{AccountID: accountID, CloudFeedID: cloudFeedID})
	if err != nil {
		return twomes.CloudFeedAuth{}, err
	}

	tokenURL, refreshToken, clientID, clientSecret, err := s.cloudFeedAuthRepo.FindOAuthInfo(accountID, cloudFeedID)
	if err != nil {
		return twomes.CloudFeedAuth{}, err
	}

	if refreshToken == "" {
		return twomes.CloudFeedAuth{}, errors.New("refresh token empty")
	}

	u, err := url.Parse(tokenURL)
	if err != nil {
		return twomes.CloudFeedAuth{}, err
	}

	form := url.Values{}
	form.Add("grant_type", "refresh_token")
	form.Add("refresh_token", refreshToken)
	form.Add("client_id", clientID)
	form.Add("client_secret", clientSecret)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return twomes.CloudFeedAuth{}, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return twomes.CloudFeedAuth{}, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return twomes.CloudFeedAuth{}, errors.New("error reading response from token endpoint")
	}

	response := struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    uint   `json:"expires_in"`
		Error        string `json:"error"`
	}{}
	respBodyReader := bytes.NewReader(respBody)
	err = json.NewDecoder(respBodyReader).Decode(&response)
	if err != nil {
		return twomes.CloudFeedAuth{}, err
	}

	if resp.StatusCode != http.StatusOK {
		// Delete auth since we can not recover from "invalid_grant" error.
		if response.Error == "invalid_grant" {
			logrus.Warnln("deleting invalid cloud feed auth for accountID", accountID, "cloudFeedID", cloudFeedID)
			err := s.cloudFeedAuthRepo.Delete(twomes.CloudFeedAuth{AccountID: accountID, CloudFeedID: cloudFeedID})
			if err != nil {
				return twomes.CloudFeedAuth{}, fmt.Errorf("error deleting invalid auth: %w", err)
			}
		}

		return twomes.CloudFeedAuth{}, fmt.Errorf("unsuccessful refresh request. request: %s", string(respBody))
	}

	cloudFeedAuth.AccessToken = response.AccessToken
	cloudFeedAuth.RefreshToken = response.RefreshToken
	cloudFeedAuth.Expiry = time.Now().Add(time.Second * time.Duration(response.ExpiresIn))

	return s.cloudFeedAuthRepo.Update(cloudFeedAuth)
}

// Run this function in a goroutine to keep tokens refreshed before they expire.
// The preRenewalDuration sets the time we need to refresh the tokens in advance of theri expiry.
func (s *CloudFeedAuthService) RefreshTokensInBackground(ctx context.Context, preRenewalDuration time.Duration) {
refreshLoop:
	for {
		accountID, cloudFeedID, expiry, err := s.cloudFeedAuthRepo.FindFirstTokenToExpire()
		if err != nil {
			logrus.Infoln("no cloud feed auths found in database. not doing anything until one is added")
			select {
			case <-s.updateChan:
				logrus.Infoln("a new cloud feed auth was added. re-checking first expiring token")
			case <-ctx.Done():
				break refreshLoop
			}
			continue
		}

		timerDuration := time.Until(expiry) - preRenewalDuration
		if timerDuration < 0 {
			// Wait 10 seconds to prevent a possible flood of refresh requests.
			time.Sleep(time.Second * 10)

			_, err = s.RefreshTokens(ctx, accountID, cloudFeedID)
			if err != nil {
				logrus.Warningln(err)
			}
			continue
		}

		expiryTimer := time.NewTimer(timerDuration)

		logrus.Infof("waiting %s to refresh first expiring token", timerDuration.String())

		select {
		case <-expiryTimer.C:
			_, err = s.RefreshTokens(ctx, accountID, cloudFeedID)
			if err != nil {
				logrus.Warningln(err)
			}
		case <-s.updateChan:
			logrus.Infoln("a new cloud feed auth was added. re-checking first expiring cloud feed auth token")
			expiryTimer.Stop()
		case <-ctx.Done():
			break refreshLoop
		}
	}
}

func exchangeAuthCode(ctx context.Context, conf *oauth2.Config, code string) (string, string, time.Time, error) {
	token, err := conf.Exchange(ctx, code, oauth2.AccessTypeOffline)
	if err != nil {
		return "", "", time.Time{}, err
	}

	return token.AccessToken, token.RefreshToken, token.Expiry, nil
}

// Run this function in a goroutine to periodically download data from the cloud feed.
// downloadStartTime is the time at which the data should be downloaded and repeated each day.
func (s *CloudFeedAuthService) DownloadInBackground(ctx context.Context, downloadStartTime time.Time) {
	waitTime := time.Until(downloadStartTime)
	startTimer := time.NewTimer(waitTime)

	logrus.Infoln("waiting", waitTime.String(), "to start downloading data from cloud feeds")

	select {
	case <-startTimer.C:
		err := s.download(ctx)
		if err != nil {
			logrus.Errorln(err)
		}
	case <-ctx.Done():
		return
	}

	ticker := time.NewTicker(CloudFeedDownloadInterval)

	for {
		waitTime = time.Until(time.Now().Add(CloudFeedDownloadInterval))
		logrus.Infoln("waiting", waitTime.String(), "to download data from cloud feeds")

		select {
		case <-ticker.C:
			err := s.download(ctx)
			if err != nil {
				logrus.Errorln(err)
			}
		case <-ctx.Done():
			return
		}
	}
}

func (s *CloudFeedAuthService) download(ctx context.Context) error {
	cloudFeedAuths, err := s.cloudFeedAuthRepo.GetAll()
	if err != nil {
		return err
	}

	logrus.Infoln("starting download of data from cloud feeds")

	for _, cfa := range cloudFeedAuths {
		logrus.Infoln("downloading data from cloud feed auth with accountID", cfa.AccountID, "cloudFeedID", cfa.CloudFeedID)

		args := enelogic.DownloadArgs{
			RequestMonthsDatapoints: enelogic.Timespan{
				From: time.Now().AddDate(-2, 0, 0),
				To:   time.Now(),
			},
			RequestDaysDatapoints:     enelogic.Timespan{},
			RequestIntervalDatapoints: enelogic.Timespan{},
		}
		measurements, err := enelogic.Download(ctx, cfa.AccessToken, args)
		if err != nil {
			if err == enelogic.ErrNoData {
				logrus.Infoln("no data found for cloud feed auth with accountID", cfa.AccountID, "cloudFeedID", cfa.CloudFeedID)
				continue
			}
			return err
		}

		if len(measurements) == 0 {
			logrus.Infoln("no data found for cloud feed auth with accountID", cfa.AccountID, "cloudFeedID", cfa.CloudFeedID)
			continue
		}

		device, err := s.cloudFeedAuthRepo.FindDevice(cfa)
		if err != nil {
			logrus.Warningln("error finding device for cloud feed auth:", err)
			continue
		}

		upload, err := s.uploadService.Create(device.ID, twomes.Time(time.Now()), measurements)
		if err != nil {
			logrus.Warningln("error creating upload:", err)
			continue
		}

		if upload.Size != len(measurements) {
			logrus.Errorln("upload size", upload.Size, "does not match number of measurements", len(measurements))
			continue
		}
	}

	return nil
}
