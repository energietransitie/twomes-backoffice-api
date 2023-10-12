package repositories

import (
	"github.com/energietransitie/twomes-backoffice-api/internal/encryption"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
	"gorm.io/gorm"
)

type CloudFeedRepository struct {
	db *gorm.DB
}

// Create a new CloudFeedRepository.
func NewCloudFeedRepository(db *gorm.DB) *CloudFeedRepository {
	return &CloudFeedRepository{
		db: db,
	}
}

// Database representation of a [twomes.CloudFeed]
type CloudFeedModel struct {
	gorm.Model
	Name             string `gorm:"unique;not null"`
	AuthorizationURL string
	TokenURL         string
	ClientID         string
	// TODO: WARNING EncryptedString still has to implement the encryption.
	ClientSecret   encryption.EncryptedString
	Scope          string
	RedirectURL    string
	CloudFeedAuths []CloudFeedAuthModel `gorm:"foreignKey:CloudFeedID"`
}

// Set the name of the table in the database.
func (CloudFeedModel) TableName() string {
	return "cloud_feed"
}

// Create a CloudFeedModel from a [twomes.CloudFeed].
func MakeCloudFeedModel(cloudFeed twomes.CloudFeed) CloudFeedModel {
	return CloudFeedModel{
		Model:            gorm.Model{ID: cloudFeed.ID},
		Name:             cloudFeed.Name,
		AuthorizationURL: cloudFeed.AuthorizationURL,
		TokenURL:         cloudFeed.TokenURL,
		ClientID:         cloudFeed.ClientID,
		ClientSecret:     encryption.EncryptedString(cloudFeed.ClientSecret),
		Scope:            cloudFeed.Scope,
		RedirectURL:      cloudFeed.RedirectURL,
	}
}

// Create a [twomes.CloudFeed] from a CloudFeedModel.
func (m *CloudFeedModel) fromModel() twomes.CloudFeed {
	return twomes.CloudFeed{
		ID:               m.Model.ID,
		Name:             m.Name,
		AuthorizationURL: m.AuthorizationURL,
		TokenURL:         m.TokenURL,
		ClientID:         m.ClientID,
		ClientSecret:     string(m.ClientSecret),
		Scope:            m.Scope,
		RedirectURL:      m.RedirectURL,
	}
}

func (r *CloudFeedRepository) Find(cloudFeed twomes.CloudFeed) (twomes.CloudFeed, error) {
	cloudFeedModel := MakeCloudFeedModel(cloudFeed)
	err := r.db.Where(&cloudFeedModel).First(&cloudFeedModel).Error
	return cloudFeedModel.fromModel(), err
}

func (r *CloudFeedRepository) GetAll() ([]twomes.CloudFeed, error) {
	var cloudFeeds []twomes.CloudFeed

	var cloudFeedModels []CloudFeedModel
	err := r.db.Find(&cloudFeedModels).Error
	if err != nil {
		return nil, err
	}

	for _, cloudFeedModel := range cloudFeedModels {
		cloudFeeds = append(cloudFeeds, cloudFeedModel.fromModel())
	}

	return cloudFeeds, nil
}

func (r *CloudFeedRepository) Create(cloudFeed twomes.CloudFeed) (twomes.CloudFeed, error) {
	cloudFeedModel := MakeCloudFeedModel(cloudFeed)
	err := r.db.Create(&cloudFeedModel).Error
	return cloudFeedModel.fromModel(), err
}

func (r *CloudFeedRepository) Update(cloudFeed twomes.CloudFeed) (twomes.CloudFeed, error) {
	cloudFeedModel := MakeCloudFeedModel(cloudFeed)
	err := r.db.Model(&cloudFeedModel).Updates(cloudFeedModel).Error
	return cloudFeedModel.fromModel(), err
}

func (r *CloudFeedRepository) Delete(cloudFeed twomes.CloudFeed) error {
	cloudFeedModel := MakeCloudFeedModel(cloudFeed)
	return r.db.Delete(&cloudFeedModel).Error
}
