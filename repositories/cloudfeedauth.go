package repositories

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/internal/encryption"
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
	// TODO: WARNING encrypted sting encryption not yet implemented.
	AccessToken    encryption.EncrpytedString
	RefreshToken   encryption.EncrpytedString
	AuthGrantToken encryption.EncrpytedString
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
		AccessToken:    encryption.EncrpytedString(cloudFeedAuth.AccessToken),
		RefreshToken:   encryption.EncrpytedString(cloudFeedAuth.RefreshToken),
		AuthGrantToken: encryption.EncrpytedString(cloudFeedAuth.AuthGrantToken),
	}
}

// Create a [twomes.CloudFeedAuth] from an CloudFeedAuthModel.
func (m *CloudFeedAuthModel) fromModel() twomes.CloudFeedAuth {
	return twomes.CloudFeedAuth{
		AccountID:      m.AccountID,
		CloudFeedID:    m.CloudFeedID,
		AccessToken:    string(m.AccessToken),
		RefreshToken:   string(m.RefreshToken),
		AuthGrantToken: string(m.AuthGrantToken),
	}
}

func (r *CloudFeedAuthRepository) Find(cloudFeedAuth twomes.CloudFeedAuth) (twomes.CloudFeedAuth, error) {
	cloudFeedAuthModel := MakeCloudFeedAuthModel(cloudFeedAuth)
	err := r.db.Where(&cloudFeedAuthModel).First(&cloudFeedAuthModel).Error
	return cloudFeedAuthModel.fromModel(), err
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
	cloudFeedAuthModel := MakeCloudFeedAuthModel(cloudFeedAuth)
	err := r.db.Create(&cloudFeedAuthModel).Error
	return cloudFeedAuthModel.fromModel(), err
}

func (r *CloudFeedAuthRepository) Delete(cloudFeedAuth twomes.CloudFeedAuth) error {
	CloudFeedAuthModel := MakeCloudFeedAuthModel(cloudFeedAuth)
	return r.db.Delete(&CloudFeedAuthModel).Error
}
