package repositories

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/internal/encryption"
	"github.com/energietransitie/twomes-backoffice-api/internal/helpers"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"gorm.io/gorm"
)

type CloudFeedAuthRepository struct {
	db *gorm.DB
}

func NewCloudFeedAuthRepository(db *gorm.DB) *CloudFeedAuthRepository {
	return &CloudFeedAuthRepository{
		db: db,
	}
}

// Database representation of [twomes.CloudFeedAuth].
type CloudFeedAuthModel struct {
	AccountID   uint `gorm:"primaryKey;autoIncrement:false"`
	CloudFeedID uint `gorm:"primaryKey;autoIncrement:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	// TODO: WARNING encrypted string encryption not yet implemented.
	AccessToken    encryption.EncryptedString
	RefreshToken   encryption.EncryptedString
	Expiry         time.Time
	AuthGrantToken encryption.EncryptedString
}

// Set the name of the table in the database.
func (CloudFeedAuthModel) TableName() string {
	return "cloud_feed_auth"
}

// Create a new CloudFeedAuthModel from a [twomes.cloudFeedAuth].
func MakeCloudFeedAuthModel(cloudFeedAuth twomes.CloudFeedAuth) CloudFeedAuthModel {
	return CloudFeedAuthModel{
		AccountID:      cloudFeedAuth.AccountID,
		CloudFeedID:    cloudFeedAuth.CloudFeedID,
		AccessToken:    encryption.EncryptedString(cloudFeedAuth.AccessToken),
		RefreshToken:   encryption.EncryptedString(cloudFeedAuth.RefreshToken),
		Expiry:         cloudFeedAuth.Expiry,
		AuthGrantToken: encryption.EncryptedString(cloudFeedAuth.AuthGrantToken),
	}
}

// Create a [twomes.CloudFeedAuth] from an CloudFeedAuthModel.
func (m *CloudFeedAuthModel) fromModel() twomes.CloudFeedAuth {
	return twomes.CloudFeedAuth{
		AccountID:      m.AccountID,
		CloudFeedID:    m.CloudFeedID,
		AccessToken:    string(m.AccessToken),
		RefreshToken:   string(m.RefreshToken),
		Expiry:         m.Expiry,
		AuthGrantToken: string(m.AuthGrantToken),
	}
}

func (r *CloudFeedAuthRepository) Find(cloudFeedAuth twomes.CloudFeedAuth) (twomes.CloudFeedAuth, error) {
	cloudFeedAuthModel := MakeCloudFeedAuthModel(cloudFeedAuth)
	err := r.db.Where(&cloudFeedAuthModel).First(&cloudFeedAuthModel).Error
	return cloudFeedAuthModel.fromModel(), err
}

func (r *CloudFeedAuthRepository) FindOAuthInfo(accountID uint, cloudFeedID uint) (string, string, string, string, error) {
	var result struct {
		TokenURL     string
		RefreshToken string
		ClientID     string
		ClientSecret string
	}
	err := r.db.Table("cloud_feed").Select("cloud_feed.token_url, cloud_feed_auth.refresh_token AS refresh_token, cloud_feed.client_id, cloud_feed.client_secret").Joins("JOIN cloud_feed_auth ON cloud_feed.id = cloud_feed_auth.cloud_feed_id").Where("cloud_feed_auth.account_id = ? AND cloud_feed_auth.cloud_feed_id = ?", accountID, cloudFeedID).Scan(&result).Error
	return result.TokenURL, result.RefreshToken, result.ClientID, result.ClientSecret, err
}

func (r *CloudFeedAuthRepository) FindFirstTokenToExpire() (uint, uint, time.Time, error) {
	var cloudFeedAuthModel CloudFeedAuthModel
	err := r.db.Order("expiry ASC").Where("expiry <> ''").First(&cloudFeedAuthModel).Error
	return cloudFeedAuthModel.AccountID, cloudFeedAuthModel.CloudFeedID, cloudFeedAuthModel.Expiry, err
}

func (r *CloudFeedAuthRepository) FindDevice(cloudFeedAuth twomes.CloudFeedAuth) (*twomes.Device, error) {
	var device DeviceModel

	// err := r.db.Table("cloud_feed_auth").
	// 	Select("device.*").
	// 	Joins("JOIN cloud_feed ON cloud_feed_auth.cloud_feed_id = cloud_feed.id").
	// 	Joins("JOIN device_type ON cloud_feed.name = device_type.name").
	// 	Joins("JOIN device ON device_type.id = device.device_type_id").
	// 	Joins("JOIN building ON device.building_id = building.id").
	// 	Where("building.account_id = cloud_feed_auth.account_id AND cloud_feed_auth.account_id = ? AND cloud_feed_auth.cloud_feed_id = ?", cloudFeedAuth.AccountID, cloudFeedAuth.CloudFeedID).
	// 	Order("cloud_feed_auth.account_id DESC").
	// 	First(&device).
	// 	Error

	err := r.db.Model(&device).
		Joins("JOIN device_type ON device_type.id = device.device_type_id").
		Joins("JOIN cloud_feed ON cloud_feed.name = device_type.name").
		Joins("JOIN cloud_feed_auth ON cloud_feed_auth.cloud_feed_id = cloud_feed.id").
		Joins("JOIN building ON building.id = device.building_id").
		Where("building.account_id = cloud_feed_auth.account_id AND cloud_feed_auth.account_id = ? AND cloud_feed_auth.cloud_feed_id = ?", cloudFeedAuth.AccountID, cloudFeedAuth.CloudFeedID).
		First(&device).
		Error

	if err != nil {
		return nil, err
	}

	d := device.fromModel()
	return &d, nil
}

func (r *CloudFeedAuthRepository) GetAll() ([]twomes.CloudFeedAuth, error) {
	var cloudFeedAuths []twomes.CloudFeedAuth

	var cloudFeedAuthModels []CloudFeedAuthModel
	err := r.db.Find(&cloudFeedAuthModels).Error
	if err != nil {
		return nil, err
	}

	for _, cloudFeedAuthModel := range cloudFeedAuthModels {
		cloudFeedAuths = append(cloudFeedAuths, cloudFeedAuthModel.fromModel())
	}

	return cloudFeedAuths, nil
}

func (r *CloudFeedAuthRepository) Create(cloudFeedAuth twomes.CloudFeedAuth) (twomes.CloudFeedAuth, error) {
	// First check if a soft deleted entry exists.
	// We need to do this because we can't create a new one if one exists.
	cfaCheck := CloudFeedAuthModel{}
	err := r.db.Unscoped().Where(&CloudFeedAuthModel{
		AccountID:   cloudFeedAuth.AccountID,
		CloudFeedID: cloudFeedAuth.CloudFeedID,
	}).First(&cfaCheck).Error

	// Return immediately if there is an error,
	// except if the error was RecordNotFound.
	if err != nil && !helpers.IsMySQLRecordNotFoundError(err) {
		return cloudFeedAuth, err
	}

	// At this point we know that there are no errors,
	// except maybe that the record was not found, so check.
	if !helpers.IsMySQLRecordNotFoundError(err) {
		// Record was found. Check if it was soft deleted.
		if cfaCheck.DeletedAt.Valid {
			// Record was soft deleted. Delete it so we can create a new one.
			err := r.db.Unscoped().Delete(&cfaCheck).Error
			if err != nil {
				return cloudFeedAuth, err
			}
		}
	}

	// If we reach this, it means there was not previous record,
	// or it was deleted, and we can just create a new one.
	cloudFeedAuthModel := MakeCloudFeedAuthModel(cloudFeedAuth)
	err = r.db.Create(&cloudFeedAuthModel).Error
	return cloudFeedAuthModel.fromModel(), err
}

func (r *CloudFeedAuthRepository) Update(cloudFeedAuth twomes.CloudFeedAuth) (twomes.CloudFeedAuth, error) {
	cloudFeedAuthModel := MakeCloudFeedAuthModel(cloudFeedAuth)
	err := r.db.Model(&cloudFeedAuthModel).Updates(cloudFeedAuthModel).Error
	return cloudFeedAuthModel.fromModel(), err
}

func (r *CloudFeedAuthRepository) Delete(cloudFeedAuth twomes.CloudFeedAuth) error {
	CloudFeedAuthModel := MakeCloudFeedAuthModel(cloudFeedAuth)
	return r.db.Delete(&CloudFeedAuthModel).Error
}
