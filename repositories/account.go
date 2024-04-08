// Package repositories implements repositories for use in services.
package repositories

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes/account"
	"github.com/energietransitie/twomes-backoffice-api/twomes/building"
	"gorm.io/gorm"
)

type AccountRepository struct {
	db *gorm.DB
}

func NewAccountRepository(db *gorm.DB) *AccountRepository {
	return &AccountRepository{
		db: db,
	}
}

// Database representation of a [account.Account].
type AccountModel struct {
	gorm.Model
	CampaignModelID uint `gorm:"column:campaign_id"`
	Campaign        CampaignModel
	ActivatedAt     *time.Time
	Buildings       []BuildingModel
	CloudFeeds      []CloudFeedModel   `gorm:"foreignKey:AccountID"`
	EnergyQueries   []EnergyQueryModel `gorm:"foreignKey:AccountID"`
}

// Set the name of the table in the database.
func (AccountModel) TableName() string {
	return "account"
}

// Create a new AccountModel from a [account.Account].
func MakeAccountModel(account account.Account) AccountModel {
	var buildingModels []BuildingModel

	for _, building := range account.Buildings {
		buildingModels = append(buildingModels, MakeBuildingModel(building))
	}

	return AccountModel{
		Model: gorm.Model{
			ID: account.ID,
		},
		CampaignModelID: account.Campaign.ID,
		Campaign:        MakeCampaignModel(account.Campaign),
		ActivatedAt:     account.ActivatedAt,
		Buildings:       buildingModels,
	}
}

// Create a [account.Account] from an AccountModel.
func (m *AccountModel) fromModel() account.Account {
	var buildings []building.Building

	for _, buildingModel := range m.Buildings {
		buildings = append(buildings, buildingModel.fromModel())
	}

	return account.Account{
		ID:          m.Model.ID,
		Campaign:    m.Campaign.fromModel(),
		ActivatedAt: m.ActivatedAt,
		Buildings:   buildings,
	}
}

func (r *AccountRepository) Find(account account.Account) (account.Account, error) {
	accountModel := MakeAccountModel(account)
	err := r.db.Preload("Campaign.App").Preload("Buildings").Where(&accountModel).First(&accountModel).Error
	return accountModel.fromModel(), err
}

func (r *AccountRepository) GetAll() ([]account.Account, error) {
	accounts := make([]account.Account, 0)

	var accountModels []AccountModel
	err := r.db.Preload("Campaign.App").Preload("Buildings").Find(&accountModels).Error
	if err != nil {
		return nil, err
	}

	for _, accountModel := range accountModels {
		accounts = append(accounts, accountModel.fromModel())
	}

	return accounts, nil
}

func (r *AccountRepository) Create(account account.Account) (account.Account, error) {
	accountModel := MakeAccountModel(account)
	err := r.db.Create(&accountModel).Error
	return accountModel.fromModel(), err
}

func (r *AccountRepository) Update(account account.Account) (account.Account, error) {
	accountModel := MakeAccountModel(account)
	err := r.db.Model(&accountModel).Updates(accountModel).Error
	return accountModel.fromModel(), err
}

func (r *AccountRepository) Delete(account account.Account) error {
	accountModel := MakeAccountModel(account)
	return r.db.Delete(&accountModel).Error
}
