// Package repositories implements repositories for use in services.
package repositories

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/twomes"
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

// Database representation of a [twomes.Account].
type AccountModel struct {
	gorm.Model
	CampaignModelID uint `gorm:"column:campaign_id"`
	Campaign        CampaignModel
	ActivatedAt     *time.Time
	Buildings       []BuildingModel
}

// Set the name of the table in the database.
func (AccountModel) TableName() string {
	return "account"
}

// Create a new AccountModel from a [twomes.Account].
func MakeAccountModel(account twomes.Account) AccountModel {
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

// Create a [twomes.Account] from an AccountModel.
func (m *AccountModel) fromModel() twomes.Account {
	var buildings []twomes.Building

	for _, buildingModel := range m.Buildings {
		buildings = append(buildings, buildingModel.fromModel())
	}

	return twomes.Account{
		ID:          m.Model.ID,
		Campaign:    m.Campaign.fromModel(),
		ActivatedAt: m.ActivatedAt,
		Buildings:   buildings,
	}
}

func (r *AccountRepository) Find(account twomes.Account) (twomes.Account, error) {
	accountModel := MakeAccountModel(account)
	err := r.db.Preload("Campaign.App").Preload("Buildings").Where(&accountModel).First(&accountModel).Error
	return accountModel.fromModel(), err
}

func (r *AccountRepository) GetAll() ([]twomes.Account, error) {
	accounts := make([]twomes.Account, 0)

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

func (r *AccountRepository) Create(account twomes.Account) (twomes.Account, error) {
	accountModel := MakeAccountModel(account)
	err := r.db.Create(&accountModel).Error
	return accountModel.fromModel(), err
}

func (r *AccountRepository) Update(account twomes.Account) (twomes.Account, error) {
	accountModel := MakeAccountModel(account)
	err := r.db.Model(&accountModel).Updates(accountModel).Error
	return accountModel.fromModel(), err
}

func (r *AccountRepository) Delete(account twomes.Account) error {
	accountModel := MakeAccountModel(account)
	return r.db.Delete(&accountModel).Error
}
