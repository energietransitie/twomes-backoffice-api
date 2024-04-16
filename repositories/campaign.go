package repositories

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/campaign"
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
	Name             string `gorm:"unique;not null"`
	AppModelID       uint   `gorm:"column:app_id"`
	App              AppModel
	InfoURL          string `gorm:"unique;not null"`
	StartTime        *time.Time
	EndTime          *time.Time
	DataSourceListID uint
}

// Set the name of the table in the database.
func (CampaignModel) TableName() string {
	return "campaign"
}

// Create a new CampaignModel from a [twomes.campaign].
func MakeCampaignModel(campaign campaign.Campaign) CampaignModel {
	return CampaignModel{
		Model: gorm.Model{
			ID: campaign.ID,
		},
		Name:             campaign.Name,
		AppModelID:       campaign.App.ID,
		App:              MakeAppModel(campaign.App),
		InfoURL:          campaign.InfoURL,
		StartTime:        campaign.StartTime,
		EndTime:          campaign.EndTime,
		DataSourceListID: campaign.DataSourceList.ID,
	}
}

// Create a [campaign.Campaign] from an CampaignModel.
func (m *CampaignModel) fromModel(db *gorm.DB) campaign.Campaign {
	return campaign.Campaign{
		ID:             m.ID,
		Name:           m.Name,
		App:            m.App.fromModel(),
		InfoURL:        m.InfoURL,
		StartTime:      m.StartTime,
		EndTime:        m.EndTime,
	}
}

func (r *CampaignRepository) Find(campaign campaign.Campaign) (campaign.Campaign, error) {
	campaignModel := MakeCampaignModel(campaign)
	err := r.db.Preload("App").Preload("DataSourceList").Where(&campaignModel).First(&campaignModel).Error
	return campaignModel.fromModel(r.db), err
}

func (r *CampaignRepository) GetAll() ([]campaign.Campaign, error) {
	var campaigns []campaign.Campaign

	var campaignModels []CampaignModel
	err := r.db.Preload("App").Preload("DataSourceList").Find(&campaignModels).Error
	if err != nil {
		return nil, err
	}

	for _, campaignModel := range campaignModels {
		campaigns = append(campaigns, campaignModel.fromModel(r.db))
	}

	return campaigns, nil
}

func (r *CampaignRepository) Create(campaign campaign.Campaign) (campaign.Campaign, error) {
	campaignModel := MakeCampaignModel(campaign)
	err := r.db.Create(&campaignModel).Error
	return campaignModel.fromModel(r.db), err
}

func (r *CampaignRepository) Delete(campaign campaign.Campaign) error {
	CampaignModel := MakeCampaignModel(campaign)
	return r.db.Delete(&CampaignModel).Error
}
