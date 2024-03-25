package repositories

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/campaign"
	"github.com/energietransitie/twomes-backoffice-api/twomes/cloudfeed"
	"gorm.io/gorm"
)

type CampaignRepository struct {
	db *gorm.DB
}

func NewCampaignRepository(db *gorm.DB) *CampaignRepository {
	return &CampaignRepository{
		db: db,
	}
}

// Database representation of [campaign.Campaign].
type CampaignModel struct {
	gorm.Model
	Name           string `gorm:"unique;not null"`
	AppModelID     uint   `gorm:"column:app_id"`
	App            AppModel
	InfoURL        string           `gorm:"unique;not null"`
	CloudFeeds     []CloudFeedModel `gorm:"many2many:campaign_cloud_feed"`
	StartTime      *time.Time
	EndTime        *time.Time
	ShoppingListID uint
	ShoppingList   ShoppingListModel `gorm:"foreignkey:ShoppingListID"`
}

// Set the name of the table in the database.
func (CampaignModel) TableName() string {
	return "campaign"
}

// Create a new CampaignModel from a [twomes.campaign].
func MakeCampaignModel(campaign campaign.Campaign) CampaignModel {
	var cloudFeedModels []CloudFeedModel

	for _, cloudFeed := range campaign.CloudFeeds {
		cloudFeedModels = append(cloudFeedModels, MakeCloudFeedModel(cloudFeed))
	}

	return CampaignModel{
		Model: gorm.Model{
			ID: campaign.ID,
		},
		Name:           campaign.Name,
		AppModelID:     campaign.App.ID,
		App:            MakeAppModel(campaign.App),
		InfoURL:        campaign.InfoURL,
		CloudFeeds:     cloudFeedModels,
		StartTime:      campaign.StartTime,
		EndTime:        campaign.EndTime,
		ShoppingListID: campaign.ShoppingList.ID,
		ShoppingList:   MakeShoppingListModel(campaign.ShoppingList),
	}
}

// Create a [campaign.Campaign] from an CampaignModel.
func (m *CampaignModel) fromModel() campaign.Campaign {
	var cloudFeeds []cloudfeed.CloudFeed

	for _, cloudFeedModel := range m.CloudFeeds {
		cloudFeeds = append(cloudFeeds, cloudFeedModel.fromModel())
	}

	return campaign.Campaign{
		ID:           m.ID,
		Name:         m.Name,
		App:          m.App.fromModel(),
		InfoURL:      m.InfoURL,
		CloudFeeds:   cloudFeeds,
		StartTime:    m.StartTime,
		EndTime:      m.EndTime,
		ShoppingList: m.ShoppingList.fromModel(),
	}
}

func (r *CampaignRepository) Find(campaign campaign.Campaign) (campaign.Campaign, error) {
	campaignModel := MakeCampaignModel(campaign)
	err := r.db.Preload("App").Preload("ShoppingList").Where(&campaignModel).First(&campaignModel).Error
	return campaignModel.fromModel(), err
}

func (r *CampaignRepository) GetAll() ([]campaign.Campaign, error) {
	var campaigns []campaign.Campaign

	var campaignModels []CampaignModel
	err := r.db.Preload("App").Preload("ShoppingList").Find(&campaignModels).Error
	if err != nil {
		return nil, err
	}

	for _, campaignModel := range campaignModels {
		campaigns = append(campaigns, campaignModel.fromModel())
	}

	return campaigns, nil
}

func (r *CampaignRepository) Create(campaign campaign.Campaign) (campaign.Campaign, error) {
	campaignModel := MakeCampaignModel(campaign)
	err := r.db.Create(&campaignModel).Error
	return campaignModel.fromModel(), err
}

func (r *CampaignRepository) Delete(campaign campaign.Campaign) error {
	CampaignModel := MakeCampaignModel(campaign)
	return r.db.Delete(&CampaignModel).Error
}
